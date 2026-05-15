package syncer

import (
	"context"
	"errors"
	"fmt"
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
	"goa.design/clue/log"
	"golang.org/x/net/publicsuffix"
)

const promptScope = "deframer"
const maxThinkRetries = 3
const thinkerBatchSize = 15
const thinkerFixerLookback = 90 * 24 * time.Hour
const thinkerFixerMinErrorCount = 4
const thinkerFixerMaxErrorCount = 6
const publicationDateGracePeriod = 10 * time.Minute

type Mode string

const (
	ModeIngester     Mode = "ingester"
	ModeThinker      Mode = "thinker"
	ModeThinkerFixer Mode = "thinker-fixer"
)

type FeedSyncer interface {
	SyncFeed(id uuid.UUID) error
	StopPolling(id uuid.UUID) error
	Poll(mode Mode)
}

type Syncer struct {
	ctx   context.Context
	repo  database.Repository
	dl    downloader.Downloader
	feeds feeds.Feeds
	think think.Think
}

func New(ctx context.Context, cfg *config.Config, repo database.Repository) (*Syncer, error) {
	th, err := think.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Syncer{
		ctx:   ctx,
		repo:  repo,
		dl:    downloader.NewDownloader(ctx, cfg),
		feeds: feeds.NewFeeds(ctx, cfg),
		think: th,
	}, nil
}

func (s *Syncer) SyncFeed(id uuid.UUID) error {
	log.Printf(s.ctx, "Syncing feed id=%s", id)
	return s.repo.EnqueueSync(id, 0)
}

func (s *Syncer) StopPolling(id uuid.UUID) error {
	log.Printf(s.ctx, "Stopping polling id=%s", id)
	return s.repo.RemoveSync(id)
}

func (s *Syncer) Poll(mode Mode) {
	log.Printf(s.ctx, "Starting poller mode=%s", mode)
	if mode == ModeThinker {
		s.pollThinker()
		return
	}
	if mode == ModeThinkerFixer {
		s.pollThinkerFixer(thinkerFixerLookback)
		return
	}
	if mode != ModeIngester {
		log.Warnf(s.ctx, "Unknown mode, defaulting to ingester mode=%s", mode)
	}
	s.pollIngester()
}

func (s *Syncer) pollIngester() {
	for {
		if s.ctx.Err() != nil {
			log.Printf(s.ctx, "Stopping poller")
			return
		}

		if s.syncNextScheduledFeed() {
			log.Printf(s.ctx, "A feed was synced")
			continue
		}

		log.Debugf(s.ctx, "Sleeping duration=%s", config.IdleSleepTime)

		select {
		case <-s.ctx.Done():
			log.Printf(s.ctx, "Stopping poller")
			return
		case <-time.After(config.IdleSleepTime):
		}
	}
}

func (s *Syncer) pollThinker() {
	for {
		if s.ctx.Err() != nil {
			log.Printf(s.ctx, "Stopping poller")
			return
		}

		if s.processThinkerBatch() {
			log.Printf(s.ctx, "Thinker batch checked")
			continue
		}

		log.Printf(s.ctx, "Thinker sleep duration=%s", config.IdleSleepTime)

		select {
		case <-s.ctx.Done():
			log.Printf(s.ctx, "Stopping poller")
			return
		case <-time.After(config.IdleSleepTime):
		}
	}
}

func (s *Syncer) pollThinkerFixer(lookback time.Duration) {
	for {
		if s.ctx.Err() != nil {
			log.Printf(s.ctx, "Stopping poller")
			return
		}

		if s.processThinkerFixerBatch(lookback) {
			log.Printf(s.ctx, "Thinker fixer batch checked")
			continue
		}

		log.Debugf(s.ctx, "Thinker fixer sleep duration=%s", config.IdleSleepTime)

		select {
		case <-s.ctx.Done():
			log.Printf(s.ctx, "Stopping poller")
			return
		case <-time.After(config.IdleSleepTime):
		}
	}
}

func (s *Syncer) processThinkerBatch() bool {
	log.Printf(s.ctx, "processThinkerBatch")
	items, err := s.repo.BeginThinkerBatch(thinkerBatchSize, time.Time{}, 0, maxThinkRetries, config.DefaultLockDuration)
	if err != nil {
		log.Errorf(s.ctx, err, "Failed to query thinker candidates")
		return false
	}
	if len(items) == 0 {
		return false
	}

	log.Printf(s.ctx, "Thinker candidates fetched count=%d", len(items))
	for i := range items {
		current := i + 1
		log.Debugf(s.ctx, "processThinkerItem item_id=%s feed_id=%s progress=%d/%d", items[i].ID, items[i].FeedID, current, len(items))
		s.thinkItem(&items[i])
	}
	return true
}

