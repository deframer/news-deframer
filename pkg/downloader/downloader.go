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
	"strings"
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
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("HTTP request failed: %s", resp.Status)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read HTTP response: %w", err)
		}
		return string(data), nil

	case "file", "":
		if d.cfg.LocalFeedFilesDir == "" {
			return "", fmt.Errorf("unsupported scheme %s", feed.Scheme)
		}

		// Use os.OpenRoot to scope file access (Go >= 1.24)
		root, err := os.OpenRoot(d.cfg.LocalFeedFilesDir)
		if err != nil {
			return "", fmt.Errorf("failed to open local feed root: %w", err)
		}
		defer func() { _ = root.Close() }()

		// URL paths start with '/', strip it to make it relative to root
		relPath := strings.TrimPrefix(feed.Path, "/")
		f, err := root.Open(relPath)
		if err != nil {
			return "", fmt.Errorf("failed to open file %q: %w", relPath, err)
		}
		defer func() { _ = f.Close() }()

		data, err := io.ReadAll(f)
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", relPath, err)
		}
		return string(data), nil

	default:
		return "", fmt.Errorf("unsupported scheme %s", feed.Scheme)
	}
}
