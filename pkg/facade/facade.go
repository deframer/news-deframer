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
)

type RSSProxyFilter struct {
	URL      string
	Lang     string
	MaxScore float64
	Embedded bool
}

type ItemResult struct {
	FeedURL  string `json:"feed_url"`
	URL      string `json:"url"`
	Hash     string `json:"hash"`
	AIResult string `json:"ai_result"`
}

type Facade interface {
	GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error)
	GetItems(ctx context.Context, u *url.URL) ([]ItemResult, error)
}

type facade struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
	repo   database.Repository
}

func New(ctx context.Context, cfg *config.Config, repo database.Repository, dl ...downloader.Downloader) Facade {
	return &facade{
		ctx:    ctx,
		cfg:    cfg,
		logger: slog.With("component", "facade"),
		repo:   repo,
	}
}

func (f *facade) GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error) {
	if filter == nil {
		return "", fmt.Errorf("filter is nil")
	}

	u, err := url.Parse(filter.URL)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	feed, err := f.repo.FindFeedByUrl(u)
	if err != nil {
		return "", err
	}

	if feed == nil {
		return "", fmt.Errorf("unknown feed")
	}

	cachedFeed, err := f.repo.FindCachedFeedById(feed.ID)
	if err != nil {
		return "", err
	}

	if cachedFeed == nil {
		feedSchedule, err := f.repo.FindFeedScheduleById(feed.ID)
		if err != nil {
			return "", err
		}
		if feedSchedule != nil && feedSchedule.LastError != nil {
			return "", fmt.Errorf("last sync attempt at %s failed with error: %s", feedSchedule.UpdatedAt.Format(time.RFC1123), *feedSchedule.LastError)
		}
		return "", fmt.Errorf("feed not updated yet - come back later")
	}

	return *cachedFeed.XMLContent, nil
}

func (f *facade) GetItems(ctx context.Context, u *url.URL) ([]ItemResult, error) {
	var results []ItemResult
	return results, nil
}
