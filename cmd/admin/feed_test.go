package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/syncer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestParseAndNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "Simple URL",
			input: "http://example.com",
			want:  "http://example.com",
		},
		{
			name:  "Trailing Slash",
			input: "http://example.com/",
			want:  "http://example.com",
		},
		{
			name:  "Whitespace",
			input: "  http://example.com  ",
			want:  "http://example.com",
		},
		{
			name:  "Whitespace and Slash",
			input: "  http://example.com/  ",
			want:  "http://example.com",
		},
		{
			name:    "Invalid URL",
			input:   "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAndNormalizeURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestFeedCommands(t *testing.T) {
	// Setup Mock
	mock := NewMockRepo()
	repo = mock
	var err error
	feedSyncer, err = syncer.New(context.Background(), &config.Config{}, mock)
	assert.NoError(t, err)

	testURL := "http://example.com/rss"

	// 1. Create Feed
	out := captureOutput(func() {
		addFeed(testURL, true, false)
	})
	assert.Contains(t, out, "Added feed")
	assert.Contains(t, out, testURL)

	// 2. Sync Feed (Create Schedule)
	out = captureOutput(func() {
		syncFeed(testURL)
	})
	assert.Contains(t, out, "Triggered sync")

	// 3. List Feeds (JSON) - Verify Schedule Present
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	var feeds []database.Feed
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.Equal(t, testURL, feeds[0].URL)
	assert.True(t, feeds[0].Enabled)
	assert.NotNil(t, feeds[0].FeedSchedule)

	// 4. Disable Feed (Should Remove Schedule)
	out = captureOutput(func() {
		disableFeed(testURL)
	})
	assert.Contains(t, out, "Disabled feed")

	// 5. List Feeds (JSON) - Verify Schedule Gone
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.False(t, feeds[0].Enabled)
	if feeds[0].FeedSchedule != nil {
		assert.Nil(t, feeds[0].FeedSchedule.NextRunAt)
	}

	// 6. Enable Feed
	out = captureOutput(func() {
		enableFeed(testURL)
	})
	assert.Contains(t, out, "Enabled feed")

	// 7. Set Polling true
	out = captureOutput(func() {
		setPolling(testURL, "true")
	})
	assert.Contains(t, out, "Set polling to true")

	// 8. Verify Schedule Created automatically
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.NotNil(t, feeds[0].FeedSchedule)

	// 9. Set Polling false (Should Remove Schedule)
	out = captureOutput(func() {
		setPolling(testURL, "false")
	})
	assert.Contains(t, out, "Set polling to false")

	// 10. List Feeds (Verify Schedule Gone)
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	if feeds[0].FeedSchedule != nil {
		assert.Nil(t, feeds[0].FeedSchedule.NextRunAt)
	}

	// 11. Sync Feed (Create Schedule for Delete test)
	captureOutput(func() { syncFeed(testURL) })

	// 12. Delete Feed (Should Remove Schedule)
	out = captureOutput(func() {
		deleteFeed(testURL)
	})
	assert.Contains(t, out, "Deleted feed")

	// 13. List Feeds (Verify Deleted and Schedule Gone)
	out = captureOutput(func() {
		listFeeds(true, true)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.True(t, feeds[0].DeletedAt.Valid)
	if feeds[0].FeedSchedule != nil {
		assert.Nil(t, feeds[0].FeedSchedule.NextRunAt)
	}

	// 14. Test Enable with Polling triggers Sync
	testURL2 := "http://example.com/rss2"
	captureOutput(func() {
		addFeed(testURL2, false, true) // Add disabled feed with polling=true
	})

	out = captureOutput(func() {
		enableFeed(testURL2)
	})
	assert.Contains(t, out, "Enabled feed")

	// Verify Schedule was created
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.Equal(t, testURL2, feeds[0].URL)
	assert.NotNil(t, feeds[0].FeedSchedule)
}

// --- Mock Repository ---

type MockRepo struct {
	feeds map[uuid.UUID]*database.Feed
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		feeds: make(map[uuid.UUID]*database.Feed),
	}
}

func (m *MockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) {
	return m.FindFeedByUrlAndAvailability(u, true)
}

func (m *MockRepo) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*database.Feed, error) {
	for _, f := range m.feeds {
		if f.URL == u.String() {
			if f.DeletedAt.Valid {
				continue
			}
			if onlyEnabled && !f.Enabled {
				continue
			}
			return f, nil
		}
	}
	return nil, nil
}

func (m *MockRepo) FindFeedById(feedID uuid.UUID) (*database.Feed, error) {
	if f, ok := m.feeds[feedID]; ok {
		return f, nil
	}
	return nil, nil
}

func (m *MockRepo) UpsertFeed(feed *database.Feed) error {
	if feed.ID == uuid.Nil {
		feed.ID = uuid.New()
		feed.CreatedAt = time.Now()
	}
	feed.UpdatedAt = time.Now()
	m.feeds[feed.ID] = feed
	return nil
}

func (m *MockRepo) UpsertItem(item *database.Item) error {
	return nil
}

func (m *MockRepo) FindItemsByUrl(u *url.URL) ([]database.Item, error) {
	return nil, nil
}

func (m *MockRepo) GetAllFeeds(deleted bool) ([]database.Feed, error) {
	var feeds []database.Feed
	for _, f := range m.feeds {
		if !deleted && f.DeletedAt.Valid {
			continue
		}
		feeds = append(feeds, *f)
	}
	return feeds, nil
}

func (m *MockRepo) DeleteFeedById(id uuid.UUID) error {
	if f, ok := m.feeds[id]; ok {
		f.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		f.UpdatedAt = time.Now()
	}
	return nil
}

func (m *MockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error {
	if f, ok := m.feeds[id]; ok {
		now := time.Now()
		f.FeedSchedule = &database.FeedSchedule{
			ID:        id,
			NextRunAt: &now,
		}
	}
	return nil
}

func (m *MockRepo) RemoveSync(id uuid.UUID) error {
	if f, ok := m.feeds[id]; ok && f.FeedSchedule != nil {
		f.FeedSchedule.NextRunAt = nil
		f.FeedSchedule.LockedUntil = nil
		f.FeedSchedule.LastError = nil
	}
	return nil
}

func (m *MockRepo) BeginFeedUpdate(lockDuration time.Duration) (*database.Feed, error) {
	return nil, nil
}

func (m *MockRepo) EndFeedUpdate(id uuid.UUID, err error, successDelay time.Duration) error {
	return nil
}

func (m *MockRepo) GetPendingHashes(feedID uuid.UUID, hashes []string) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, h := range hashes {
		res[h] = true
	}
	return res, nil
}

func (m *MockRepo) GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]database.Item, error) {
	return nil, nil
}

func (m *MockRepo) UpsertCachedFeed(cachedFeed *database.CachedFeed) error {
	return nil
}

func (m *MockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	return nil, nil
}

func (m *MockRepo) FindFeedScheduleById(feedID uuid.UUID) (*database.FeedSchedule, error) {
	if f, ok := m.feeds[feedID]; ok {
		return f.FeedSchedule, nil
	}
	return nil, nil
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}
