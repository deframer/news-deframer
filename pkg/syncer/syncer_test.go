package syncer

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	enqueueSyncCalled bool
	lastId            uuid.UUID
	removeSyncCalled  bool
}

// Implement database.Repository interface stubs
func (m *mockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) { return nil, nil }
func (m *mockRepo) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*database.Feed, error) {
	return nil, nil
}
func (m *mockRepo) FindFeedById(feedID uuid.UUID) (*database.Feed, error) { return nil, nil }
func (m *mockRepo) UpsertFeed(feed *database.Feed) error                  { return nil }
func (m *mockRepo) FindItemsByUrl(u *url.URL) ([]database.Item, error)    { return nil, nil }
func (m *mockRepo) GetAllFeeds(deleted bool) ([]database.Feed, error)     { return nil, nil }
func (m *mockRepo) DeleteFeedById(id uuid.UUID) error                     { return nil }
func (m *mockRepo) RemoveSync(id uuid.UUID) error {
	m.removeSyncCalled = true
	m.lastId = id
	return nil
}
func (m *mockRepo) BeginFeedUpdate(lockDuration time.Duration) (*database.Feed, error) {
	return nil, nil
}
func (m *mockRepo) EndFeedUpdate(id uuid.UUID, err error, successDelay time.Duration) error {
	return nil
}
func (m *mockRepo) GetPendingHashes(feedID uuid.UUID, hashes []string) (map[string]bool, error) {
	res := make(map[string]bool, len(hashes))
	for _, h := range hashes {
		res[h] = true
	}
	return res, nil
}
func (m *mockRepo) UpsertItem(item *database.Item) error { return nil }
func (m *mockRepo) GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]database.Item, error) {
	return nil, nil
}
func (m *mockRepo) UpsertCachedFeed(cachedFeed *database.CachedFeed) error { return nil }
func (m *mockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	return nil, nil
}
func (m *mockRepo) FindFeedScheduleById(feedID uuid.UUID) (*database.FeedSchedule, error) {
	return nil, nil
}

// Implement EnqueueSync
func (m *mockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error {
	m.enqueueSyncCalled = true
	m.lastId = id
	return nil
}

func TestSyncFeed(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	id := uuid.New()
	err = s.SyncFeed(id)

	assert.NoError(t, err)
	assert.True(t, repo.enqueueSyncCalled)
	assert.Equal(t, id, repo.lastId)
}

func TestPoll(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to test exit

	s, err := New(ctx, cfg, repo)
	assert.NoError(t, err)
	s.Poll()
}

func TestStopPolling(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	id := uuid.New()
	err = s.StopPolling(id)

	assert.NoError(t, err)
	assert.True(t, repo.removeSyncCalled)
	assert.Equal(t, id, repo.lastId)
}

func TestSyncFeedInternal(t *testing.T) {
	// developer trigger
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping test on GitHub Actions")
	}

	t.Skip("Skipping test")

	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := database.NewRepository(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		url     string
		wantErr bool
		enabled bool
	}{
		{
			name:    "Localhost Feed",
			url:     "http://localhost:8003/feed",
			wantErr: false,
			enabled: true,
		},
		{
			name:    "WordPress Feed",
			url:     "http://wordpress/feed",
			wantErr: true,
			enabled: false,
		},
	}

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.enabled {
				t.Skip("Test case disabled")
			}

			u, err := url.Parse(tt.url)
			assert.NoError(t, err)

			feed, err := repo.FindFeedByUrl(u)
			assert.NoError(t, err)
			if assert.NotNil(t, feed) {
				err = s.updatingFeed(feed)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestWantedDomains(t *testing.T) {
	repo := &mockRepo{}
	cfg := &config.Config{}
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		feed     *database.Feed
		expected []string
		wantErr  bool
	}{
		{
			name: "Enforce Disabled",
			feed: &database.Feed{
				EnforceFeedDomain: false,
				URL:               "http://example.com/feed",
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Simple Domain",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://example.com/feed",
			},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Subdomain",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://blog.example.com/feed",
			},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Public Suffix",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://example.co.uk/feed",
			},
			expected: []string{"example.co.uk"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Localhost",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://localhost:8003/feed",
			},
			expected: []string{"localhost"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - WordPress",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://wordpress/feed",
			},
			expected: []string{"wordpress"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains, err := s.wantedDomains(tt.feed)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, domains)
			}
		})
	}
}
