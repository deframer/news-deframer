package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"goa.design/clue/log"
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
		ctx: ctx,
		cfg: cfg,
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

	log.Printf(ctx, "downloading feed url=%s", feed.String())

	switch feed.Scheme {
	case "http", "https":
		// HTTP download
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for URL %q: %w", feed.String(), err)
		}

		// Mimic a browser to avoid some anti-bot protections
		req.Header.Set("User-Agent", d.userAgent)
		req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8, */*;q=0.1")
		req.Header.Set("Sec-CH-UA", d.sec_CH_UA)
		req.Header.Set("Priority", d.priority)

		// #nosec G704 -- feed URL is an explicit user-configured target.
		resp, err := d.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch URL %q: %w", feed.String(), err)
		}

		if resp.StatusCode != http.StatusOK {
			// we have very stupid rss feeds
			// we get a 404 - but a RSS feed data
			// we don't even see this in the browser - this is a workaround
			peek, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
			if err == nil && looksLikeRSSFeed(peek) {
				return &responseReadCloser{
					Reader: io.MultiReader(bytes.NewReader(peek), resp.Body),
					Closer: resp.Body,
				}, nil
			}

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
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8, */*;q=0.1")
	req.Header.Set("Sec-CH-UA", d.sec_CH_UA)
	req.Header.Set("Priority", d.priority)

	// #nosec G704 -- redirect resolution intentionally requests caller-provided URL.
	resp, err := d.client.Do(req)
	if err != nil {
		return targetURL, err
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.Request.URL.String(), nil
}

type responseReadCloser struct {
	io.Reader
	Closer io.Closer
}

func (r *responseReadCloser) Close() error {
	return r.Closer.Close()
}

func looksLikeRSSFeed(peek []byte) bool {
	trimmed := bytes.TrimSpace(bytes.ToLower(peek))
	return bytes.HasPrefix(trimmed, []byte("<?xml")) || bytes.HasPrefix(trimmed, []byte("<rss")) || bytes.HasPrefix(trimmed, []byte("<feed"))
}
