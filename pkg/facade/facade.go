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

type DomainEntry struct {
	Domain    string  `json:"domain"`
	Language  string  `json:"language"`
	PortalUrl *string `json:"portal_url,omitempty"`
}

type Facade interface {
	GetRssProxyFeed(ctx context.Context, filter *RSSProxyFilter) (string, error)
	GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error)
	GetFirstItemForUrl(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error)
	GetRootDomains(ctx context.Context) ([]DomainEntry, error)
	GetTopTrendByDomain(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error)
	GetContextByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error)
	GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error)
	GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error)
	GetArticlesByTrend(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error)
	GetSentimentsByTrend(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error)
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

func (f *facade) GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error) {
	dbItems, err := f.repo.FindAnalyzedItemsByRootDomain(rootDomain, MaxItemsForRootDomain)
	if err != nil {
		return nil, err
	}

	items := make([]database.AnalyzedItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		if maxScore > 0 && dbItem.ThinkRating > maxScore {
			continue
		}
		items = append(items, dbItem)
	}
	return items, nil
}

func (f *facade) GetFirstItemForUrl(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error) {
	return f.repo.FindFirstAnalyzedItemByUrl(u)
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

			var portalUrl *string
			if existing, ok := entryMap[key]; ok {
				portalUrl = existing.PortalUrl
			}
			if feed.PortalUrl != nil {
				portalUrl = feed.PortalUrl
			}
			entryMap[key] = DomainEntry{Domain: *feed.RootDomain, Language: lang, PortalUrl: portalUrl}
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

func (f *facade) GetTopTrendByDomain(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
	return f.repo.GetTopTrendByDomain(domain, language, date, days)
}

func (f *facade) GetContextByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
	return f.repo.GetContextByDomain(term, domain, language, date, days)
}

func (f *facade) GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
	return f.repo.GetLifecycleByDomain(term, domain, language, date, days)
}

func (f *facade) GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error) {
	return f.repo.GetDomainComparison(domainA, domainB, language, date, days, database.DomainComparisonUtilityThreshold, database.DomainComparisonOutlierRatioThreshold, database.DomainComparisonLimit)
}

func (f *facade) GetArticlesByTrend(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
	return f.repo.GetArticlesByTrend(term, domain, date, days, offset, limit)
}

func (f *facade) GetSentimentsByTrend(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
	return f.repo.GetSentimentsByTrend(term, domain, date, days)
}
