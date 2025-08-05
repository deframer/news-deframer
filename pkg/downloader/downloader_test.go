package downloader

//go:generate mockgen -destination=./mock_downloader/mocks.go github.com/egandro/news-deframer/pkg/downloader Downloader

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDownloader(t *testing.T) {
	d := NewDownloader()
	assert.NotNil(t, d, "Downloader should be initialized")
}

func TestNewUpdateFeeds(t *testing.T) {
	d := NewDownloader()
	assert.NotNil(t, d, "Downloader should be initialized")

	tests := []struct {
		name        string
		feed        string
		expected    string
		expectError bool
	}{
		// HTTP feed
		func() struct {
			name        string
			feed        string
			expected    string
			expectError bool
		} {
			expectedHTTP := "<rss>http feed</rss>"
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(expectedHTTP))
			}))
			t.Cleanup(ts.Close)
			return struct {
				name        string
				feed        string
				expected    string
				expectError bool
			}{
				name:     "HTTP feed",
				feed:     ts.URL,
				expected: expectedHTTP,
			}
		}(),
		// file:// feed
		func() struct {
			name        string
			feed        string
			expected    string
			expectError bool
		} {
			expectedFile := "<rss>file feed</rss>"
			tmpFile := filepath.Join(t.TempDir(), "file1.xml")
			assert.NoError(t, os.WriteFile(tmpFile, []byte(expectedFile), 0644))
			return struct {
				name        string
				feed        string
				expected    string
				expectError bool
			}{
				name:     "file:// feed",
				feed:     "file://" + tmpFile,
				expected: expectedFile,
			}
		}(),
		// plain file path feed
		func() struct {
			name        string
			feed        string
			expected    string
			expectError bool
		} {
			expectedPlain := "<rss>plain file</rss>"
			tmpFile := filepath.Join(t.TempDir(), "file2.xml")
			assert.NoError(t, os.WriteFile(tmpFile, []byte(expectedPlain), 0644))
			return struct {
				name        string
				feed        string
				expected    string
				expectError bool
			}{
				name:     "plain file path feed",
				feed:     tmpFile,
				expected: expectedPlain,
			}
		}(),
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := d.DownloadRSSFeed(tc.feed)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}
