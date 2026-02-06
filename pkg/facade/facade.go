package facade

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sort"
	"strings"
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

type DomainEntry struct {
	Domain   string `json:"domain"`
	Language string `json:"language"`
}

type Facade interface {
	GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error)
	GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]AnalyzedItem, error)
	GetFirstItemForUrl(ctx context.Context, u *url.URL) (*AnalyzedItem, error)
	GetRootDomains(ctx context.Context) ([]DomainEntry, error)
	GetTopTrendByDomain(ctx context.Context, domain string, language string, daysInPast int) ([]database.TrendMetric, error)
	GetContextByDomain(ctx context.Context, term string, domain string, language string, daysInPast int) ([]database.TrendContext, error)
	GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, daysInPast int) ([]database.Lifecycle, error)
	GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, daysInPast int) ([]database.DomainComparison, error)
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

func (f *facade) GetRootDomains(ctx context.Context) ([]DomainEntry, error) {
	feeds, err := f.repo.GetAllFeeds(false)
	if err != nil {
		return nil, err
	}

	entryMap := make(map[string]DomainEntry)
	for _, feed := range feeds {
		if !feed.Enabled {
			continue
		}
		if feed.RootDomain != nil && *feed.RootDomain != "" {
			lang := ""
			if feed.Language != nil {
				lang = *feed.Language
			}
			key := *feed.RootDomain + "|" + lang
			entryMap[key] = DomainEntry{Domain: *feed.RootDomain, Language: lang}
		}
	}

	domains := make([]DomainEntry, 0, len(entryMap))
	for _, entry := range entryMap {
		domains = append(domains, entry)
	}
	sort.Slice(domains, func(i, j int) bool {
		d1 := strings.ToLower(domains[i].Domain)
		d2 := strings.ToLower(domains[j].Domain)
		if d1 == d2 {
			return domains[i].Language < domains[j].Language
		}
		return d1 < d2
	})

	return domains, nil
}

func (f *facade) GetTopTrendByDomain(ctx context.Context, domain string, language string, daysInPast int) ([]database.TrendMetric, error) {
	return f.repo.GetTopTrendByDomain(domain, language, daysInPast)
}

func (f *facade) GetContextByDomain(ctx context.Context, term string, domain string, language string, daysInPast int) ([]database.TrendContext, error) {
	return f.repo.GetContextByDomain(term, domain, language, daysInPast)
}

func (f *facade) GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, daysInPast int) ([]database.Lifecycle, error) {
	return f.repo.GetLifecycleByDomain(term, domain, language, daysInPast)
}

func (f *facade) GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, daysInPast int) ([]database.DomainComparison, error) {
	return f.repo.GetDomainComparison(domainA, domainB, language, daysInPast, database.DomainComparisonUtilityThreshold, database.DomainComparisonOutlierRatioThreshold, database.DomainComparisonLimit)
}
