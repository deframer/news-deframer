package facade

import (
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockDownloader struct {
	downloadRSSFeed func(ctx context.Context, feed *url.URL) (io.ReadCloser, error)
	resolveRedirect func(ctx context.Context, targetURL string) (string, error)
}

func (m *mockDownloader) DownloadRSSFeed(ctx context.Context, feed *url.URL) (io.ReadCloser, error) {
	if m.downloadRSSFeed != nil {
		return m.downloadRSSFeed(ctx, feed)
	}
	return nil, nil
}

func (m *mockDownloader) ResolveRedirect(ctx context.Context, targetURL string) (string, error) {
	if m.resolveRedirect != nil {
		return m.resolveRedirect(ctx, targetURL)
	}
	return targetURL, nil
}

type mockRepo struct {
	findFeedByUrl                func(u *url.URL) (*database.Feed, error)
	findFeedByUrlAndAvailability func(u *url.URL, onlyEnabled bool) (*database.Feed, error)
	findFeedById                 func(feedID uuid.UUID) (*database.Feed, error)
	upsertFeed                   func(feed *database.Feed) error
	findItemsByUrl               func(u *url.URL) ([]database.Item, error)
	getAllFeeds                  func(deleted bool) ([]database.Feed, error)
	deleteFeedById               func(id uuid.UUID) error
	enqueueSync                  func(id uuid.UUID, pollingInterval time.Duration) error
	removeSync                   func(id uuid.UUID) error
	beginFeedUpdate              func(lockDuration time.Duration) (*database.Feed, error)
	endFeedUpdate                func(id uuid.UUID, err error, successDelay time.Duration) error
	getPendingItems              func(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error)
	upsertItem                   func(item *database.Item) error
	getItemsByHashes             func(feedID uuid.UUID, hashes []string) ([]database.Item, error)
	beginThinkFixerBatch         func(limit int, since time.Time, minErrorCount int, maxErrorCount int, lockDuration time.Duration) ([]database.Item, error)
	upsertCachedFeed             func(cachedFeed *database.CachedFeed) error
	findCachedFeedById           func(feedID uuid.UUID) (*database.CachedFeed, error)
	findFeedScheduleById         func(feedID uuid.UUID) (*database.FeedSchedule, error)
	findItemsByRootDomain        func(rootDomain string, limit int) ([]database.Item, error)
	getTopTrendByDomain          func(domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error)
	getContextByDomain           func(term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error)
	getLifecycleByDomain         func(term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error)
	getDomainComparison          func(domainA string, domainB string, language string, date *time.Time, days int, utilityThreshold float64, outlierRatioThreshold float64, limit int) ([]database.DomainComparison, error)
	getArticlesByTrend           func(term string, domain string, date *time.Time, days int) ([]database.AnalyzedArticle, error)
}

func mustParseTestDate(raw string) *time.Time {
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		panic(err)
	}
	normalized := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized
}

func (m *mockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) {
	if m.findFeedByUrl != nil {
		return m.findFeedByUrl(u)
	}
	return nil, nil
}

func (m *mockRepo) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*database.Feed, error) {
	if m.findFeedByUrlAndAvailability != nil {
		return m.findFeedByUrlAndAvailability(u, onlyEnabled)
	}
	return nil, nil
}

func (m *mockRepo) FindFeedById(feedID uuid.UUID) (*database.Feed, error) {
	if m.findFeedById != nil {
		return m.findFeedById(feedID)
	}
	return nil, nil
}

func (m *mockRepo) UpsertFeed(feed *database.Feed) error {
	if m.upsertFeed != nil {
		return m.upsertFeed(feed)
	}
	return nil
}

func (m *mockRepo) FindItemsByUrl(u *url.URL) ([]database.Item, error) {
	if m.findItemsByUrl != nil {
		return m.findItemsByUrl(u)
	}
	return nil, nil
}

func (m *mockRepo) GetAllFeeds(deleted bool) ([]database.Feed, error) {
	if m.getAllFeeds != nil {
		return m.getAllFeeds(deleted)
	}
	return nil, nil
}

func (m *mockRepo) DeleteFeedById(id uuid.UUID) error {
	if m.deleteFeedById != nil {
		return m.deleteFeedById(id)
	}
	return nil
}

func (m *mockRepo) PurgeFeedById(id uuid.UUID) error {
	return nil
}

