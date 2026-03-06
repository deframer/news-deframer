package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	"github.com/stretchr/testify/assert"
)

type mockFacade struct {
	getArticlesByTrend func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error)
}

func (m *mockFacade) GetRssProxyFeed(ctx context.Context, filter *facade.RSSProxyFilter) (string, error) {
	return "", nil
}

func (m *mockFacade) GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]facade.AnalyzedItem, error) {
	return nil, nil
}

func (m *mockFacade) GetFirstItemForUrl(ctx context.Context, u *url.URL) (*facade.AnalyzedItem, error) {
	return nil, nil
}

func (m *mockFacade) GetRootDomains(ctx context.Context) ([]facade.DomainEntry, error) {
	return nil, nil
}

func (m *mockFacade) GetTopTrendByDomain(ctx context.Context, domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
	return nil, nil
}

func (m *mockFacade) GetContextByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
	return nil, nil
}

func (m *mockFacade) GetLifecycleByDomain(ctx context.Context, term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
	return nil, nil
}

func (m *mockFacade) GetDomainComparison(ctx context.Context, domainA string, domainB string, language string, date *time.Time, days int) ([]database.DomainComparison, error) {
	return nil, nil
}

func (m *mockFacade) GetArticlesByTrend(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
	if m.getArticlesByTrend != nil {
		return m.getArticlesByTrend(ctx, term, domain, date, days, offset, limit)
	}
	return nil, nil
}

func TestHandleArticles(t *testing.T) {
	ctx := context.Background()
	today := time.Now().UTC().Format("2006-01-02")

	t.Run("success", func(t *testing.T) {
		title := "T"
		rating := 0.2
		todayPtr, err := parseOptionalDateParam(today)
		assert.NoError(t, err)
		mockF := &mockFacade{
			getArticlesByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
				assert.Equal(t, "ai", term)
				assert.Equal(t, "example.com", domain)
				assert.Equal(t, todayPtr, date)
				assert.Equal(t, 1, days)
				assert.Equal(t, 0, offset)
				assert.Equal(t, 20, limit)
				return []database.AnalyzedArticle{{URL: "https://example.com/a", Rating: &rating, Title: &title}}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/articles?root=example.com&term=ai&date="+today, nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "https://example.com/a")
	})

	t.Run("missing params", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/api/articles?root=example.com", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid date", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/api/articles?root=example.com&term=ai&date=03-03-2026", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("facade error", func(t *testing.T) {
		mockF := &mockFacade{
			getArticlesByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error) {
				return nil, errors.New("boom")
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/articles?root=example.com&term=ai&date="+today, nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