func (s *Syncer) processThinkerFixerBatch(lookback time.Duration) bool {
	log.Debugf(s.ctx, "processThinkerFixerBatch")
	since := time.Time{}
	if lookback > 0 {
		since = time.Now().Add(-lookback)
	}
	items, err := s.repo.BeginThinkerFixerBatch(thinkerBatchSize, since, thinkerFixerMinErrorCount, thinkerFixerMaxErrorCount, config.DefaultLockDuration)
	if err != nil {
		log.Errorf(s.ctx, err, "Failed to query thinker fixer candidates")
		return false
	}
	if len(items) == 0 {
		return false
	}

	log.Printf(s.ctx, "Thinker fixer candidates fetched count=%d", len(items))
	for i := range items {
		current := i + 1
		log.Debugf(s.ctx, "processThinkerItem item_id=%s feed_id=%s progress=%d/%d", items[i].ID, items[i].FeedID, current, len(items))
		s.thinkItem(&items[i])
	}
	return true
}

// syncNextScheduledFeed return true if this has updated entries
func (s *Syncer) syncNextScheduledFeed() bool {
	log.Debugf(s.ctx, "syncNextScheduledFeed")
	feed, err := s.repo.BeginFeedUpdate(config.DefaultLockDuration)
	if err != nil {
		log.Errorf(s.ctx, err, "Failed query the next feed to sync")
		return false
	}
	if feed == nil {
		return false
	}
	err = s.updatingFeed(feed)
	if err := s.repo.EndFeedUpdate(feed.ID, err, config.PollingInterval); err != nil {
		log.Errorf(s.ctx, err, "Failed to end feed update")
	}
	return true
}

func (s *Syncer) updatingFeed(feed *database.Feed) error {
	log.Printf(s.ctx, "Updating feed id=%s url=%s", feed.ID, feed.URL)

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
	log.Debugf(s.ctx, "items len=%d", len(items))

	hashes := feeds.GetHashes(items)
	log.Debugf(s.ctx, "hashes len=%d", len(hashes))

	count, err := s.syncPendingFeedItems(feed, parsedFeed, items, hashes)
	if err != nil {
		return err
	}

	if count == 0 {
		log.Debugf(s.ctx, "all items are processed - not updating the feed")
		return nil
	}

	return nil
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
	log.Errorf(s.ctx, fmt.Errorf("no language specified in feed or database"), "No language specified in feed or database, defaulting to 'en' feed_url=%s feed_id=%s", feed.URL, feed.ID)
	return "en"
}

func (s *Syncer) syncPendingFeedItems(feed *database.Feed, parsedFeed *gofeed.Feed, items []feeds.ItemHashPair, hashes []string) (int, error) {
	pendingItems, err := s.repo.GetPendingItems(feed.ID, hashes, maxThinkRetries)
	if err != nil {
		return 0, err
	}
	count := len(pendingItems)
	log.Debugf(s.ctx, "pending items len=%d", count)

	if count == 0 {
		log.Debugf(s.ctx, "All items are already processed")
		return 0, nil
	}

	language := s.determineLanguage(feed, parsedFeed)
	if feed.Language == nil || *feed.Language == "" {
		log.Printf(s.ctx, "Feed language is not set, updating it. feed_id=%s language=%s", feed.ID, language)
		feed.Language = &language
		if err := s.repo.UpsertFeed(feed); err != nil {
			// Log the error but continue, as this is not a critical failure
			log.Errorf(s.ctx, err, "Failed to update feed language feed_id=%s", feed.ID)
		}
	}

	if feed.ResolveItemUrl {
		log.Debugf(s.ctx, "resolving remaining item urls url=%s", feed.URL)
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
					log.Warnf(s.ctx, "failed to resolve redirect url=%s error=%v", item.Link, err)
				} else {
					item.Link = resolved
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
		if _, ok := pendingItems[item.Hash]; ok {
			count++
			log.Debugf(s.ctx, "syncItem feed=%s hash=%s progress=%d/%d", feed.ID, item.Hash, count, total)
			s.syncItem(feed, item.Hash, item.Item, language)
		}
	}

	return count, nil
}

