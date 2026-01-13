package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/feeds"
	"github.com/egandro/news-deframer/pkg/think"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"golang.org/x/net/publicsuffix"
)

const defaultLockDuration = 5 * time.Minute
const idleSleepTime = 10 * time.Second
const pollingInterval = 10 * time.Minute
const feedTitlePrefix = "Deframer: "
const promptScope = "deframer"
const customPrefix = "deframer"
const customNamespace = "https://github.com/egandro/news-deframer/"

type FeedSyncer interface {
	SyncFeed(id uuid.UUID) error
	StopPolling(id uuid.UUID) error
	Poll()
}

type Syncer struct {
	ctx    context.Context
	repo   database.Repository
	logger *slog.Logger
	dl     downloader.Downloader
	feeds  feeds.Feeds
	think  think.Think
}

func New(ctx context.Context, cfg *config.Config, repo database.Repository) (*Syncer, error) {
	th, err := think.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Syncer{
		ctx:    ctx,
		repo:   repo,
		logger: slog.With("component", "syncer"),
		dl:     downloader.NewDownloader(ctx, cfg),
		feeds:  feeds.NewFeeds(ctx, cfg),
		think:  th,
	}, nil
}

func (s *Syncer) SyncFeed(id uuid.UUID) error {
	s.logger.Info("Syncing feed", "id", id)
	return s.repo.EnqueueSync(id, 0, defaultLockDuration)
}

func (s *Syncer) StopPolling(id uuid.UUID) error {
	s.logger.Info("Stopping polling", "id", id)
	return s.repo.RemoveSync(id)
}

func (s *Syncer) Poll() {
	s.logger.Info("Starting poller")

	for {
		if s.ctx.Err() != nil {
			s.logger.Info("Stopping poller")
			return
		}

		if s.syncNextScheduledFeed() {
			s.logger.Info("A feed was synced")
			continue
		}

		s.logger.Info("Sleeping...", "duration", idleSleepTime)

		select {
		case <-s.ctx.Done():
			s.logger.Info("Stopping poller")
			return
		case <-time.After(idleSleepTime):
		}
	}
}

// syncNextScheduledFeed return true if this has updated entries
func (s *Syncer) syncNextScheduledFeed() bool {
	s.logger.Info("syncNextScheduledFeed")
	feed, err := s.repo.BeginFeedUpdate(defaultLockDuration)
	if err != nil {
		s.logger.Error("Failed query the next feed to sync", "error", err)
		return false
	}
	if feed == nil {
		return false
	}
	err = s.updatingFeed(feed)
	if err := s.repo.EndFeedUpdate(feed.ID, err, pollingInterval); err != nil {
		s.logger.Error("Failed to end feed update", "error", err)
	}
	return true
}

func (s *Syncer) updatingFeed(feed *database.Feed) error {
	s.logger.Info("Updating feed", "id", feed.ID, "url", feed.URL)

	u, err := url.Parse(feed.URL)
	if err != nil {
		return err
	}

	rc, err := s.dl.DownloadRSSFeed(s.ctx, u)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	parsedFeed, err := s.feeds.ParseFeed(s.ctx, rc)
	if err != nil {
		return err
	}

	domains, err := s.wantedDomains(feed)
	if err != nil {
		return err
	}

	// items we can calculate a hash and it's urls are on our wanted domain list
	items := s.feeds.FilterItems(s.ctx, parsedFeed, domains)
	s.logger.Debug("items", "len", len(items))

	hashes := feeds.GetHashes(items)
	s.logger.Debug("hashes", "len", len(hashes))

	count, err := s.processItems(feed, parsedFeed, items, hashes)
	if err != nil {
		return err
	}

	if count == 0 {
		s.logger.Debug("all items are processed - not updating the feed")
		return nil
	}

	return s.updateCacheFeed(feed, parsedFeed, hashes)
}

func (s *Syncer) wantedDomains(feed *database.Feed) ([]string, error) {
	if !feed.EnforceFeedDomain {
		return nil, nil
	}

	u, err := url.Parse(feed.URL)
	if err != nil {
		return nil, err
	}

	baseDomain, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname())
	if err != nil {
		baseDomain = u.Hostname()
	}

	return []string{baseDomain}, nil
}

