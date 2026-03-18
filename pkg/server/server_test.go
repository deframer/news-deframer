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
	getItemsForRootDomain func(ctx context.Context, rootDomain string, maxScore float64) ([]facade.AnalyzedItem, error)
	getFirstItemForUrl    func(ctx context.Context, u *url.URL) (*facade.AnalyzedItem, error)
	getArticlesByTrend    func(ctx context.Context, term string, domain string, date *time.Time, days int, offset int, limit int) ([]database.AnalyzedArticle, error)
	getSentimentsByTrend  func(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error)
}

func (m *mockFacade) GetRssProxyFeed(ctx context.Context, filter *facade.RSSProxyFilter) (string, error) {
	return "", nil
}

func (m *mockFacade) GetItemsForRootDomain(ctx context.Context, rootDomain string, maxScore float64) ([]facade.AnalyzedItem, error) {
	if m.getItemsForRootDomain != nil {
		return m.getItemsForRootDomain(ctx, rootDomain, maxScore)
	}
	return nil, nil
}

func (m *mockFacade) GetFirstItemForUrl(ctx context.Context, u *url.URL) (*facade.AnalyzedItem, error) {
	if m.getFirstItemForUrl != nil {
		return m.getFirstItemForUrl(ctx, u)
	}
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

func (m *mockFacade) GetSentimentsByTrend(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
	if m.getSentimentsByTrend != nil {
		return m.getSentimentsByTrend(ctx, term, domain, date, days)
	}
	return nil, nil
}

func TestHandleArticles(t *testing.T) {
	ctx := context.Background()
	today := time.Now().UTC().Format("2006-01-02")

	t.Run("success", func(t *testing.T) {
		title := "T"
		rating := 0.2
		authors := database.StringArray{"Jane Doe"}
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
				return []database.AnalyzedArticle{{URL: "https://example.com/a", Rating: &rating, Title: &title, Authors: authors}}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/articles?root=example.com&term=ai&date="+today, nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "https://example.com/a")
		assert.Contains(t, rr.Body.String(), `"authors":["Jane Doe"]`)
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

func TestHandleSentiments(t *testing.T) {
	ctx := context.Background()
	today := time.Now().UTC().Format("2006-01-02")

	t.Run("success", func(t *testing.T) {
		todayPtr, err := parseOptionalDateParam(today)
		assert.NoError(t, err)
		valence := 0.7
		joy := 0.9

		mockF := &mockFacade{
			getSentimentsByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
				assert.Equal(t, "ai", term)
				assert.Equal(t, "example.com", domain)
				assert.Equal(t, todayPtr, date)
				assert.Equal(t, 1, days)

				return &database.SentimentItem{Sentiments: &database.SentimentScores{Valence: valence, Joy: joy}}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/sentiments?root=example.com&term=ai&date="+today, nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"valence":0.7`)
	})

	t.Run("missing params", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/api/sentiments?root=example.com", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("invalid date", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/api/sentiments?root=example.com&term=ai&date=03-03-2026", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("facade error", func(t *testing.T) {
		mockF := &mockFacade{
			getSentimentsByTrend: func(ctx context.Context, term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
				return nil, errors.New("boom")
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/sentiments?root=example.com&term=ai&date="+today, nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}

func TestHandleItem(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockF := &mockFacade{
			getFirstItemForUrl: func(ctx context.Context, u *url.URL) (*facade.AnalyzedItem, error) {
				assert.Equal(t, "https://example.com/a", u.String())
				return &facade.AnalyzedItem{
					Hash:    "hash1",
					URL:     "https://example.com/a",
					Authors: database.StringArray{"Jane Doe", "John Roe"},
					ThinkResult: &database.ThinkResult{
						TitleCorrected: "Corrected Title",
					},
				}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/item?url=https://example.com/a", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"authors":["Jane Doe","John Roe"]`)
	})

	t.Run("mobile api alias", func(t *testing.T) {
		mockF := &mockFacade{
			getFirstItemForUrl: func(ctx context.Context, u *url.URL) (*facade.AnalyzedItem, error) {
				assert.Equal(t, "https://example.com/a", u.String())
				return &facade.AnalyzedItem{
					Hash:    "hash1",
					URL:     "https://example.com/a",
					Authors: database.StringArray{"Jane Doe"},
					ThinkResult: &database.ThinkResult{
						TitleCorrected: "Corrected Title",
					},
				}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/mobile/api/item?url=https://example.com/a", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"authors":["Jane Doe"]`)
	})
}

func TestHandleSite(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockF := &mockFacade{
			getItemsForRootDomain: func(ctx context.Context, rootDomain string, maxScore float64) ([]facade.AnalyzedItem, error) {
				assert.Equal(t, "example.com", rootDomain)
				return []facade.AnalyzedItem{{
					Hash:    "hash1",
					URL:     "https://example.com/a",
					Authors: database.StringArray{"Jane Doe"},
					ThinkResult: &database.ThinkResult{
						TitleCorrected: "Corrected Title",
					},
				}}, nil
			},
		}

		s := New(ctx, &config.Config{DisableETag: true}, mockF)
		req := httptest.NewRequest(http.MethodGet, "/api/site?root=example.com", nil)
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"authors":["Jane Doe"]`)
	})
}

func TestCORS(t *testing.T) {
	ctx := context.Background()

	t.Run("wildcard origin", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true, CORSAllowedOrigins: "*"}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://192.168.192.124:8090")
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Authorization, Content-Type", rr.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("allowed origin from list", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true, CORSAllowedOrigins: "http://localhost:8090, http://192.168.192.124:8090"}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://192.168.192.124:8090")
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "http://192.168.192.124:8090", rr.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, rr.Header().Values("Vary"), "Origin")
	})

	t.Run("disallowed origin", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true, CORSAllowedOrigins: "http://localhost:8090"}, &mockFacade{})
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		req.Header.Set("Origin", "http://192.168.192.124:8090")
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("preflight request", func(t *testing.T) {
		s := New(ctx, &config.Config{DisableETag: true, CORSAllowedOrigins: "*"}, &mockFacade{})
		req := httptest.NewRequest(http.MethodOptions, "/api/domains", nil)
		req.Header.Set("Origin", "http://192.168.192.124:8090")
		rr := httptest.NewRecorder()

		s.httpServer.Handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	})
}
