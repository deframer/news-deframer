package downloader

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestDownloadRSSFeed_HTTP(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("mock content"))
	}))
	defer server.Close()

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	rc, err := d.DownloadRSSFeed(ctx, u)
	assert.NoError(t, err)
	defer func() { _ = rc.Close() }()
	content, err := io.ReadAll(rc)
	assert.NoError(t, err)
	assert.Equal(t, "mock content", string(content))
}

func TestDownloadRSSFeed_UnsupportedScheme(t *testing.T) {
	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	tests := []struct {
		url string
		err string
	}{
		{"file:///tmp/test", "unsupported scheme file"},
		{"invalid://test", "unsupported scheme invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			assert.NoError(t, err)

			rc, err := d.DownloadRSSFeed(ctx, u)
			assert.Error(t, err)
			assert.Nil(t, rc)
			assert.EqualError(t, err, tt.err)
		})
	}
}

func TestResolveRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/initial":
			http.Redirect(w, r, "/redirect1", http.StatusMovedPermanently)
		case "/redirect1":
			http.Redirect(w, r, "/final", http.StatusFound)
		case "/final":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	t.Run("FollowRedirects", func(t *testing.T) {
		finalURL, err := d.ResolveRedirect(ctx, ts.URL+"/initial")
		assert.NoError(t, err)
		assert.Equal(t, ts.URL+"/final", finalURL)
	})

	t.Run("NoRedirect", func(t *testing.T) {
		finalURL, err := d.ResolveRedirect(ctx, ts.URL+"/final")
		assert.NoError(t, err)
		assert.Equal(t, ts.URL+"/final", finalURL)
	})
}

func TestNewDownloader_TransportConfig(t *testing.T) {
	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	impl, ok := d.(*downloader)
	assert.True(t, ok, "NewDownloader should return *downloader")
	assert.NotNil(t, impl.client)
	assert.NotNil(t, impl.client.Transport)

	transport, ok := impl.client.Transport.(*http.Transport)
	assert.True(t, ok, "Transport should be *http.Transport")

	assert.Equal(t, 100, transport.MaxIdleConns)
	assert.Equal(t, 20, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 90*time.Second, transport.IdleConnTimeout)
	assert.Equal(t, 10*time.Second, transport.TLSHandshakeTimeout)
	assert.Equal(t, 1*time.Second, transport.ExpectContinueTimeout)
	assert.True(t, transport.ForceAttemptHTTP2)
}

func TestResolveRedirect_Concurrent(t *testing.T) {
	// Setup a server that redirects
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/final", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	concurrency := 50 // Exceeds MaxIdleConnsPerHost to stress it slightly
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			target := ts.URL + "/redirect"
			resolved, err := d.ResolveRedirect(ctx, target)
			assert.NoError(t, err)
			assert.Equal(t, ts.URL+"/final", resolved)
		}()
	}
	wg.Wait()
}

func TestDownloadRSSFeed_WashingtonPost(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("Skipping integration test in short mode")
	// }

	// they are "special" so they need an "official" user agent
	t.Skip()

	ctx := context.Background()
	d := NewDownloader(ctx, nil)

	u, err := url.Parse("https://www.washingtonpost.com/arcio/rss/")
	assert.NoError(t, err)

	rc, err := d.DownloadRSSFeed(ctx, u)
	assert.NoError(t, err)
	if rc != nil {
		defer func() { _ = rc.Close() }()
		content, err := io.ReadAll(rc)
		assert.NoError(t, err)
		assert.NotEmpty(t, content)
	}
}
