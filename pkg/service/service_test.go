package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"testing"
	"time"

	web "github.com/deframer/news-deframer/gen/web"
	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFacade struct {
	getItemsForRootDomain func(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error)
	getFirstItemForUrl    func(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error)
	getArticlesByTrend    func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error)
	getSentimentsByTrend  func(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error)
	getRootDomains        func(ctx context.Context) ([]facade.DomainEntry, error)
	getTopTrendByDomain   func(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error)
	getContextByDomain    func(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error)
	getLifecycleByDomain  func(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error)
	getDomainComparison   func(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error)
}

func (m *mockFacade) GetRssProxyFeed(ctx context.Context, filter *facade.RSSProxyFilter) (string, error) {
	return "", nil
}

func (m *mockFacade) GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error) {
	if m.getItemsForRootDomain != nil {
		return m.getItemsForRootDomain(ctx, rootDomain, maxScore)
	}
	return nil, nil
}

func (m *mockFacade) GetFirstItemForUrl(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error) {
	if m.getFirstItemForUrl != nil {
		return m.getFirstItemForUrl(ctx, u)
	}
	return nil, nil
}

func (m *mockFacade) GetRootDomains(ctx context.Context) ([]facade.DomainEntry, error) {
	if m.getRootDomains != nil {
		return m.getRootDomains(ctx)
	}
	return nil, nil
}

func (m *mockFacade) GetTopTrendByDomain(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
	if m.getTopTrendByDomain != nil {
		return m.getTopTrendByDomain(ctx, domain, language, date, days)
	}
	return nil, nil
}

func (m *mockFacade) GetContextByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
	if m.getContextByDomain != nil {
		return m.getContextByDomain(ctx, term, domain, language, date, days)
	}
	return nil, nil
}

func (m *mockFacade) GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
	if m.getLifecycleByDomain != nil {
		return m.getLifecycleByDomain(ctx, term, domain, language, date, days)
	}
	return nil, nil
}

func (m *mockFacade) GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error) {
	if m.getDomainComparison != nil {
		return m.getDomainComparison(ctx, domainA, domainB, language, date, days)
	}
	return nil, nil
}

func (m *mockFacade) GetArticlesByTrend(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
	if m.getArticlesByTrend != nil {
		return m.getArticlesByTrend(ctx, term, domain, date, days, offset, limit)
	}
	return nil, nil
}

func (m *mockFacade) GetSentimentsByTrend(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
	if m.getSentimentsByTrend != nil {
		return m.getSentimentsByTrend(ctx, term, domain, date, days)
	}
	return nil, nil
}

func testService(mf *mockFacade) *Service {
	return &Service{
		facade: mf,
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		cfg:    &config.Config{},
	}
}

func TestItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		timeVal := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
		rating := 0.5
		mediaURL := "https://cdn.example.com/media.jpg"
		mediaType := "image/jpeg"
		mediaMedium := "image"
		mediaTitle := "Media title"
		mediaDesc := "Media description"
		credit := "Photo credit"
		thumbURL := "https://cdn.example.com/thumb.jpg"
		thumb := &database.MediaThumbnail{URL: thumbURL, Height: 120, Width: 200}
		mockF := &mockFacade{
			getFirstItemForUrl: func(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error) {
				require.Equal(t, "https://example.com/a", u.String())
				return &database.AnalyzedItem{
					Hash:    "hash1",
					URL:     "https://example.com/a",
					Authors: database.StringArray{"Jane Doe", "John Roe"},
					Sentiments: &database.SentimentScores{Valence: 0.7, Joy: 0.9},
					SentimentsDeframed: &database.SentimentScores{Valence: 0.2},
					MediaContent: &database.MediaContent{
						URL:         mediaURL,
						Type:        mediaType,
						Medium:      mediaMedium,
						Height:      640,
						Width:       480,
						Title:       mediaTitle,
						Description: mediaDesc,
						Thumbnail:   thumb,
						Credit:      credit,
					},
					ThinkResult: &database.ThinkResult{TitleCorrected: "Corrected Title"},
					ThinkRating: rating,
					PubDate:     timeVal,
				}, nil
			},
		}

		res, err := testService(mockF).Item(ctx, &web.ItemPayload{URL: "https://example.com/a"})
		require.NoError(t, err)
		require.NotNil(t, res)
		assert.Equal(t, "hash1", res.Hash)
		assert.Equal(t, "https://example.com/a", res.URL)
		assert.Equal(t, rating, res.Rating)
		assert.Equal(t, []string{"Jane Doe", "John Roe"}, res.Authors)
		assert.Equal(t, timeVal.Format(time.RFC3339), res.PubDate)
		require.NotNil(t, res.Sentiments)
		assert.NotNil(t, res.Sentiments.Valence)
		assert.Equal(t, 0.7, *res.Sentiments.Valence)
		require.NotNil(t, res.Media)
		assert.Equal(t, mediaURL, res.Media.URL)
		assert.Equal(t, mediaType, *res.Media.Type)
		assert.Equal(t, mediaMedium, *res.Media.Medium)
		require.NotNil(t, res.Media.Thumbnail)
		assert.Equal(t, thumbURL, res.Media.Thumbnail.URL)
	})

	t.Run("missing url", func(t *testing.T) {
		_, err := testService(&mockFacade{}).Item(ctx, &web.ItemPayload{URL: ""})
		assert.Error(t, err)
	})

	t.Run("invalid url", func(t *testing.T) {
		_, err := testService(&mockFacade{}).Item(ctx, &web.ItemPayload{URL: "::bad"})
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockF := &mockFacade{getFirstItemForUrl: func(ctx context.Context, u *url.URL) (*database.AnalyzedItem, error) {
			return nil, nil
		}}
		_, err := testService(mockF).Item(ctx, &web.ItemPayload{URL: "https://example.com/a"})
		assert.Error(t, err)
	})
}

func TestSite(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		rating := 0.2
		mockF := &mockFacade{getItemsForRootDomain: func(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error) {
			require.Equal(t, "example.com", rootDomain)
			require.Equal(t, 0.3, maxScore)
			return []database.AnalyzedItem{{URL: "https://example.com/a", ThinkRating: rating, PubDate: time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)}}, nil
		}}
		res, err := testService(mockF).Site(ctx, &web.SitePayload{Root: "example.com", MaxScore: 0.3})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "https://example.com/a", res[0].URL)
		assert.Equal(t, rating, res[0].Rating)
	})

	t.Run("missing root", func(t *testing.T) {
		_, err := testService(&mockFacade{}).Site(ctx, &web.SitePayload{})
		assert.Error(t, err)
	})

	t.Run("facade error", func(t *testing.T) {
		mockF := &mockFacade{getItemsForRootDomain: func(ctx context.Context, rootDomain string, maxScore float64) ([]database.AnalyzedItem, error) {
			return nil, errors.New("boom")
		}}
		_, err := testService(mockF).Site(ctx, &web.SitePayload{Root: "example.com"})
		assert.Error(t, err)
	})
}