func (s *Syncer) syncItem(feed *database.Feed, hash string, item *gofeed.Item, language string) {
	pubDate := time.Now()
	if item.PublishedParsed != nil {
		if item.PublishedParsed.After(time.Now().Add(publicationDateGracePeriod)) {
			log.Warnf(s.ctx, "ignoring future publication date hash=%s item_url=%s pub_date=%s", hash, item.Link, item.PublishedParsed)
		} else {
			pubDate = *item.PublishedParsed
		}
	}

	content, err := s.feeds.RenderItem(s.ctx, item)
	if err != nil {
		log.Errorf(s.ctx, err, "failed to render item hash=%s", hash)
		return
	}

	mediaContent, err := extractMediaContent(item)
	if err != nil {
		log.Errorf(s.ctx, err, "failed to extract media content hash=%s", hash)
	}

	dbItem := &database.Item{
		Hash:            hash,
		FeedID:          feed.ID,
		URL:             netutil.NormalizeURL(item.Link),
		Language:        &language,
		Content:         content,
		PubDate:         pubDate,
		MediaContent:    mediaContent,
		ThinkErrorCount: 0,
		Categories:      emptyStringArray(s.feeds.ExtractCategories(item)),
		Authors:         emptyStringArray(s.extractAndNormalizeAuthors(item, language)),
	}

	if err := s.repo.UpsertItem(dbItem); err != nil {
		log.Errorf(s.ctx, err, "failed to create item hash=%s", hash)
	}

	// move the lock time in the future (we also extend the sleep time to have a fair execution window)
	if err := s.repo.EnqueueSync(feed.ID, config.IdleSleepTime); err != nil {
		log.Errorf(s.ctx, err, "failed to extend the lock duration item hash=%s", hash)
	}
}

func (s *Syncer) thinkItem(dbItem *database.Item) {
	if dbItem == nil {
		return
	}

	parsedItem, err := s.parseItemContent(dbItem.Content)
	if err != nil {
		log.Errorf(s.ctx, err, "Failed to parse item content item_id=%s", dbItem.ID)
		return
	}

	language := "en"
	if dbItem.Language != nil && *dbItem.Language != "" {
		language = *dbItem.Language
	}

	result, err := s.renderThoughtsAndItem(parsedItem, language, dbItem.ThinkErrorCount, "item_id", dbItem.ID, "item_url", parsedItem.Link)
	if err != nil {
		return
	}

	if result.pubDate != nil {
		if result.pubDate.After(time.Now().Add(publicationDateGracePeriod)) {
			log.Warnf(s.ctx, "ignoring future publication date item_id=%s item_url=%s pub_date=%s", dbItem.ID, parsedItem.Link, result.pubDate)
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
		log.Errorf(s.ctx, err, "failed to update item item_id=%s", dbItem.ID)
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

func (s *Syncer) renderThoughtsAndItem(parsedItem *gofeed.Item, language string, currentErrorCount int, logKeys ...any) (*thinkerOutcome, error) {
	if parsedItem == nil {
		return nil, fmt.Errorf("parsed item is nil")
	}

	req := think.Request{
		Title:       text.StripHTML(parsedItem.Title),
		Description: text.StripHTML(parsedItem.Description),
	}
	futureErrorCount := currentErrorCount + 1
	ignoreCategoryErrors := futureErrorCount > maxThinkRetries || futureErrorCount > thinkerFixerMaxErrorCount
	res, err := s.think.Run(promptScope, language, req, ignoreCategoryErrors)

	var thinkError *string
	var mediaContent *database.MediaContent
	var thinkRating float64
	nextErrorCount := currentErrorCount

	if err != nil {
		if errors.Is(err, context.Canceled) || s.ctx.Err() != nil {
			log.Debugf(s.ctx, "%s", formatLogKeys("context canceled", logKeys...))
			return nil, err
		}
		log.Errorf(s.ctx, err, "%s", formatLogKeys("analysis failed", logKeys...))
		errStr := err.Error()
		thinkError = &errStr
		nextErrorCount++
	} else {
		err = updateContent(parsedItem, res)
		if err != nil {
			log.Errorf(s.ctx, err, "%s", formatLogKeys("update content failed", logKeys...))
			return nil, err
		}
		mediaContent, err = extractMediaContent(parsedItem)
		if err != nil {
			log.Errorf(s.ctx, err, "%s", formatLogKeys("extract media content failed", logKeys...))
		}
		if res != nil {
			thinkRating = res.Overall
		}

		applyFancyRatingText(parsedItem, res, thinkRating, language)
		nextErrorCount = 0
	}

	content, err := s.feeds.RenderItem(s.ctx, parsedItem)
	if err != nil {
		log.Errorf(s.ctx, err, "%s", formatLogKeys("failed to render item", logKeys...))
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

func (s *Syncer) parseItemContent(content string) (*gofeed.Item, error) {
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

func formatLogKeys(msg string, logKeys ...any) string {
	if len(logKeys) == 0 {
		return msg
	}

	var b strings.Builder
	b.WriteString(msg)
	for i := 0; i < len(logKeys); i += 2 {
		b.WriteByte(' ')
		fmt.Fprint(&b, logKeys[i])
		b.WriteByte('=')
		if i+1 < len(logKeys) {
			fmt.Fprint(&b, logKeys[i+1])
		} else {
			b.WriteString("<missing>")
		}
	}
	return b.String()
}
