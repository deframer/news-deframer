package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
)

type downloader struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
	client *http.Client
}

type Downloader interface {
	DownloadRSSFeed(ctx context.Context, feed *url.URL) (string, error)
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
func (d *downloader) DownloadRSSFeed(ctx context.Context, feed *url.URL) (string, error) {
	if feed == nil {
		return "", errors.New("feed cannot be nil")
	}

	d.logger.DebugContext(ctx, "downloading feed", "url", feed.String())

	switch feed.Scheme {
	case "http", "https":
		// HTTP download
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.String(), nil)
		if err != nil {
			return "", fmt.Errorf("failed to create request for URL %q: %w", feed.String(), err)
		}

		resp, err := d.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL %q: %w", feed.String(), err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP request failed: %s", resp.Status)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read HTTP response: %w", err)
		}
		return string(data), nil

	default:
		if !d.cfg.LocalFileFeeds {
			return "", fmt.Errorf("unsupported scheme %s", feed.Scheme)
		}
		// Local file handling (with or without file:// prefix)
		path := feed.Path
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", path, err)
		}
		return string(data), nil
	}
}
