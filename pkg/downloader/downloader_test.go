package downloader

import (
	"context"
	"io"
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

func TestDownloadRSSFeed_File(t *testing.T) {
	// Create temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.xml")
	err := os.WriteFile(tmpFile, []byte("file content"), 0644)
	assert.NoError(t, err)

	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{
		LocalFeedFilesDir: tmpDir,
	})

	// Test with file path (relative to root)
	u1, err := url.Parse("test.xml")
	assert.NoError(t, err)
	rc1, err := d.DownloadRSSFeed(ctx, u1)
	assert.NoError(t, err)
	defer func() { _ = rc1.Close() }()
	content1, err := io.ReadAll(rc1)
	assert.NoError(t, err)
	assert.Equal(t, "file content", string(content1))

	// Test with file:// prefix (path /test.xml relative to root)
	u2, err := url.Parse("file:///test.xml")
	assert.NoError(t, err)
	rc2, err := d.DownloadRSSFeed(ctx, u2)
	assert.NoError(t, err)
	defer func() { _ = rc2.Close() }()
	content2, err := io.ReadAll(rc2)
	assert.NoError(t, err)
	assert.Equal(t, "file content", string(content2))
}

func TestDownloadRSSFeed_File_Disabled(t *testing.T) {
	ctx := context.Background()
	d := NewDownloader(ctx, &config.Config{
		LocalFeedFilesDir: "",
	})

	// Even if the file existed, it should fail due to config
	u, err := url.Parse("file:///tmp/test.xml")
	assert.NoError(t, err)
	_, err = d.DownloadRSSFeed(ctx, u)
	assert.Error(t, err)
}
