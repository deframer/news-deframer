package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
)

// fake a real chrome browser

// https://github.com/brave/brave-browser/wiki/User-Agents
const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
const defaultSec_CH_CA = `"Chromium";v="136", "Brave";v="136", "Not.A/Brand";v="99"`
const defaultPriority = "priority: u=0, i"

type DownloaderOpts struct {
	UserAgent string
	Sec_CH_UA string
	Priority  string
}

type downloader struct {
	ctx       context.Context
	cfg       *config.Config
	logger    *slog.Logger
	client    *http.Client
	userAgent string
	sec_CH_UA string
	priority  string
}

type Downloader interface {
	DownloadRSSFeed(ctx context.Context, feed *url.URL) (io.ReadCloser, error)
	ResolveRedirect(ctx context.Context, targetURL string) (string, error)
}

// NewDownloader initializes a new downloader
func NewDownloader(ctx context.Context, cfg *config.Config, opts ...DownloaderOpts) Downloader {
	// Custom transport to allow more idle connections per host
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   20, // Optimized for concurrent feed item resolution
	}

	userAgent := defaultUserAgent
	if len(opts) > 0 && opts[0].UserAgent != "" {
		userAgent = opts[0].UserAgent
	}

	sec_CH_CA := defaultSec_CH_CA
	if len(opts) > 0 && opts[0].Sec_CH_UA != "" {
		sec_CH_CA = opts[0].Sec_CH_UA
	}

	priority := defaultPriority
	if len(opts) > 0 && opts[0].Priority != "" {
		priority = opts[0].Priority
	}

	return &downloader{
		ctx:    ctx,
		cfg:    cfg,
		logger: slog.With("component", "downloader"),
		client: &http.Client{
			Timeout:   15 * time.Second,
			Transport: transport,
		},
		userAgent: userAgent,
		sec_CH_UA: sec_CH_CA,
		priority:  priority,
	}
}

// DownloadRSSFeed downloads from http/https URLs
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

		// Mimic a browser to avoid some anti-bot protections
		req.Header.Set("User-Agent", d.userAgent)
		req.Header.Set("Sec-CH-UA", d.sec_CH_UA)
		req.Header.Set("Priority", d.priority)

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
	req.Header.Set("User-Agent", d.userAgent)
	req.Header.Set("Sec-CH-UA", d.sec_CH_UA)
	req.Header.Set("Priority", d.priority)

	resp, err := d.client.Do(req)
	if err != nil {
		return targetURL, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.Request.URL.String(), nil
}
