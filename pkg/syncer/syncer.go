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
	"sync"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/downloader"
	"github.com/deframer/news-deframer/pkg/feeds"
	"github.com/deframer/news-deframer/pkg/think"
	"github.com/deframer/news-deframer/pkg/util/netutil"
	"github.com/deframer/news-deframer/pkg/util/text"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"golang.org/x/net/publicsuffix"
)

const feedTitlePrefix = "Deframer: "
const promptScope = "deframer"
const customPrefix = "deframer"
const customNamespace = "https://github.com/deframer/news-deframer/"
const maxThinkRetries = 3

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
	return s.repo.EnqueueSync(id, 0)
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

		s.logger.Info("Sleeping...", "duration", config.IdleSleepTime)

		select {
		case <-s.ctx.Done():
			s.logger.Info("Stopping poller")
			return
		case <-time.After(config.IdleSleepTime):
		}
	}
}

// syncNextScheduledFeed return true if this has updated entries
func (s *Syncer) syncNextScheduledFeed() bool {
	s.logger.Info("syncNextScheduledFeed")
	feed, err := s.repo.BeginFeedUpdate(config.DefaultLockDuration)
	if err != nil {
		s.logger.Error("Failed query the next feed to sync", "error", err)
		return false
	}
	if feed == nil {
		return false
	}
	err = s.updatingFeed(feed)
	if err := s.repo.EndFeedUpdate(feed.ID, err, config.PollingInterval); err != nil {
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

	if feed.ResolveItemUrl {
		var wg sync.WaitGroup
		// Limit concurrency to avoid overwhelming the server or local resources
		sem := make(chan struct{}, 10)

		for _, item := range parsedFeed.Items {
			wg.Add(1)
			go func(item *gofeed.Item) {
				defer wg.Done()
				sem <- struct{}{}        // Acquire token
				defer func() { <-sem }() // Release token

				resolved, err := s.dl.ResolveRedirect(s.ctx, item.Link)
				if err != nil {
					s.logger.Warn("failed to resolve redirect", "url", item.Link, "error", err)
				} else {
					item.Link = resolved
					// s.logger.Debug("resolved redirect", "url", item.Link)
				}
				item.Link = netutil.NormalizeURL(item.Link)
			}(item)
		}
		wg.Wait()
	} else {
		for _, item := range parsedFeed.Items {
			item.Link = netutil.NormalizeURL(item.Link)
		}
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

	domains := []string{baseDomain}
	if feed.RootDomain != nil {
		if rd := strings.TrimSpace(*feed.RootDomain); rd != "" {
			domains = append(domains, rd)
		}
	}

	return domains, nil
}

func (s *Syncer) determineLanguage(feed *database.Feed, parsedFeed *gofeed.Feed) string {
	if l := strings.TrimSpace(parsedFeed.Language); l != "" {
		// from feed
		return strings.ToLower(strings.Split(l, "-")[0])
	}
	if feed.Language != nil && *feed.Language != "" {
		// user override
		return *feed.Language
	}
	// this is a wild guess and will not match!
	s.logger.Error("No language specified in feed or database, defaulting to 'en'", "feed_url", feed.URL, "feed_id", feed.ID)
	return "en"
}

func (s *Syncer) processItems(feed *database.Feed, parsedFeed *gofeed.Feed, items []feeds.ItemHashPair, hashes []string) (int, error) {
	pendingItems, err := s.repo.GetPendingItems(feed.ID, hashes, maxThinkRetries)
	if err != nil {
		return 0, err
	}
	count := len(pendingItems)
	s.logger.Debug("pending items", "len", count)

	if count == 0 {
		s.logger.Debug("All items are already processed")
		return 0, nil
	}

	language := s.determineLanguage(feed, parsedFeed)
	if feed.Language == nil || *feed.Language == "" {
		s.logger.Info("Feed language is not set, updating it.", "feed_id", feed.ID, "language", language)
		feed.Language = &language
		if err := s.repo.UpsertFeed(feed); err != nil {
			// Log the error but continue, as this is not a critical failure
			s.logger.Error("Failed to update feed language", "error", err, "feed_id", feed.ID)
		}
	}

	count = 0
	for _, item := range items {
		if s.ctx.Err() != nil {
			// context might be canceled
			// we have processed x items - so report a 0,nil that
			// don't create a result in postgres - the next tick will pick it up
			return 0, nil
		}
		if errorCount, ok := pendingItems[item.Hash]; ok {
			s.processItem(feed, item.Hash, item.Item, language, errorCount)
			count++
		}
	}

	return count, nil
}

func (s *Syncer) processItem(feed *database.Feed, hash string, item *gofeed.Item, language string, currentErrorCount int) {
	req := think.Request{
		Title:       text.StripHTML(item.Title),
		Description: text.StripHTML(item.Description),
	}
	res, err := s.think.Run(promptScope, language, req)

	var thinkError *string
	var mediaContent *database.MediaContent
	var thinkRating float64
	nextErrorCount := currentErrorCount

	// put it in the item as "deframer:title_original"
	// we technically don't need it here - but helps decoupling the trend miner loop
	// usually we don't persist the xml with the "deframer:..." this is an exception here
	// so that we don't technically need a thinker result
	feeds.SetExtension(item, customPrefix, "title_original", item.Title)
	feeds.SetExtension(item, customPrefix, "description_original", item.Description)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			// worker is shutting down don't store this as an error
			s.logger.Debug("context canceled", "hash", hash)
			return
		}
		s.logger.Error("analysis failed", "error", err, "hash", hash, "item_url", item.Link)
		errStr := err.Error()
		thinkError = &errStr
		nextErrorCount++
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
		if res != nil {
			thinkRating = res.Overall
		}

		applyFancyRatingText(item, res, thinkRating, language)
		// Success reset error count (optional, but good practice if we allow re-processing of succeeded items later)
		nextErrorCount = 0
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
		Hash:            hash,
		FeedID:          feed.ID,
		URL:             netutil.NormalizeURL(item.Link),
		Language:        &language,
		Content:         content,
		PubDate:         pubDate,
		ThinkResult:     res,
		MediaContent:    mediaContent,
		ThinkError:      thinkError,
		ThinkErrorCount: nextErrorCount,
		ThinkRating:     thinkRating,
		Categories:      s.feeds.ExtractCategories(item),
	}

	s.logger.Debug("processItem", "feed", feed.ID, "hash", hash)

	if err := s.repo.UpsertItem(dbItem); err != nil {
		s.logger.Error("failed to create item", "error", err, "hash", hash)
	}

	// move the lock time in the future (we also extend the sleep time to have a fair execution window)
	if err := s.repo.EnqueueSync(feed.ID, config.IdleSleepTime); err != nil {
		s.logger.Error("failed to extend the lock duration item", "error", err, "hash", hash)
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
	if strings.TrimSpace(res.TitleCorrected) != "" {
		item.Title = res.TitleCorrected
	}
	if strings.TrimSpace(res.DescriptionCorrected) != "" {
		item.Description = res.DescriptionCorrected
	}

	// Append the reason to the description
	item.Description = fmt.Sprintf("%s<br/><br/>%s", item.Description, res.OverallReason)

	var mediaData *MediaData

	// If we don't have existing media content, try to extract it
	if len(item.Extensions["media"]["content"]) == 0 {
		if mediaData = s.extractFromContentTag(item); mediaData != nil {
			item.Content = ""
		} else if mediaData = s.extractFromEnclosureTag(item); mediaData != nil {
			// found in enclosure
		} else {
			mediaData = s.extractFromDescriptionFallback(res)
		}
	}

	if mediaData != nil {
		s.applyMediaData(item, mediaData)
	}

	if item.Extensions != nil {
		if mediaExt, ok := item.Extensions["media"]; ok {
			if _, hasGroup := mediaExt["group"]; hasGroup {
				// sometimes media is organized in a media:group
				delete(mediaExt, "group")
				s.logger.Info("removed media:group tag", "url", item.Link)
			}
		}
	}

	return nil
}

func (s *Syncer) extractFromContentTag(item *gofeed.Item) *MediaData {
	if item.Content == "" {
		return nil
	}
	mediaData, err := transformContent(item.Content)
	if err == nil && mediaData.URL != "" {
		return &mediaData
	}
	return nil
}

func (s *Syncer) extractFromEnclosureTag(item *gofeed.Item) *MediaData {
	for _, enc := range item.Enclosures {
		if strings.HasPrefix(enc.Type, "image/") && enc.URL != "" {
			w, h := parseDimensions(enc.URL)
			return &MediaData{
				URL:    enc.URL,
				Medium: "image",
				Width:  w,
				Height: h,
			}
		}
	}
	return nil
}

func (s *Syncer) extractFromDescriptionFallback(res *database.ThinkResult) *MediaData {
	target := res.DescriptionOriginal
	mediaData, err := transformContent(target)
	if err != nil || mediaData.URL == "" {
		target = res.TitleOriginal
		mediaData, err = transformContent(target)
	}

	if err == nil && mediaData.URL != "" {
		return &mediaData
	}
	return nil
}

func (s *Syncer) applyMediaData(item *gofeed.Item, data *MediaData) {
	if item.Extensions == nil {
		item.Extensions = make(map[string]map[string][]ext.Extension)
	}
	if item.Extensions["media"] == nil {
		item.Extensions["media"] = make(map[string][]ext.Extension)
	}

	// <media:credit> always delete this
	delete(item.Extensions["media"], "credit")

	attrs := map[string]string{
		"url":    data.URL,
		"medium": data.Medium,
	}
	if data.Width > 0 {
		attrs["width"] = fmt.Sprintf("%d", data.Width)
	}
	if data.Height > 0 {
		attrs["height"] = fmt.Sprintf("%d", data.Height)
	}

	item.Extensions["media"]["content"] = []ext.Extension{{
		Name:  "content",
		Attrs: attrs,
	}}

	if len(data.Alt) > 0 {
		item.Extensions["media"]["description"] = []ext.Extension{{
			Name:  "description",
			Value: data.Alt,
		}}
	}
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
		// Fallback: check for media:thumbnail
		if thumbs, ok := mediaExt["thumbnail"]; ok && len(thumbs) > 0 {
			for _, t := range thumbs {
				if t.Attrs["url"] != "" {
					val := t
					// Deep copy attrs to avoid side effects and allow modification
					newAttrs := make(map[string]string)
					for k, v := range val.Attrs {
						newAttrs[k] = v
					}
					val.Attrs = newAttrs

					if _, ok := val.Attrs["medium"]; !ok {
						val.Attrs["medium"] = "image"
					}
					contentExt = &val
					break
				}
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