func (m *mockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration) error {
	if m.enqueueSync != nil {
		return m.enqueueSync(id, pollingInterval)
	}
	return nil
}

func (m *mockRepo) EnqueueMine(id uuid.UUID, miningInterval time.Duration) error {
	return nil
}

func (m *mockRepo) RemoveMine(id uuid.UUID) error {
	return nil
}

func (m *mockRepo) RemoveSync(id uuid.UUID) error {
	if m.removeSync != nil {
		return m.removeSync(id)
	}
	return nil
}

func (m *mockRepo) BeginFeedUpdate(lockDuration time.Duration) (*database.Feed, error) {
	if m.beginFeedUpdate != nil {
		return m.beginFeedUpdate(lockDuration)
	}
	return nil, nil
}

func (m *mockRepo) EndFeedUpdate(id uuid.UUID, err error, successDelay time.Duration) error {
	if m.endFeedUpdate != nil {
		return m.endFeedUpdate(id, err, successDelay)
	}
	return nil
}

func (m *mockRepo) GetPendingItems(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error) {
	if m.getPendingItems != nil {
		return m.getPendingItems(feedID, hashes, maxRetries)
	}
	res := make(map[string]int)
	for _, h := range hashes {
		res[h] = 0
	}
	return res, nil
}

func (m *mockRepo) UpsertItem(item *database.Item) error {
	if m.upsertItem != nil {
		return m.upsertItem(item)
	}
	return nil
}

func (m *mockRepo) GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]database.Item, error) {
	if m.getItemsByHashes != nil {
		return m.getItemsByHashes(feedID, hashes)
	}
	return nil, nil
}

func (m *mockRepo) BeginThinkFixerBatch(limit int, since time.Time, minErrorCount int, maxErrorCount int, lockDuration time.Duration) ([]database.Item, error) {
	if m.beginThinkFixerBatch != nil {
		return m.beginThinkFixerBatch(limit, since, minErrorCount, maxErrorCount, lockDuration)
	}
	return nil, nil
}

func (m *mockRepo) UpsertCachedFeed(cachedFeed *database.CachedFeed) error {
	if m.upsertCachedFeed != nil {
		return m.upsertCachedFeed(cachedFeed)
	}
	return nil
}

func (m *mockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	if m.findCachedFeedById != nil {
		return m.findCachedFeedById(feedID)
	}
	return nil, nil
}

func (m *mockRepo) FindFeedScheduleById(feedID uuid.UUID) (*database.FeedSchedule, error) {
	if m.findFeedScheduleById != nil {
		return m.findFeedScheduleById(feedID)
	}
	return nil, nil
}

func (m *mockRepo) CreateFeedSchedule(feedID uuid.UUID) error {
	return nil
}

func (m *mockRepo) FindItemsByRootDomain(rootDomain string, limit int) ([]database.Item, error) {
	if m.findItemsByRootDomain != nil {
		return m.findItemsByRootDomain(rootDomain, limit)
	}
	return nil, nil
}

func (m *mockRepo) GetTopTrendByDomain(domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
	if m.getTopTrendByDomain != nil {
		return m.getTopTrendByDomain(domain, language, date, days)
	}
	return nil, nil
}

func (m *mockRepo) GetContextByDomain(term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
	if m.getContextByDomain != nil {
		return m.getContextByDomain(term, domain, language, date, days)
	}
	return nil, nil
}

func (m *mockRepo) GetLifecycleByDomain(term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
	if m.getLifecycleByDomain != nil {
		return m.getLifecycleByDomain(term, domain, language, date, days)
	}
	return nil, nil
}

func (m *mockRepo) GetDomainComparison(domainA string, domainB string, language string, date *time.Time, days int, utilityThreshold float64, outlierRatioThreshold float64, limit int) ([]database.DomainComparison, error) {
	if m.getDomainComparison != nil {
		return m.getDomainComparison(domainA, domainB, language, date, days, utilityThreshold, outlierRatioThreshold, limit)
	}
	return nil, nil
}

func (m *mockRepo) GetArticlesByTrend(term string, domain string, date *time.Time, days int) ([]database.AnalyzedArticle, error) {
	if m.getArticlesByTrend != nil {
		return m.getArticlesByTrend(term, domain, date, days)
	}
	return nil, nil
}

