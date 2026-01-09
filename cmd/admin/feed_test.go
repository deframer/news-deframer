package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/database"
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

	testURL := "http://example.com/rss"

	// 1. Create Feed
	out := captureOutput(func() {
		addFeed(testURL, true)
	})
	assert.Contains(t, out, "Added feed")
	assert.Contains(t, out, testURL)

	// Verify in Mock
	u, _ := url.Parse(testURL)
	feed, err := mock.FindFeedByUrl(u)
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	assert.True(t, feed.Enabled)
	feedID := feed.ID

	// 2. List Feeds (JSON)
	out = captureOutput(func() {
		listFeeds(true)
	})
	var feeds []database.Feed
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.Equal(t, testURL, feeds[0].URL)
	assert.True(t, feeds[0].Enabled)

	// 3. Disable Feed
	out = captureOutput(func() {
		disableFeed(testURL)
	})
	assert.Contains(t, out, "Disabled feed")

	// Verify in Mock
	feed, _ = mock.FindFeedById(feedID)
	assert.False(t, feed.Enabled)

	// 4. List Feeds (JSON) check disabled
	out = captureOutput(func() {
		listFeeds(true)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.False(t, feeds[0].Enabled)

	// 5. Enable Feed
	out = captureOutput(func() {
		enableFeed(testURL)
	})
	assert.Contains(t, out, "Enabled feed")

	// Verify in Mock
	feed, _ = mock.FindFeedById(feedID)
	assert.True(t, feed.Enabled)

	// 6. Delete Feed
	out = captureOutput(func() {
		deleteFeed(testURL)
	})
	assert.Contains(t, out, "Deleted feed")

	// Verify in Mock
	feed, _ = mock.FindFeedById(feedID)
	assert.True(t, feed.DeletedAt.Valid)

	// 7. List Feeds (JSON) check deleted
	out = captureOutput(func() {
		listFeeds(true)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 1)
	assert.True(t, feeds[0].DeletedAt.Valid)
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

func (m *MockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	return nil, nil
}

func (m *MockRepo) FindItemsByCachedFeedFeedId(feedID uuid.UUID) ([]database.Item, error) {
	return nil, nil
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
