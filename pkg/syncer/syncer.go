package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
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
	"golang.org/x/net/publicsuffix"
)

const feedTitlePrefix = "Deframer: "
const promptScope = "deframer"
const customPrefix = "deframer"
const customNamespace = "https://github.com/deframer/news-deframer/"
const maxThinkRetries = 3
const thinkFixerBatchSize = 15
const thinkFixerLookback = 30 * 24 * time.Hour
const thinkFixerStopThreshold = maxThinkRetries
const publicationDateGracePeriod = 10 * time.Minute

type Mode string

const (
	ModeWorker     Mode = "worker"
	ModeThinkFixer Mode = "think-fixer"
)

type FeedSyncer interface {
	SyncFeed(id uuid.UUID) error
	StopPolling(id uuid.UUID) error
	Poll(mode Mode)
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

func (s *Syncer) Poll(mode Mode) {
	s.logger.Info("Starting poller", "mode", mode)
	if mode == ModeThinkFixer {
		s.pollThinkerFixerMode()
		return
	}
	if mode != ModeWorker {
		s.logger.Warn("Unknown mode, defaulting to worker", "mode", mode)
	}
	s.pollWorkerMode()
}

func (s *Syncer) pollWorkerMode() {
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

func (s *Syncer) pollThinkerFixerMode() {
	for {
		if s.ctx.Err() != nil {
			s.logger.Info("Stopping poller")
			return
		}

		if s.fixNextThinkerBatch() {
			s.logger.Info("A thinker-fixer batch was rechecked")
			continue
		}

		s.logger.Info("Think-fixer sleeping...", "duration", config.IdleSleepTime)

		select {
		case <-s.ctx.Done():
			s.logger.Info("Stopping poller")
			return
		case <-time.After(config.IdleSleepTime):
		}
	}
}

// fixNextThinkerBatch return true if this has updated entries
func (s *Syncer) fixNextThinkerBatch() bool {
	s.logger.Info("fixNextThinkerBatch")
	since := time.Now().Add(-thinkFixerLookback)
	maxErrorCount := thinkFixerStopThreshold * 2
	items, err := s.repo.BeginThinkFixerBatch(thinkFixerBatchSize, since, thinkFixerStopThreshold, maxErrorCount, config.DefaultLockDuration)
	if err != nil {
		s.logger.Error("Failed to query thinker-fixer candidates", "error", err)
		return false
	}
	if len(items) == 0 {
		return false
	}

	s.logger.Info("Think-fixer candidates fetched", "count", len(items))
	for i := range items {
		current := i + 1
		s.logger.Debug("processThinkerItem", "item_id", items[i].ID, "feed_id", items[i].FeedID, "progress", fmt.Sprintf("%d/%d", current, len(items)))
		s.processThinkerItem(&items[i])
	}
	return true
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

	for _, item := range parsedFeed.Items {
		item.Link = netutil.NormalizeURL(item.Link)
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

	if feed.ResolveItemUrl {
		s.logger.Debug("resolving remaining item urls", "url", feed.URL)
		var wg sync.WaitGroup
		// Limit concurrency to avoid overwhelming the server or local resources
		sem := make(chan struct{}, 10)

		for _, item := range items {
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
			}(item.Item)
		}
		wg.Wait()
	}

	total := len(pendingItems)
	count = 0
	for _, item := range items {
		if s.ctx.Err() != nil {
			// context might be canceled
			// we have processed x items - so report a 0,nil that
			// don't create a result in feed_schedules table - the next tick will pick it up
			return 0, nil
		}
		if errorCount, ok := pendingItems[item.Hash]; ok {
			count++
			s.logger.Debug("processItem", "feed", feed.ID, "hash", item.Hash, "progress", fmt.Sprintf("%d/%d", count, total))
			s.processItem(feed, item.Hash, item.Item, language, errorCount)
		}
	}

	return count, nil
}

func (s *Syncer) processItem(feed *database.Feed, hash string, item *gofeed.Item, language string, currentErrorCount int) {
	result, err := s.thinkRenderAndExtract(item, language, currentErrorCount, "hash", hash, "item_url", item.Link)
	if err != nil {
		return
	}

	pubDate := time.Now()
	if result.pubDate != nil {
		if result.pubDate.After(time.Now().Add(publicationDateGracePeriod)) {
			s.logger.Warn("ignoring future publication date", "hash", hash, "item_url", item.Link, "pub_date", *result.pubDate)
		} else {
			pubDate = *result.pubDate
		}
	}

	dbItem := &database.Item{
		Hash:            hash,
		FeedID:          feed.ID,
		URL:             netutil.NormalizeURL(item.Link),
		Language:        &language,
		Content:         result.content,
		PubDate:         pubDate,
		ThinkResult:     result.thinkResult,
		MediaContent:    result.mediaContent,
		ThinkError:      result.thinkError,
		ThinkErrorCount: result.nextErrorCount,
		ThinkRating:     result.thinkRating,
		Categories:      emptyStringArray(result.categories),
		Authors:         emptyStringArray(result.authors),
	}

	if err := s.repo.UpsertItem(dbItem); err != nil {
		s.logger.Error("failed to create item", "error", err, "hash", hash)
	}

	// move the lock time in the future (we also extend the sleep time to have a fair execution window)
	if err := s.repo.EnqueueSync(feed.ID, config.IdleSleepTime); err != nil {
		s.logger.Error("failed to extend the lock duration item", "error", err, "hash", hash)
	}
}

func (s *Syncer) processThinkerItem(dbItem *database.Item) {
	if dbItem == nil {
		return
	}

	parsedItem, err := s.parseItemFromContent(dbItem.Content)
	if err != nil {
		s.logger.Error("Failed to parse item content", "error", err, "item_id", dbItem.ID)
		return
	}

	language := "en"
	if dbItem.Language != nil && *dbItem.Language != "" {
		language = *dbItem.Language
	}

	result, err := s.thinkRenderAndExtract(parsedItem, language, dbItem.ThinkErrorCount, "item_id", dbItem.ID, "item_url", parsedItem.Link)
	if err != nil {
		return
	}

	if result.pubDate != nil {
		if result.pubDate.After(time.Now().Add(publicationDateGracePeriod)) {
			s.logger.Warn("ignoring future publication date", "item_id", dbItem.ID, "item_url", parsedItem.Link, "pub_date", *result.pubDate)
		} else {
			dbItem.PubDate = *result.pubDate
		}
	}

	dbItem.Content = result.content
	dbItem.ThinkResult = result.thinkResult
	dbItem.MediaContent = result.mediaContent
	dbItem.ThinkError = result.thinkError
	dbItem.ThinkErrorCount = result.nextErrorCount
	dbItem.ThinkRating = result.thinkRating
	dbItem.Categories = emptyStringArray(result.categories)
	dbItem.Authors = emptyStringArray(result.authors)

	// let the trend miner recreate it as it now has access to the thinker results
	if err := s.repo.UpsertItemWithTrendInvalidation(dbItem); err != nil {
		s.logger.Error("failed to update item", "error", err, "item_id", dbItem.ID)
	}
}

type thinkerOutcome struct {
	content        string
	thinkResult    *database.ThinkResult
	mediaContent   *database.MediaContent
	thinkError     *string
	nextErrorCount int
	thinkRating    float64
	categories     []string
	authors        database.StringArray
	pubDate        *time.Time
}

func emptyStringArray(values []string) database.StringArray {
	if values == nil {
		return database.StringArray{}
	}
	return database.StringArray(values)
}

func (s *Syncer) thinkRenderAndExtract(parsedItem *gofeed.Item, language string, currentErrorCount int, logKeys ...any) (*thinkerOutcome, error) {
	if parsedItem == nil {
		return nil, fmt.Errorf("parsed item is nil")
	}

	req := think.Request{
		Title:       text.StripHTML(parsedItem.Title),
		Description: text.StripHTML(parsedItem.Description),
	}
	res, err := s.think.Run(promptScope, language, req)

	var thinkError *string
	var mediaContent *database.MediaContent
	var thinkRating float64
	nextErrorCount := currentErrorCount

	feeds.SetExtension(parsedItem, customPrefix, "title_original", parsedItem.Title)
	feeds.SetExtension(parsedItem, customPrefix, "description_original", parsedItem.Description)

	if err != nil {
		if errors.Is(err, context.Canceled) || s.ctx.Err() != nil {
			s.logger.Debug("context canceled", logKeys...)
			return nil, err
		}
		args := append([]any{"error", err}, logKeys...)
		s.logger.Error("analysis failed", args...)
		errStr := err.Error()
		thinkError = &errStr
		nextErrorCount++
	} else {
		err = s.updateContent(parsedItem, res)
		if err != nil {
			args := append([]any{"error", err}, logKeys...)
			s.logger.Error("update content failed", args...)
			return nil, err
		}
		mediaContent, err = extractMediaContent(parsedItem)
		if err != nil {
			args := append([]any{"error", err}, logKeys...)
			s.logger.Error("extract media content failed", args...)
		}
		if res != nil {
			thinkRating = res.Overall
		}

		applyFancyRatingText(parsedItem, res, thinkRating, language)
		nextErrorCount = 0
	}

	content, err := s.feeds.RenderItem(s.ctx, parsedItem)
	if err != nil {
		args := append([]any{"error", err}, logKeys...)
		s.logger.Error("failed to render item", args...)
		return nil, err
	}

	var pubDate *time.Time
	if parsedItem.PublishedParsed != nil {
		pubDate = parsedItem.PublishedParsed
	}

	return &thinkerOutcome{
		content:        content,
		thinkResult:    res,
		mediaContent:   mediaContent,
		thinkError:     thinkError,
		nextErrorCount: nextErrorCount,
		thinkRating:    thinkRating,
		categories:     s.feeds.ExtractCategories(parsedItem),
		authors:        s.extractAndNormalizeAuthors(parsedItem, language),
		pubDate:        pubDate,
	}, nil
}

func (s *Syncer) parseItemFromContent(content string) (*gofeed.Item, error) {
	wrapped := `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel>` + content + `</channel></rss>`
	pf, err := gofeed.NewParser().Parse(strings.NewReader(wrapped))
	if err != nil {
		return nil, err
	}
	if len(pf.Items) == 0 {
		return nil, fmt.Errorf("no items found in content")
	}
	return pf.Items[0], nil
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
		if mediaData = extractFromContentTag(item); mediaData != nil {
			item.Content = ""
		} else if mediaData = extractFromEnclosureTag(item); mediaData != nil {
			// found in enclosure
		} else {
			mediaData = extractFromDescriptionFallback(res)
		}
	}

	if mediaData != nil {
		applyMediaData(item, mediaData)
	}

	if item.Extensions != nil {
		if mediaExt, ok := item.Extensions["media"]; ok {
			// sometimes media is organized in a media:group
			delete(mediaExt, "group")
		}
	}

	return nil
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