func TestGetRssProxyFeed(t *testing.T) {
	ctx := context.Background()
	rssContent := `
	<rss version="2.0">
		<channel>
			<title>Dummy Feed</title>
		</channel>
	</rss>`
	mockDL := &mockDownloader{
		downloadRSSFeed: func(ctx context.Context, feed *url.URL) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(rssContent)), nil
		},
	}

	mockR := &mockRepo{
		findFeedByUrl: func(u *url.URL) (*database.Feed, error) {
			return &database.Feed{}, nil
		},
		findCachedFeedById: func(id uuid.UUID) (*database.CachedFeed, error) {
			return &database.CachedFeed{XMLContent: &rssContent}, nil
		},
	}
	f := New(ctx, nil, mockR, mockDL)

	filter := RSSProxyFilter{
		URL:      "http://example.com",
		Lang:     "en",
		Max:      0.75,
		Embedded: true,
	}

	xmlData, err := f.GetRssProxyFeed(ctx, &filter)
	assert.NoError(t, err)
	assert.Contains(t, xmlData, "<title>Dummy Feed</title>")
}

func TestGetItemsForRootDomain(t *testing.T) {
	ctx := context.Background()
	rootDomain := "example.com"

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		expectedItems := []database.Item{
			{
				Hash: "hash1",
				URL:  "http://example.com/1",
				ThinkResult: &database.ThinkResult{
					TitleCorrected: "Corrected Title 1",
				},
				MediaContent: &database.MediaContent{
					URL: "http://example.com/img1.jpg",
				},
				ThinkRating: 0.5,
				PubDate:     now,
			},
			{
				Hash: "hash2",
				URL:  "http://example.com/2",
				// Nil ThinkResult and MediaContent to test handling
				PubDate: now.Add(-1 * time.Hour),
			},
		}

		mockR := &mockRepo{
			findItemsByRootDomain: func(domain string, limit int) ([]database.Item, error) {
				assert.Equal(t, rootDomain, domain)
				assert.Equal(t, MaxItemsForRootDomain, limit)
				return expectedItems, nil
			},
		}

		f := New(ctx, nil, mockR)

		items, err := f.GetItemsForRootDomain(ctx, rootDomain, 0.0)
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		// Verify Item 1
		assert.Equal(t, "hash1", items[0].Hash)
		assert.Equal(t, "http://example.com/1", items[0].URL)
		assert.Equal(t, "Corrected Title 1", items[0].TitleCorrected)
		assert.Equal(t, "http://example.com/img1.jpg", items[0].MediaContent.URL)
		assert.Equal(t, 0.5, items[0].ThinkRating)
		assert.Equal(t, now, items[0].PubDate)

		// Verify Item 2
		assert.Equal(t, "hash2", items[1].Hash)
		assert.Equal(t, "http://example.com/2", items[1].URL)
		assert.Empty(t, items[1].TitleCorrected)
		assert.Nil(t, items[1].MediaContent)
		assert.Equal(t, 0.0, items[1].ThinkRating)
		assert.Equal(t, now.Add(-1*time.Hour), items[1].PubDate)
	})

	t.Run("WithFilter", func(t *testing.T) {
		expectedItems := []database.Item{
			{
				Hash:        "hash1",
				ThinkRating: 0.8,
			},
			{
				Hash:        "hash2",
				ThinkRating: 0.2,
			},
		}

		mockR := &mockRepo{
			findItemsByRootDomain: func(domain string, limit int) ([]database.Item, error) {
				return expectedItems, nil
			},
		}
		f := New(ctx, nil, mockR)
		items, err := f.GetItemsForRootDomain(ctx, rootDomain, 0.5)
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "hash2", items[0].Hash)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			findItemsByRootDomain: func(domain string, limit int) ([]database.Item, error) {
				return nil, assert.AnError
			},
		}
		f := New(ctx, nil, mockR)
		items, err := f.GetItemsForRootDomain(ctx, rootDomain, 0.0)
		assert.Error(t, err)
		assert.Nil(t, items)
	})
}

