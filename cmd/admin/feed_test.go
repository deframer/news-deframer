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

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/syncer"
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
		addFeed(testURL, true, false, false, false, "", false, []string{})
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
	assert.NotNil(t, feeds[0].RootDomain)
	assert.Equal(t, "example.com", *feeds[0].RootDomain)

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
		assert.Nil(t, feeds[0].FeedSchedule.NextThinkerAt)
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
		assert.Nil(t, feeds[0].FeedSchedule.NextThinkerAt)
	}

	// 11. Sync Feed (Create Schedule for Delete test)
	captureOutput(func() { syncFeed(testURL) })

	// 12. Delete Feed (Should Remove Schedule)
	out = captureOutput(func() {
		deleteFeed(testURL, false)
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
		assert.Nil(t, feeds[0].FeedSchedule.NextThinkerAt)
	}

	// 13a. Purge Feed (Must resolve deleted feeds too)
	out = captureOutput(func() {
		id := feeds[0].ID
		deleteFeed(id.String(), true)
	})
	assert.Contains(t, out, "Purged feed")

	// 13b. List Feeds (Verify Gone)
	out = captureOutput(func() {
		listFeeds(true, true)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)
	assert.Len(t, feeds, 0)

	// 14. Test Enable with Polling triggers Sync
	testURL2 := "http://example.com/rss2"
	captureOutput(func() {
		addFeed(testURL2, false, true, false, false, "", false, []string{}) // Add disabled feed with polling=true
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

	// 15. Test --no-root-domain
	testURL3 := "http://no-root.com/rss"
	captureOutput(func() {
		addFeed(testURL3, true, false, false, true, "", false, []string{})
	})

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundNoRoot *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL3 {
			foundNoRoot = &feeds[i]
		}
	}
	assert.NotNil(t, foundNoRoot)
	assert.Nil(t, foundNoRoot.RootDomain)

	// 16. Test Root Domain Extraction (Subdomain)
	testURL4 := "http://blog.example.co.uk/rss"
	captureOutput(func() {
		addFeed(testURL4, true, false, false, false, "", false, []string{})
	})

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundSub *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL4 {
			foundSub = &feeds[i]
		}
	}
	assert.NotNil(t, foundSub)
	assert.NotNil(t, foundSub.RootDomain)
	assert.Equal(t, "example.co.uk", *foundSub.RootDomain)

	// 17. Test Language Commands
	testURL5 := "http://example.com/rss5"
	captureOutput(func() {
		addFeed(testURL5, true, false, false, false, "en", false, []string{})
	})

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundLang *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL5 {
			foundLang = &feeds[i]
		}
	}
	assert.NotNil(t, foundLang)
	assert.NotNil(t, foundLang.Language)
	assert.Equal(t, "en", *foundLang.Language)

	out = captureOutput(func() {
		setLanguage(testURL5, "de")
	})
	assert.Contains(t, out, "Set language to de")

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	for i := range feeds {
		if feeds[i].URL == testURL5 {
			foundLang = &feeds[i]
		}
	}
	assert.NotNil(t, foundLang)
	assert.NotNil(t, foundLang.Language)
	assert.Equal(t, "de", *foundLang.Language)

	out = captureOutput(func() {
		deleteLanguage(testURL5)
	})
	assert.Contains(t, out, "Deleted language")

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	for i := range feeds {
		if feeds[i].URL == testURL5 {
			foundLang = &feeds[i]
		}
	}
	assert.NotNil(t, foundLang)
	assert.Nil(t, foundLang.Language)

	// 17. Sync All (should trigger sync for all enabled feeds)
	// Reset schedule for testURL2
	if f, err := mock.FindFeedByUrlAndAvailability(&url.URL{Scheme: "http", Host: "example.com", Path: "/rss2"}, true); err == nil && f != nil {
		f.FeedSchedule.NextThinkerAt = nil
	}

	out = captureOutput(func() {
		syncAllFeeds()
	})
	assert.Contains(t, out, "Triggered sync for url=http://example.com/rss2")
	// Deleted/Disabled feeds should not be synced. testURL is deleted.
	// We add " with" to ensure we don't match partial URLs (like rss vs rss2)
	assert.NotContains(t, out, "Triggered sync for url="+testURL+" with")

	// 18. Test ResolveItemUrl Commands
	testURL6 := "http://example.com/rss6"
	captureOutput(func() {
		addFeed(testURL6, true, false, false, false, "", true, []string{})
	})

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundResolve *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL6 {
			foundResolve = &feeds[i]
		}
	}
	assert.NotNil(t, foundResolve)
	assert.True(t, foundResolve.ResolveItemUrl)

	out = captureOutput(func() {
		setResolveItemUrl(testURL6, "false")
	})
	assert.Contains(t, out, "Set resolve_item_url to false")

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	for i := range feeds {
		if feeds[i].URL == testURL6 {
			foundResolve = &feeds[i]
		}
	}
	assert.NotNil(t, foundResolve)
	assert.False(t, foundResolve.ResolveItemUrl)

	// 19. Test Mining Commands
	testURL7 := "http://example.com/rss7"
	captureOutput(func() {
		addFeed(testURL7, true, false, true, false, "", false, []string{})
	})

	out = captureOutput(func() {
		mineFeed(testURL7)
	})
	assert.Contains(t, out, "Triggered mining for url=http://example.com/rss7")

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundMine *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL7 {
			foundMine = &feeds[i]
		}
	}
	assert.NotNil(t, foundMine)
	assert.True(t, foundMine.Mining)
	assert.NotNil(t, foundMine.FeedSchedule)
	assert.NotNil(t, foundMine.FeedSchedule.NextMiningAt)

	out = captureOutput(func() {
		setMining(testURL7, "false")
	})
	assert.Contains(t, out, "Set mining to false")

	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	for i := range feeds {
		if feeds[i].URL == testURL7 {
			foundMine = &feeds[i]
		}
	}
	assert.NotNil(t, foundMine)
	assert.False(t, foundMine.Mining)

	// 20. Test mine-all
	// Reset schedule for testURL7
	if f, err := mock.FindFeedByUrlAndAvailability(&url.URL{Scheme: "http", Host: "example.com", Path: "/rss7"}, true); err == nil && f != nil {
		f.FeedSchedule.NextMiningAt = nil
		setMining(f.URL, "true")
	}

	// Add another feed that is enabled but not for mining
	testURL8 := "http://example.com/rss8"
	captureOutput(func() {
		addFeed(testURL8, true, false, false, false, "", false, []string{})
	})

	out = captureOutput(func() {
		mineAllFeeds()
	})
	assert.Contains(t, out, "Triggered mining for url=http://example.com/rss7")
	assert.NotContains(t, out, "Triggered mining for url=http://example.com/rss8")

	// 21. Test sync on non-polling feed
	testURL9 := "http://example.com/rss9"
	captureOutput(func() {
		addFeed(testURL9, true, false, false, false, "", false, []string{}) // Polling is false
	})

	out = captureOutput(func() {
		syncFeed(testURL9)
	})
	assert.Contains(t, out, "Triggered sync for url=http://example.com/rss9")

	// Verify Schedule was created for one-time sync
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundSync *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL9 {
			foundSync = &feeds[i]
		}
	}
	assert.NotNil(t, foundSync)
	assert.False(t, foundSync.Polling, "Polling should remain false")
	assert.NotNil(t, foundSync.FeedSchedule)
	assert.NotNil(t, foundSync.FeedSchedule.NextThinkerAt, "NextThinkerAt should be set for the immediate sync")

	// 22. Test mine on non-mining feed
	testURL10 := "http://example.com/rss10"
	captureOutput(func() {
		addFeed(testURL10, true, false, false, false, "", false, []string{}) // Mining is false
	})

	out = captureOutput(func() {
		mineFeed(testURL10)
	})
	assert.Contains(t, out, "Triggered mining for url=http://example.com/rss10")

	// Verify Schedule was created for one-time mine
	out = captureOutput(func() {
		listFeeds(true, false)
	})
	err = json.Unmarshal([]byte(out), &feeds)
	assert.NoError(t, err)

	var foundMineNoMining *database.Feed
	for i := range feeds {
		if feeds[i].URL == testURL10 {
			foundMineNoMining = &feeds[i]
		}
	}
	assert.NotNil(t, foundMineNoMining)
	assert.False(t, foundMineNoMining.Mining, "Mining should remain false")
	assert.NotNil(t, foundMineNoMining.FeedSchedule)
	assert.NotNil(t, foundMineNoMining.FeedSchedule.NextMiningAt, "NextMiningAt should be set for the immediate mine")
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

