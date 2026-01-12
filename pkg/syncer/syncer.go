package syncer

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/feeds"
	"github.com/egandro/news-deframer/pkg/think"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
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
	th, err := think.New(cfg)
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
		if pendingHashes[item.Hash] {
			// // A11Y - foo <span role="img" aria-label="Rating: 2 out of 4">◑</span>

			// 	Filled: ❶ ❷ ❸ ❹ ❺ (U+2776 to U+277A)
			// 	Outline: ① ② ③ ④ ⑤ (U+2460 to U+2464)
			// 	Negative Squares: 1️⃣ 2️⃣ 3️⃣ (Emoji style)

			// 	0/4: ○ (U+25CB)
			// 	1/4: ◔ (U+25D4)
			// 	2/4: ◑ (U+25D1)
			// 	3/4: ◕ (U+25D5)
			// 	4/4: ● (U+25CF)[2]
			// item.Item.Title = "★★★☆☆" + item.Item.Title
			// item.Item.Description = "★☆☆☆☆" + item.Item.Description
			s.processItem(feed, item.Hash, item.Item, language)
			count++
		}
	}

	return count, nil
}

func (s *Syncer) processItem(feed *database.Feed, hash string, item *gofeed.Item, language string) {
	data := map[string]interface{}{
		"title":       item.Title,
		"description": item.Description,
	}
	res, err := s.think.Run(promptScope, language, data)
	if err != nil {
		s.logger.Error("analysis failed", "error", err, "hash", hash, "item_url", item.Link)
		res = map[string]interface{}{
			"error": err.Error(),
		}
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

	analyzerResult := database.JSONB(res)
	dbItem := &database.Item{
		Hash:           hash,
		FeedID:         feed.ID,
		URL:            item.Link,
		AnalyzerResult: &analyzerResult,
		Content:        content,
		PubDate:        pubDate,
	}

	s.logger.Debug("processItem", "feed", feed.ID, "hash", hash)

	if err := s.repo.UpsertItem(dbItem); err != nil {
		s.logger.Error("failed to create item", "error", err, "hash", hash)
	}
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

		if item.AnalyzerResult != nil {
			res := *item.AnalyzerResult
			isErrorOnly := false
			if len(res) == 1 {
				_, isErrorOnly = res["error"]
			}

			if len(res) > 0 && !isErrorOnly {
				for k, v := range res {
					feeds.SetExtension(current, customPrefix, k, fmt.Sprintf("%v", v))
				}
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