func (s *Syncer) processItems(feed *database.Feed, parsedFeed *gofeed.Feed, items []feeds.ItemHashPair, hashes []string) (int, error) {
	pendingHashes, err := s.repo.GetPendingHashes(feed.ID, hashes)
	if err != nil {
		return 0, err
	}
	count := len(pendingHashes)
	s.logger.Debug("pending hashes", "len", count)

	if count == 0 {
		s.logger.Debug("All items are already processed")
		return 0, nil
	}

	language := "en"
	if l := strings.TrimSpace(parsedFeed.Language); l != "" {
		// this is an iso code for the language split at the "-" we only want the first place
		language = strings.ToLower(strings.Split(l, "-")[0])
	}

	count = 0
	for _, item := range items {
		if s.ctx.Err() != nil {
			// context might be canceled
			// we have processed x items - so report a 0,nil that
			// don't create a result in postgres - the next tick will pick it up
			return 0, nil
		}
		if pendingHashes[item.Hash] {
			s.processItem(feed, item.Hash, item.Item, language)
			count++
		}
	}

	return count, nil
}

func (s *Syncer) processItem(feed *database.Feed, hash string, item *gofeed.Item, language string) {
	req := think.Request{
		Title:       item.Title,
		Description: item.Description,
	}
	res, err := s.think.Run(promptScope, language, req)

	var thinkError *string
	var mediaContent *database.MediaContent
	var thinkRating float64

	if err != nil {
		if errors.Is(err, context.Canceled) {
			// worker is shutting down don't store this as an error
			s.logger.Debug("context canceled", "hash", hash)
			return
		}
		s.logger.Error("analysis failed", "error", err, "hash", hash, "item_url", item.Link)
		errStr := err.Error()
		thinkError = &errStr
	} else {
		err = s.updateContent(item, res)
		if err != nil {
			s.logger.Error("update content failed", "error", err)
			return
		}
		mediaContent, err = s.extractMediaContent(item)
		if err != nil {
			s.logger.Error("extract media content failed", "error", err)
		}
		thinkRating = calculateHybrid(res, 0.7)
		applyFancyRatingText(item, res, thinkRating, language)
	}

	content, err := s.feeds.RenderItem(s.ctx, item)
	if err != nil {
		s.logger.Error("failed to render item", "error", err, "hash", hash)
		return
	}

	pubDate := time.Now()
	if item.PublishedParsed != nil {
		pubDate = *item.PublishedParsed
	}

	dbItem := &database.Item{
		Hash:         hash,
		FeedID:       feed.ID,
		URL:          item.Link,
		Content:      content,
		PubDate:      pubDate,
		ThinkResult:  res,
		MediaContent: mediaContent,
		ThinkError:   thinkError,
		ThinkRating:  thinkRating,
	}

	s.logger.Debug("processItem", "feed", feed.ID, "hash", hash)

	if err := s.repo.UpsertItem(dbItem); err != nil {
		s.logger.Error("failed to create item", "error", err, "hash", hash)
	}
}

func (s *Syncer) updateContent(item *gofeed.Item, res *database.ThinkResult) error {
	if item == nil || res == nil {
		return nil
	}

	// save the original title
	res.TitleOriginal = item.Title
	res.DescriptionOriginal = item.Description

	// correct the title
	item.Title = res.TitleCorrected
	item.Description = res.DescriptionCorrected

	// Add the reason
	item.Description = fmt.Sprintf("%s<br/><br/>%s", item.Description, res.OverallReason)

	if item.Content != "" {
		// tagesschau uses "content:encoded"
		// let's try to grab any image and keep the title/description
		mediaData, err := transformContent(item.Content)
		if err == nil && mediaData.URL != "" {
			// s.logger.Debug("removed content:encoded", "url", item.Link)
			item.Content = ""

			if item.Extensions == nil {
				item.Extensions = make(map[string]map[string][]ext.Extension)
			}
			if item.Extensions["media"] == nil {
				item.Extensions["media"] = make(map[string][]ext.Extension)
			}

			// <media:credit> always delete this
			delete(item.Extensions["media"], "credit")

			attrs := map[string]string{
				"url":    mediaData.URL,
				"medium": mediaData.Medium,
			}
			if mediaData.Width > 0 {
				attrs["width"] = fmt.Sprintf("%d", mediaData.Width)
			}
			if mediaData.Height > 0 {
				attrs["height"] = fmt.Sprintf("%d", mediaData.Height)
			}

			item.Extensions["media"]["content"] = []ext.Extension{{
				Name:  "content",
				Attrs: attrs,
			}}

			if len(mediaData.Alt) > 0 {
				item.Extensions["media"]["description"] = []ext.Extension{{
					Name:  "description",
					Value: mediaData.Alt,
				}}
			}
		}
	}

	if item.Extensions != nil {
		if mediaExt, ok := item.Extensions["media"]; ok {
			if _, hasGroup := mediaExt["group"]; hasGroup {
				delete(mediaExt, "group")
				s.logger.Info("removed media:group tag", "url", item.Link)
			}
		}
	}

	return nil
}

