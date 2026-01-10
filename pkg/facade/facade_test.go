package facade

import (
	"context"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/valkey"
	"github.com/google/uuid"
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
	getFeedByUrl     func(u *url.URL) (*valkey.FeedUrlToUUID, error)
	updateFeed       func(state valkey.FeedUrlToUUID, info valkey.FeedInfo, ttl time.Duration) error
	tryLockFeedByUrl func(u *url.URL, state valkey.FeedUrlToUUID, ttl time.Duration) (bool, error)
	drainFeed        func(feedID uuid.UUID) error
	close            func() error
}

func (m *mockValkey) GetFeedByUrl(u *url.URL) (*valkey.FeedUrlToUUID, error) {
	if m.getFeedByUrl != nil {
		return m.getFeedByUrl(u)
	}
	return nil, nil
}

func (m *mockValkey) UpdateFeed(state valkey.FeedUrlToUUID, info valkey.FeedInfo, ttl time.Duration) error {
	if m.updateFeed != nil {
		return m.updateFeed(state, info, ttl)
	}
	return nil
}

func (m *mockValkey) TryLockFeedByUrl(u *url.URL, state valkey.FeedUrlToUUID, ttl time.Duration) (bool, error) {
	if m.tryLockFeedByUrl != nil {
		return m.tryLockFeedByUrl(u, state, ttl)
	}
	return true, nil
}

func (m *mockValkey) DrainFeed(feedID uuid.UUID) error {
	if m.drainFeed != nil {
		return m.drainFeed(feedID)
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
	findFeedByUrl                func(u *url.URL) (*database.Feed, error)
	findFeedByUrlAndAvailability func(u *url.URL, onlyEnabled bool) (*database.Feed, error)
	findFeedById                 func(feedID uuid.UUID) (*database.Feed, error)
	upsertFeed                   func(feed *database.Feed) error
	findCachedFeedById           func(feedID uuid.UUID) (*database.CachedFeed, error)
	findItemsByCachedFeedFeedId  func(feedID uuid.UUID) ([]database.Item, error)
	findItemsByUrl               func(u *url.URL) ([]database.Item, error)
	getAllFeeds                  func(deleted bool) ([]database.Feed, error)
	deleteFeedById               func(id uuid.UUID) error
}

func (m *mockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) {
	if m.findFeedByUrl != nil {
		return m.findFeedByUrl(u)
	}
	return nil, nil
}

func (m *mockRepo) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*database.Feed, error) {
	if m.findFeedByUrlAndAvailability != nil {
		return m.findFeedByUrlAndAvailability(u, onlyEnabled)
	}
	return nil, nil
}

func (m *mockRepo) FindFeedById(feedID uuid.UUID) (*database.Feed, error) {
	if m.findFeedById != nil {
		return m.findFeedById(feedID)
	}
	return nil, nil
}

func (m *mockRepo) UpsertFeed(feed *database.Feed) error {
	if m.upsertFeed != nil {
		return m.upsertFeed(feed)
	}
	return nil
}

func (m *mockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	if m.findCachedFeedById != nil {
		return m.findCachedFeedById(feedID)
	}
	return nil, nil
}

func (m *mockRepo) FindItemsByCachedFeedFeedId(feedID uuid.UUID) ([]database.Item, error) {
	if m.findItemsByCachedFeedFeedId != nil {
		return m.findItemsByCachedFeedFeedId(feedID)
	}
	return nil, nil
}

func (m *mockRepo) FindItemsByUrl(u *url.URL) ([]database.Item, error) {
	if m.findItemsByUrl != nil {
		return m.findItemsByUrl(u)
	}
	return nil, nil
}

func (m *mockRepo) GetAllFeeds(deleted bool) ([]database.Feed, error) {
	if m.getAllFeeds != nil {
		return m.getAllFeeds(deleted)
	}
	return nil, nil
}

func (m *mockRepo) DeleteFeedById(id uuid.UUID) error {
	if m.deleteFeedById != nil {
		return m.deleteFeedById(id)
	}
	return nil
}

func TestHasFeedByUrl(t *testing.T) {
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) {
						return &valkey.FeedUrlToUUID{
							Cache: valkey.Ok,
							UUID:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						}, nil
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) {
						return &valkey.FeedUrlToUUID{
							Cache: valkey.ValueUnknown,
						}, nil
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) { return nil, nil },
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) { return nil, nil },
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) {
						calls++
						if calls == 1 {
							return &valkey.FeedUrlToUUID{
								Cache: valkey.Updating,
							}, nil
						}
						return &valkey.FeedUrlToUUID{
							Cache: valkey.Ok,
							UUID:  uuid.MustParse("00000000-0000-0000-0000-000000000001"),
						}, nil
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
					getFeedByUrl: func(u *url.URL) (*valkey.FeedUrlToUUID, error) {
						return &valkey.FeedUrlToUUID{
							Cache: valkey.Updating,
						}, nil
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
			got, err := f.HasFeedByUrl(ctx, testURL)
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