func TestGetFirstItemForUrl(t *testing.T) {
	ctx := context.Background()
	targetURL := "http://example.com/article"
	u, _ := url.Parse(targetURL)

	t.Run("Found", func(t *testing.T) {
		now := time.Now()
		expectedItem := database.Item{
			Hash: "hash1",
			URL:  targetURL,
			ThinkResult: &database.ThinkResult{
				TitleCorrected: "Corrected Title",
			},
			ThinkRating: 0.8,
			PubDate:     now,
		}

		mockR := &mockRepo{
			findItemsByUrl: func(u *url.URL) ([]database.Item, error) {
				assert.Equal(t, targetURL, u.String())
				return []database.Item{expectedItem}, nil
			},
		}

		f := New(ctx, nil, mockR)

		item, err := f.GetFirstItemForUrl(ctx, u)
		assert.NoError(t, err)
		assert.NotNil(t, item)
		assert.Equal(t, "hash1", item.Hash)
		assert.Equal(t, "Corrected Title", item.TitleCorrected)
		assert.Equal(t, now, item.PubDate)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockR := &mockRepo{
			findItemsByUrl: func(u *url.URL) ([]database.Item, error) {
				return []database.Item{}, nil
			},
		}
		f := New(ctx, nil, mockR)
		item, err := f.GetFirstItemForUrl(ctx, u)
		assert.NoError(t, err)
		assert.Nil(t, item)
	})
}

func TestGetRootDomains(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		domainA := "a.com"
		domainB := "b.com"
		domainC := "c.com"
		langEn := "en"
		langDe := "de"
		portalUrl := "http://portal.b.com"

		feeds := []database.Feed{
			{
				Base:       database.Base{},
				Enabled:    true,
				RootDomain: &domainB, // b.com
				Language:   &langEn,
				PortalUrl:  &portalUrl,
			},
			{
				Base:       database.Base{},
				Enabled:    true,
				RootDomain: &domainA, // a.com
				Language:   &langDe,
			},
			{
				Base:       database.Base{},
				Enabled:    true,
				RootDomain: &domainB, // Duplicate b.com
				Language:   &langEn,
			},
			{
				Base:       database.Base{},
				Enabled:    false,
				RootDomain: &domainC, // Disabled, should be ignored
			},
			{
				Base:       database.Base{},
				Enabled:    true,
				RootDomain: nil, // Nil, should be ignored
			},
		}

		mockR := &mockRepo{
			getAllFeeds: func(deleted bool) ([]database.Feed, error) {
				assert.False(t, deleted)
				return feeds, nil
			},
		}

		f := New(ctx, nil, mockR)
		domains, err := f.GetRootDomains(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []DomainEntry{
			{Domain: "a.com", Language: "de"},
			{Domain: "b.com", Language: "en", PortalUrl: &portalUrl},
		}, domains)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			getAllFeeds: func(deleted bool) ([]database.Feed, error) {
				return nil, assert.AnError
			},
		}

		f := New(ctx, nil, mockR)
		domains, err := f.GetRootDomains(ctx)
		assert.Error(t, err)
		assert.Nil(t, domains)
	})
}

