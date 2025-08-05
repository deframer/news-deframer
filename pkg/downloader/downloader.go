package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type downloader struct {
}

type Downloader interface {
	DownloadRSSFeed(feed string) (string, error)
}

// NewDownloader initializes a new downloader
func NewDownloader() Downloader {
	res := &downloader{}

	return res
}

func (d *downloader) DownloadRSSFeed(feed string) (string, error) {
	if feed == "" {
		return "", errors.New("feed cannot be empty")
	}

	switch {
	case strings.HasPrefix(feed, "http://") || strings.HasPrefix(feed, "https://"):
		// HTTP download
		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Get(feed)
		if err != nil {
			return "", fmt.Errorf("failed to fetch URL %q: %w", feed, err)
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
		// Local file handling (with or without file:// prefix)
		path := strings.TrimPrefix(feed, "file://")
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read file %q: %w", path, err)
		}
		return string(data), nil
	}
}
