package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
)

const userAgent = "Mozilla/5.0 (compatible; Deframer/1.0; +https://github.com/deframer/news-deframer)"

type downloader struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
	client *http.Client
}

type Downloader interface {
	DownloadRSSFeed(ctx context.Context, feed *url.URL) (io.ReadCloser, error)
	ResolveRedirect(ctx context.Context, targetURL string) (string, error)
}

// NewDownloader initializes a new downloader
func NewDownloader(ctx context.Context, cfg *config.Config) Downloader {
	return &downloader{
		ctx:    ctx,
		cfg:    cfg,
		logger: slog.With("component", "downloader"),
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// DownloadRSSFeed downloads from http/https/file URLs
func (d *downloader) DownloadRSSFeed(ctx context.Context, feed *url.URL) (io.ReadCloser, error) {
	if feed == nil {
		return nil, errors.New("feed cannot be nil")
	}

	d.logger.DebugContext(ctx, "downloading feed", "url", feed.String())

	switch feed.Scheme {
	case "http", "https":
		// HTTP download
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for URL %q: %w", feed.String(), err)
		}
		req.Header.Set("User-Agent", userAgent)

		resp, err := d.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL %q: %w", feed.String(), err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("HTTP request failed: %s", resp.Status)
		}

		return resp.Body, nil

	default:
		return nil, fmt.Errorf("unsupported scheme %s", feed.Scheme)
	}
}

// ResolveRedirect performs a HEAD request to resolve the final URL after redirects.
func (d *downloader) ResolveRedirect(ctx context.Context, targetURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		return targetURL, err
	}

	// Mimic a browser to avoid some anti-bot protections
	req.Header.Set("User-Agent", userAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return targetURL, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.Request.URL.String(), nil
}
