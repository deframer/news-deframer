package facade

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/feeds"
	"github.com/egandro/news-deframer/pkg/valkey"
	"golang.org/x/net/publicsuffix"
)

var (
	maxPendingTimeout = 30 * time.Second
	checkInterval     = 5 * time.Second
	maxTimeout        = 5 * time.Minute
)

type Facade interface {
	HasFeed(ctx context.Context, u *url.URL) (bool, error)
	GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error)
	GetItems(ctx context.Context, u *url.URL) ([]ItemResult, error)
}

type facade struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
	valkey valkey.Valkey
	repo   database.Repository
	dl     downloader.Downloader
	feeds  feeds.Feeds
}

func New(ctx context.Context, cfg *config.Config, v valkey.Valkey, repo database.Repository, dl ...downloader.Downloader) Facade {
	var d downloader.Downloader
	if len(dl) > 0 {
		d = dl[0]
	} else {
		d = downloader.NewDownloader(ctx, cfg)
	}

	return &facade{
		ctx:    ctx,
		cfg:    cfg,
		logger: slog.With("component", "facade"),
		valkey: v,
		repo:   repo,
		dl:     d,
		feeds:  feeds.NewFeeds(ctx, cfg),
	}
}

func (f *facade) HasFeed(ctx context.Context, u *url.URL) (bool, error) {
	val, err := f.valkey.GetFeedByUrl(u)
	if err != nil {
		return false, err
	}

	// If we have a definitive result, return it. (= cache hit)
	if val != nil {
		switch val.Cache {
		case valkey.Ok:
			return true, nil
		case valkey.ValueUnknown:
			return false, nil
		}
	}

	// If it was missing (nil), try to acquire the lock to be the one fetching it.
	if val == nil {
		acquired, err := f.valkey.TryLockFeedByUrl(u, valkey.FeedUrlToUUID{Cache: valkey.Updating}, maxPendingTimeout)
		if err != nil {
			return false, err
		}
		if !acquired {
			// We failed to acquire lock, meaning someone else just set it to pending.
			// Fall through to the wait loop.
			goto WAIT_LOOP
		}

		// We acquired the lock. Proceed to DB lookup.
		return f.fetchAndCache(u)
	}

	// It is pending (either from GetFeedUUID or failed UpdateFeedUUID). Wait for it.
WAIT_LOOP:
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	timeout := time.After(maxPendingTimeout)

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for pending feed")
		case <-ticker.C:
			val, err := f.valkey.GetFeedByUrl(u)
			if err != nil {
				return false, err
			}
			if val != nil {
				switch val.Cache {
				case valkey.Ok:
					return true, nil
				case valkey.ValueUnknown:
					return false, nil
				}
			}
			// If val is nil here, it means the pending key expired or was deleted without result.
			// we don't care and wait for the next query
		}
	}
}

func (f *facade) fetchAndCache(u *url.URL) (bool, error) {
	feed, err := f.repo.FindFeedByUrl(u)
	if err != nil {
		return false, err
	}

	if feed == nil {
		return false, nil
	}

	state := valkey.FeedUrlToUUID{Cache: valkey.ValueUnknown}
	var info valkey.FeedInfo
	info.URL = u.String()
	if feed != nil {
		var baseDomains []string
		if feed.EnforceFeedDomain {
			// we need to enforce the base domain
			if bd, err := publicsuffix.EffectiveTLDPlusOne(u.Hostname()); err == nil {
				baseDomains = append(baseDomains, bd)
			} else {
				baseDomains = append(baseDomains, u.Hostname())
			}
			// TODO we might have an allow list table later for each feed
		}
		state = valkey.FeedUrlToUUID{
			Cache: valkey.Ok,
			UUID:  feed.ID,
		}
		info.BaseDomain = baseDomains
	}

	if err := f.valkey.UpdateFeedByUrl(state, info, maxTimeout); err != nil {
		return false, err
	}
	return state.Cache != valkey.ValueUnknown, nil
}

type RSSProxyFilter struct {
	URL      string
	Lang     string
	MaxScore float64
	Embedded bool
}

func (f *facade) GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error) {
	if filter == nil {
		filter = &RSSProxyFilter{}
	}

	u, err := url.Parse(filter.URL)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	f.logger.DebugContext(ctx, "downloading feed", "url", filter.URL)

	rc, err := f.dl.DownloadRSSFeed(ctx, u)
	if err != nil {
		return "", err
	}
	defer func() { _ = rc.Close() }()

	feed, err := f.feeds.ParseFeed(ctx, rc)
	if err != nil {
		return "", err
	}

	feed.Title = "Proxied: " + feed.Title

	newFeed, items := f.feeds.GetValidItems(f.ctx, feed)
	f.logger.Debug("items", "len", len(items))

	// timeStr := time.Now().Format(time.RFC1123Z)

	// dummy := &gofeed.Item{
	// 	Title:       "News Deframer Proxy",
	// 	Description: fmt.Sprintf("Proxied feed for: %s", filter.URL),
	// 	Link:        "https://github.com/egandro/news-deframer",
	// 	Published:   timeStr,
	// 	GUID:        "news-deframer-status: " + timeStr,
	// }
	// newFeed.Items = append([]*gofeed.Item{dummy}, newFeed.Items...)

	return f.feeds.RenderFeed(ctx, newFeed)
}

type ItemResult struct {
	FeedURL  string `json:"feed_url"`
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	AIResult string `json:"ai_result"`
}

func (f *facade) GetItems(ctx context.Context, u *url.URL) ([]ItemResult, error) {
	var results []ItemResult
	return results, nil
}
