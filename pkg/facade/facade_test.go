package facade

import (
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/stretchr/testify/assert"
)

type mockDownloader struct {
	downloadRSSFeed func(ctx context.Context, feed *url.URL) (io.ReadCloser, error)
}

func (m *mockDownloader) DownloadRSSFeed(ctx context.Context, feed *url.URL) (io.ReadCloser, error) {
	if m.downloadRSSFeed != nil {
		return m.downloadRSSFeed(ctx, feed)
	}
	return nil, nil
}

type mockValkey struct {
	getFeedUrl     func(u *url.URL) (*string, error)
	updateFeedUrl  func(u *url.URL, value string, ttl time.Duration) error
	tryLockFeedUrl func(u *url.URL, value string, ttl time.Duration) (bool, error)
	deleteFeedUrl  func(u *url.URL) error
	close          func() error
}

func (m *mockValkey) GetFeedUrl(u *url.URL) (*string, error) {
	if m.getFeedUrl != nil {
		return m.getFeedUrl(u)
	}
	return nil, nil
}

func (m *mockValkey) UpdateFeedUrl(u *url.URL, value string, ttl time.Duration) error {
	if m.updateFeedUrl != nil {
		return m.updateFeedUrl(u, value, ttl)
	}
	return nil
}

func (m *mockValkey) TryLockFeedUrl(u *url.URL, value string, ttl time.Duration) (bool, error) {
	if m.tryLockFeedUrl != nil {
		return m.tryLockFeedUrl(u, value, ttl)
	}
	return true, nil
}

func (m *mockValkey) DeleteFeedUrl(u *url.URL) error {
	if m.deleteFeedUrl != nil {
		return m.deleteFeedUrl(u)
	}
	return nil
}

func (m *mockValkey) Close() error {
	if m.close != nil {
		return m.close()
	}
	return nil
}

type mockRepo struct {
	getTime       func() (time.Time, error)
	findFeedByUrl func(u *url.URL) (*database.Feed, error)
}

func (m *mockRepo) GetTime() (time.Time, error) {
	if m.getTime != nil {
		return m.getTime()
	}
	return time.Now(), nil
}

func (m *mockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) {
	if m.findFeedByUrl != nil {
		return m.findFeedByUrl(u)
	}
	return nil, nil
}

func TestHasFeed(t *testing.T) {
	// Override timeouts for testing to ensure fast execution
	origMaxPendingTimeout := maxPendingTimeout
	origCheckInterval := checkInterval
	defer func() {
		maxPendingTimeout = origMaxPendingTimeout
		checkInterval = origCheckInterval
	}()
	maxPendingTimeout = 100 * time.Millisecond
	checkInterval = 10 * time.Millisecond

	ctx := context.Background()
	testURL, _ := url.Parse("http://example.com/feed.xml")

	tests := []struct {
		name          string
		setupValkey   func() *mockValkey
		setupRepo     func() *mockRepo
		expected      bool
		expectedError bool
	}{
		{
			name: "Cached Valid",
			setupValkey: func() *mockValkey {
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) {
						val := statusValid
						return &val, nil
					},
				}
			},
			setupRepo: func() *mockRepo { return &mockRepo{} },
			expected:  true,
		},
		{
			name: "Cached Invalid",
			setupValkey: func() *mockValkey {
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) {
						val := statusInvalid
						return &val, nil
					},
				}
			},
			setupRepo: func() *mockRepo { return &mockRepo{} },
			expected:  false,
		},
		{
			name: "Not Cached, DB Found",
			setupValkey: func() *mockValkey {
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) { return nil, nil },
				}
			},
			setupRepo: func() *mockRepo {
				return &mockRepo{
					findFeedByUrl: func(u *url.URL) (*database.Feed, error) {
						return &database.Feed{URL: u.String()}, nil
					},
				}
			},
			expected: true,
		},
		{
			name: "Not Cached, DB Not Found",
			setupValkey: func() *mockValkey {
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) { return nil, nil },
				}
			},
			setupRepo: func() *mockRepo {
				return &mockRepo{
					findFeedByUrl: func(u *url.URL) (*database.Feed, error) { return nil, nil },
				}
			},
			expected: false,
		},
		{
			name: "Pending then Valid",
			setupValkey: func() *mockValkey {
				calls := 0
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) {
						calls++
						if calls == 1 {
							val := statusPending
							return &val, nil
						}
						val := statusValid
						return &val, nil
					},
				}
			},
			setupRepo: func() *mockRepo { return &mockRepo{} },
			expected:  true,
		},
		{
			name: "Pending Timeout",
			setupValkey: func() *mockValkey {
				return &mockValkey{
					getFeedUrl: func(u *url.URL) (*string, error) {
						val := statusPending
						return &val, nil
					},
				}
			},
			setupRepo:     func() *mockRepo { return &mockRepo{} },
			expected:      false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(ctx, nil, tt.setupValkey(), tt.setupRepo(), &mockDownloader{})
			got, err := f.HasFeed(ctx, testURL)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

// func TestHasArticle(t *testing.T) {
// 	ctx := context.Background()

// 	f := New(ctx, nil, &mockValkey{}, &mockRepo{})

// 	u, err := url.Parse("http://example.com/article")
// 	assert.NoError(t, err)
// 	exists, err := f.HasArticle(ctx, u)
// 	assert.NoError(t, err)
// 	assert.False(t, exists)
// }

func TestGetRssProxyFeed(t *testing.T) {
	ctx := context.Background()
	rssContent := `
	<rss version="2.0">
		<channel>
			<title>Test Feed</title>
		</channel>
	</rss>`
	mockDL := &mockDownloader{
		downloadRSSFeed: func(ctx context.Context, feed *url.URL) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(rssContent)), nil
		},
	}
	f := New(ctx, nil, &mockValkey{}, &mockRepo{}, mockDL)

	filter := RSSProxyFilter{
		URL:      "http://example.com",
		Lang:     "en",
		MaxScore: 0.75,
		Embedded: true,
	}

	xmlData, err := f.GetRssProxyFeed(ctx, &filter)
	assert.NoError(t, err)
	assert.Contains(t, xmlData, "<title>Proxied: Test Feed</title>")
}