func TestArticlesAndSentiments(t *testing.T) {
	ctx := context.Background()
	today := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	todayStr := today.Format("2006-01-02")

	t.Run("articles success", func(t *testing.T) {
		title := "T"
		rating := 0.2
		mockF := &mockFacade{getArticlesByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
			require.Equal(t, "ai", term)
			require.Equal(t, "example.com", domain)
			require.Equal(t, today, *date)
			require.Equal(t, 1, days)
			require.Equal(t, 0, offset)
			require.Equal(t, 20, limit)
			return []database.AnalyzedArticle{{URL: "https://example.com/a", Rating: &rating, Title: &title, Authors: database.StringArray{"Jane Doe"}, PubDate: time.Date(2026, 5, 1, 11, 0, 0, 0, time.UTC)}}, nil
		}}
		res, err := testService(mockF).Articles(ctx, &web.ArticlesPayload{Root: "example.com", Term: "ai", Date: &todayStr, Days: 1, Offset: 0, Limit: 20})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "https://example.com/a", res[0].URL)
		assert.Equal(t, "2026-05-01T11:00:00Z", res[0].PubDate)
	})

	t.Run("sentiments success", func(t *testing.T) {
		valence := 0.7
		joy := 0.9
		mockF := &mockFacade{getSentimentsByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
			require.Equal(t, "ai", term)
			require.Equal(t, "example.com", domain)
			require.Equal(t, today, *date)
			require.Equal(t, 1, days)
			return &database.SentimentItem{Sentiments: &database.SentimentScores{Valence: valence, Joy: joy}}, nil
		}}
		res, err := testService(mockF).Sentiments(ctx, &web.SentimentsPayload{Root: "example.com", Term: "ai", Date: &todayStr, Days: 1})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotNil(t, res.Sentiments)
		assert.Equal(t, valence, *res.Sentiments.Valence)
		assert.Equal(t, joy, *res.Sentiments.Joy)
	})
}

func TestDomainsAndTrends(t *testing.T) {
	ctx := context.Background()
	today := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	todayStr := today.Format("2006-01-02")

	t.Run("domains success", func(t *testing.T) {
		portal := "https://example.com"
		mockF := &mockFacade{getRootDomains: func(ctx context.Context) ([]facade.DomainEntry, error) {
			return []facade.DomainEntry{{Domain: "example.com", Language: "en", PortalUrl: &portal}}, nil
		}}
		res, err := testService(mockF).Domains(ctx, &web.DomainsPayload{})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "example.com", res[0].Domain)
		assert.Equal(t, portal, *res[0].PortalURL)
	})

	t.Run("top trends success", func(t *testing.T) {
		mockF := &mockFacade{getTopTrendByDomain: func(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
			require.Equal(t, "example.com", domain)
			require.Equal(t, "en", language)
			require.Equal(t, today, *date)
			require.Equal(t, 1, days)
			return []database.TrendMetric{{TrendTopic: "t", Frequency: 2, Utility: 3, OutlierRatio: 1.5, TimeSlice: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)}}, nil
		}}
		res, err := testService(mockF).TopTrendsByDomain(ctx, &web.TopTrendsByDomainPayload{Domain: "example.com", Lang: "en", Date: &todayStr, Days: 1})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "t", res[0].TrendTopic)
	})

	t.Run("context success", func(t *testing.T) {
		mockF := &mockFacade{getContextByDomain: func(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
			require.Equal(t, "ai", term)
			require.Equal(t, "example.com", domain)
			require.Equal(t, "en", language)
			return []database.TrendContext{{Context: "ctx", Frequency: 2}}, nil
		}}
		res, err := testService(mockF).ContextByDomain(ctx, &web.ContextByDomainPayload{Term: "ai", Domain: "example.com", Lang: "en", Days: 1})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "ctx", res[0].Context)
	})

	t.Run("lifecycle success", func(t *testing.T) {
		mockF := &mockFacade{getLifecycleByDomain: func(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
			return []database.Lifecycle{{TimeSlice: today, Frequency: 4, Velocity: 5}}, nil
		}}
		res, err := testService(mockF).LifecycleByDomain(ctx, &web.LifecycleByDomainPayload{Term: "ai", Domain: "example.com", Lang: "en", Days: 1})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "2026-05-01T00:00:00Z", res[0].TimeSlice)
	})

	t.Run("comparison success", func(t *testing.T) {
		mockF := &mockFacade{getDomainComparison: func(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error) {
			return []database.DomainComparison{{Classification: "INTERSECT", RankGroup: 1, TrendTopic: "t", ScoreA: 0.7, ScoreB: 0.8}}, nil
		}}
		res, err := testService(mockF).DomainComparisonEndpoint(ctx, &web.DomainComparisonPayload{DomainA: "a", DomainB: "b", Lang: "en", Days: 1})
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "INTERSECT", res[0].Classification)
	})
}