func TestGetTopTrendByDomain(t *testing.T) {
	ctx := context.Background()
	domain := "example.com"
	lang := "en"
	date := mustParseTestDate("2026-03-03")
	days := 7

	t.Run("Success", func(t *testing.T) {
		expected := []database.TrendMetric{{TrendTopic: "topic1", Frequency: 10}}
		mockR := &mockRepo{
			getTopTrendByDomain: func(d, l string, dt *time.Time, days int) ([]database.TrendMetric, error) {
				assert.Equal(t, domain, d)
				assert.Equal(t, lang, l)
				assert.Equal(t, date, dt)
				return expected, nil
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetTopTrendByDomain(ctx, domain, lang, date, days)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			getTopTrendByDomain: func(d, l string, dt *time.Time, days int) ([]database.TrendMetric, error) {
				return nil, assert.AnError
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetTopTrendByDomain(ctx, domain, lang, date, days)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetContextByDomain(t *testing.T) {
	ctx := context.Background()
	term := "test"
	domain := "example.com"
	lang := "en"
	date := mustParseTestDate("2026-03-03")
	days := 7

	t.Run("Success", func(t *testing.T) {
		expected := []database.TrendContext{{Context: "ctx1", Frequency: 5}}
		mockR := &mockRepo{
			getContextByDomain: func(tm, d, l string, dt *time.Time, days int) ([]database.TrendContext, error) {
				assert.Equal(t, term, tm)
				assert.Equal(t, domain, d)
				assert.Equal(t, date, dt)
				return expected, nil
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetContextByDomain(ctx, term, domain, lang, date, days)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestGetLifecycleByDomain(t *testing.T) {
	ctx := context.Background()
	term := "test-term"
	domain := "example.com"
	language := "en"
	date := mustParseTestDate("2026-03-03")
	days := 30

	t.Run("Success", func(t *testing.T) {
		expected := []database.Lifecycle{{Frequency: 100, Velocity: 50}}
		mockR := &mockRepo{
			getLifecycleByDomain: func(argTerm string, argDomain string, argLang string, argDate *time.Time, argDays int) ([]database.Lifecycle, error) {
				assert.Equal(t, term, argTerm)
				assert.Equal(t, domain, argDomain)
				assert.Equal(t, language, argLang)
				assert.Equal(t, date, argDate)
				assert.Equal(t, days, argDays)
				return expected, nil
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetLifecycleByDomain(ctx, term, domain, language, date, days)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			getLifecycleByDomain: func(argTerm string, argDomain string, argLang string, argDate *time.Time, argDays int) ([]database.Lifecycle, error) {
				return nil, assert.AnError
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetLifecycleByDomain(ctx, term, domain, language, date, days)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetDomainComparison(t *testing.T) {
	ctx := context.Background()
	domainA := "a.com"
	domainB := "b.com"
	lang := "en"
	date := mustParseTestDate("2026-03-03")
	days := 7

	t.Run("Success", func(t *testing.T) {
		expected := []database.DomainComparison{{TrendTopic: "topic1"}}
		mockR := &mockRepo{
			getDomainComparison: func(dA, dB, l string, dt *time.Time, d int, ut, ort float64, lim int) ([]database.DomainComparison, error) {
				assert.Equal(t, domainA, dA)
				assert.Equal(t, domainB, dB)
				assert.Equal(t, lang, l)
				assert.Equal(t, date, dt)
				assert.Equal(t, days, d)
				assert.Equal(t, database.DomainComparisonUtilityThreshold, ut)
				assert.Equal(t, database.DomainComparisonOutlierRatioThreshold, ort)
				assert.Equal(t, database.DomainComparisonLimit, lim)
				return expected, nil
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetDomainComparison(ctx, domainA, domainB, lang, date, days)
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			getDomainComparison: func(dA, dB, l string, dt *time.Time, d int, ut, ort float64, lim int) ([]database.DomainComparison, error) {
				return nil, assert.AnError
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetDomainComparison(ctx, domainA, domainB, lang, date, days)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetArticlesByTrend(t *testing.T) {
	ctx := context.Background()
	term := "test-term"
	domain := "example.com"
	date := mustParseTestDate("2026-03-03")
	days := 0

	t.Run("Success", func(t *testing.T) {
		title1 := "Corrected Title 1"
		rating1 := 0.5
		expectedItems := []database.AnalyzedArticle{{
			URL:    "http://example.com/1",
			Title:  &title1,
			Rating: &rating1,
		}, {
			URL:    "http://example.com/2",
			Title:  nil,
			Rating: nil,
		}}

		mockR := &mockRepo{
			getArticlesByTrend: func(argTerm string, argDomain string, argDate *time.Time, argDays int) ([]database.AnalyzedArticle, error) {
				assert.Equal(t, term, argTerm)
				assert.Equal(t, domain, argDomain)
				assert.Equal(t, date, argDate)
				assert.Equal(t, days, argDays)
				return expectedItems, nil
			},
		}
		f := New(ctx, nil, mockR)

		items, err := f.GetArticlesByTrend(ctx, term, domain, date, days)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		assert.Equal(t, "http://example.com/1", items[0].URL)
		if assert.NotNil(t, items[0].Title) {
			assert.Equal(t, "Corrected Title 1", *items[0].Title)
		}
		if assert.NotNil(t, items[0].Rating) {
			assert.Equal(t, 0.5, *items[0].Rating)
		}
		assert.Equal(t, "http://example.com/2", items[1].URL)
		assert.Nil(t, items[1].Title)
		assert.Nil(t, items[1].Rating)
	})

	t.Run("RepoError", func(t *testing.T) {
		mockR := &mockRepo{
			getArticlesByTrend: func(argTerm string, argDomain string, argDate *time.Time, argDays int) ([]database.AnalyzedArticle, error) {
				return nil, assert.AnError
			},
		}
		f := New(ctx, nil, mockR)
		res, err := f.GetArticlesByTrend(ctx, term, domain, date, days)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
