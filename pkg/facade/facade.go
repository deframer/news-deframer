package facade

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/downloader"
)

type RSSProxyFilter struct {
	URL      string
	Lang     string
	Max      float64
	Embedded bool
}

const MaxItemsForRootDomain = 30

type AnalyzedItem struct {
	Hash string `json:"hash"`
	URL  string `json:"url"`
	database.ThinkResult
	MediaContent *database.MediaContent `json:"media,omitempty"`
	ThinkRating  float64                `json:"rating"`
	PubDate      time.Time              `json:"pubDate"`
}

type Facade interface {
	GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error)
	GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]AnalyzedItem, error)
	GetFirstItemForUrl(ctx context.Context, u *url.URL) (*AnalyzedItem, error)
	GetRootDomains(ctx context.Context) ([]string, error)
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
		if feedSchedule == nil {
			return "", fmt.Errorf("feed not updated yet - come back later")
		}
	}

	return *cachedFeed.XMLContent, nil
}

func (f *facade) GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]AnalyzedItem, error) {
	dbItems, err := f.repo.FindItemsByRootDomain(rootDomain, MaxItemsForRootDomain)
	if err != nil {
		return nil, err
	}

	items := make([]AnalyzedItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		if maxScore > 0 && dbItem.ThinkRating > maxScore {
			continue
		}

		item := AnalyzedItem{
			Hash:         dbItem.Hash,
			URL:          dbItem.URL,
			MediaContent: dbItem.MediaContent,
			ThinkRating:  dbItem.ThinkRating,
			PubDate:      dbItem.PubDate,
		}
		if dbItem.ThinkResult != nil {
			item.ThinkResult = *dbItem.ThinkResult
		}
		items = append(items, item)
	}
	return items, nil
}

func (f *facade) GetFirstItemForUrl(ctx context.Context, u *url.URL) (*AnalyzedItem, error) {
	items, err := f.repo.FindItemsByUrl(u)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	dbItem := items[0]
	item := AnalyzedItem{
		Hash:         dbItem.Hash,
		URL:          dbItem.URL,
		MediaContent: dbItem.MediaContent,
		ThinkRating:  dbItem.ThinkRating,
		PubDate:      dbItem.PubDate,
	}
	if dbItem.ThinkResult != nil {
		item.ThinkResult = *dbItem.ThinkResult
	}

	return &item, nil
}

func (f *facade) GetRootDomains(ctx context.Context) ([]string, error) {
	feeds, err := f.repo.GetAllFeeds(false)
	if err != nil {
		return nil, err
	}

	domainMap := make(map[string]bool)
	for _, feed := range feeds {
		if !feed.Enabled {
			continue
		}
		if feed.RootDomain != nil && *feed.RootDomain != "" {
			domainMap[*feed.RootDomain] = true
		}
	}

	domains := make([]string, 0, len(domainMap))
	for d := range domainMap {
		domains = append(domains, d)
	}
	sort.Strings(domains)

	return domains, nil
}