func (s *Syncer) extractMediaContent(item *gofeed.Item) (*database.MediaContent, error) {
	if item == nil || item.Extensions == nil {
		return nil, nil
	}

	mediaExt, ok := item.Extensions["media"]
	if !ok {
		return nil, nil
	}

	// Find media:content
	var contentExt *ext.Extension
	if contents, ok := mediaExt["content"]; ok {
		for _, c := range contents {
			if c.Attrs["url"] != "" {
				val := c
				contentExt = &val
				break
			}
		}
	}

	if contentExt == nil {
		return nil, nil
	}

	mc := &database.MediaContent{
		URL:    contentExt.Attrs["url"],
		Type:   contentExt.Attrs["type"],
		Medium: contentExt.Attrs["medium"],
	}

	// Dimensions
	w, _ := strconv.Atoi(contentExt.Attrs["width"])
	h, _ := strconv.Atoi(contentExt.Attrs["height"])

	if w == 0 || h == 0 {
		w, h = parseDimensions(mc.URL)
	}
	mc.Width = w
	mc.Height = h

	// Metadata (Title/Description)
	// Check children first
	if titles, ok := contentExt.Children["title"]; ok && len(titles) > 0 {
		mc.Title = titles[0].Value
	}
	if descs, ok := contentExt.Children["description"]; ok && len(descs) > 0 {
		mc.Description = descs[0].Value
	}
	if credits, ok := contentExt.Children["credit"]; ok && len(credits) > 0 {
		mc.Credit = credits[0].Value
	}

	// Fallback to group/item level tags
	if mc.Title == "" {
		if titles, ok := mediaExt["title"]; ok && len(titles) > 0 {
			mc.Title = titles[0].Value
		}
	}
	if mc.Description == "" {
		if descs, ok := mediaExt["description"]; ok && len(descs) > 0 {
			mc.Description = descs[0].Value
		}
	}
	if mc.Credit == "" {
		if credits, ok := mediaExt["credit"]; ok && len(credits) > 0 {
			mc.Credit = credits[0].Value
		}
	}

	// Thumbnail
	if thumbs, ok := mediaExt["thumbnail"]; ok && len(thumbs) > 0 {
		tExt := thumbs[0]
		if tURL := tExt.Attrs["url"]; tURL != "" {
			mt := &database.MediaThumbnail{
				URL: tURL,
			}
			tw, _ := strconv.Atoi(tExt.Attrs["width"])
			th, _ := strconv.Atoi(tExt.Attrs["height"])

			if tw == 0 || th == 0 {
				tw, th = parseDimensions(tURL)
			}
			mt.Width = tw
			mt.Height = th
			mc.Thumbnail = mt
		}
	}

	return mc, nil
}

func (s *Syncer) updateCacheFeed(feed *database.Feed, parsedFeed *gofeed.Feed, hashes []string) error {
	// load all items from database by their hashes and with the

	items, err := s.repo.GetItemsByHashes(feed.ID, hashes)
	if err != nil {
		return err
	}
	s.logger.Debug("items", "len", len(items))

	// delete the items - we will replace them
	parsedFeed.Items = []*gofeed.Item{}
	// custom namespace
	feeds.AddNamespace(parsedFeed, "xmlns:"+customPrefix, customNamespace)
	// media namespace
	feeds.AddNamespace(parsedFeed, "xmlns:media", "http://search.yahoo.com/mrss/")

	fp := gofeed.NewParser()
	for _, item := range items {
		wrapped := `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel>` + item.Content + `</channel></rss>`
		pf, err := fp.Parse(strings.NewReader(wrapped))
		if err != nil {
			s.logger.Error("failed to parse item xml", "error", err, "hash", item.Hash)
			continue
		}

		if len(pf.Items) == 0 {
			continue
		}

		current := pf.Items[0]

		if item.ThinkResult != nil && item.ThinkError == nil {
			res := *item.ThinkResult
			var m map[string]interface{}
			b, _ := json.Marshal(res)
			_ = json.Unmarshal(b, &m)
			for k, v := range m {
				feeds.SetExtension(current, customPrefix, k, fmt.Sprintf("%v", v))
			}
		}

		parsedFeed.Items = append(parsedFeed.Items, current)
	}

	parsedFeed.Title = feedTitlePrefix + parsedFeed.Title

	xmlContent, err := s.feeds.RenderFeed(s.ctx, parsedFeed)
	if err != nil {
		return err
	}
	s.logger.Info("Rendered feed", "len", len(xmlContent))

	cachedFeed := &database.CachedFeed{
		ID:         feed.ID,
		XMLContent: &xmlContent,
		ItemRefs:   database.StringArray(hashes),
	}

	return s.repo.UpsertCachedFeed(cachedFeed)
}