func (m *MockRepo) PurgeFeedById(id uuid.UUID) error {
	delete(m.feeds, id)
	return nil
}

func (m *MockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration) error {
	if f, ok := m.feeds[id]; ok {
		now := time.Now()
		f.FeedSchedule = &database.FeedSchedule{
			ID:            id,
			NextThinkerAt: &now,
		}
	}
	return nil
}

func (m *MockRepo) EnqueueMine(id uuid.UUID, miningInterval time.Duration) error {
	if f, ok := m.feeds[id]; ok {
		now := time.Now()
		if f.FeedSchedule == nil {
			f.FeedSchedule = &database.FeedSchedule{ID: id}
		}
		f.FeedSchedule.NextMiningAt = &now
	}
	return nil
}

func (m *MockRepo) RemoveSync(id uuid.UUID) error {
	if f, ok := m.feeds[id]; ok && f.FeedSchedule != nil {
		f.FeedSchedule.NextThinkerAt = nil
		f.FeedSchedule.ThinkerLockedUntil = nil
	}
	return nil
}

func (m *MockRepo) RemoveMine(id uuid.UUID) error {
	if f, ok := m.feeds[id]; ok && f.FeedSchedule != nil {
		f.FeedSchedule.NextMiningAt = nil
		f.FeedSchedule.MiningLockedUntil = nil
	}
	return nil
}

func (m *MockRepo) BeginFeedUpdate(lockDuration time.Duration) (*database.Feed, error) {
	return nil, nil
}

func (m *MockRepo) EndFeedUpdate(id uuid.UUID, err error, successDelay time.Duration) error {
	return nil
}

func (m *MockRepo) GetPendingItems(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error) {
	res := make(map[string]int)
	for _, h := range hashes {
		res[h] = 0
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

func (m *MockRepo) CreateFeedSchedule(feedID uuid.UUID) error {
	if f, ok := m.feeds[feedID]; ok {
		if f.FeedSchedule == nil {
			f.FeedSchedule = &database.FeedSchedule{ID: feedID}
		}
	}
	return nil
}

func (m *MockRepo) FindItemsByRootDomain(rootDomain string, limit int) ([]database.Item, error) {
	return nil, nil
}

func (m *MockRepo) GetTopTrendByDomain(domain string, language string, daysInPast int) ([]database.TrendMetric, error) {
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
