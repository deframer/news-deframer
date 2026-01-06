package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestDownloadRSSFeed_HTTP(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock content"))
	}))
	defer server.Close()

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{})

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)
	content, err := d.DownloadRSSFeed(ctx, u)
	assert.NoError(t, err)
	assert.Equal(t, "mock content", content)
}

func TestDownloadRSSFeed_File(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.xml")
	err := os.WriteFile(tmpFile, []byte("file content"), 0644)
	assert.NoError(t, err)

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{
		LocalFileFeeds: true,
	})

	// Test with file path
	u1, err := url.Parse(tmpFile)
	assert.NoError(t, err)
	content, err := d.DownloadRSSFeed(ctx, u1)
	assert.NoError(t, err)
	assert.Equal(t, "file content", content)

	// Test with file:// prefix
	u2, err := url.Parse("file://" + tmpFile)
	assert.NoError(t, err)
	content, err = d.DownloadRSSFeed(ctx, u2)
	assert.NoError(t, err)
	assert.Equal(t, "file content", content)
}

func TestDownloadRSSFeed_File_Disabled(t *testing.T) {
	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{
		LocalFileFeeds: false,
	})

	// Even if the file existed, it should fail due to config
	u, err := url.Parse("file:///tmp/test.xml")
	assert.NoError(t, err)
	_, err = d.DownloadRSSFeed(ctx, u)
	assert.Error(t, err)
}
