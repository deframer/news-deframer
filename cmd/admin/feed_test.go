package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/downloader"
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

func TestFeedImportHelpers(t *testing.T) {
	u, err := url.Parse("http://example.com/rss")
	assert.NoError(t, err)

	country := "US"
	feed := createNewFeed(u, ImportFeed{URL: u.String(), Country: &country})
	assert.NotNil(t, feed)
	assert.Empty(t, feed.Tags)
	assert.Equal(t, "US", feed.Country)
	assert.Nil(t, feed.Language)

	existing := &database.Feed{Tags: []string{"old"}, Country: "DE"}
	updatedCountry := "FR"
	updateExistingFeed(existing, ImportFeed{Tags: []string{"alpha", "beta"}, Country: &updatedCountry})
	assert.Equal(t, database.StringArray{"alpha", "beta"}, existing.Tags)
	assert.Equal(t, "FR", existing.Country)
}

func TestValidateImportFeeds(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	existing := &database.Feed{URL: "http://existing.example/rss", Enabled: true, Polling: true, Mining: true}
	assert.NoError(t, mock.UpsertFeed(existing))

	country := "US"
	language := "en"
	rootDomain := "new.example"
	feeds := []ImportFeed{
		{
			URL:        "http://new.example/rss/",
			Country:    &country,
			Language:   &language,
			RootDomain: &rootDomain,
			Polling:    boolPtr(true),
			Mining:     boolPtr(false),
			Enabled:    boolPtr(true),
			Categories: []string{"politics", "culture"},
			Tags:       []string{"lead"},
		},
		{
			URL:        "http://existing.example/rss",
			Country:    &country,
			Language:   &language,
			Categories: []string{"politics"},
		},
	}

	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, false)
	assert.False(t, hadErrors)
	assert.Contains(t, out.String(), "1 OK would insert url=http://new.example/rss")
	assert.Contains(t, out.String(), "2 OK would update url=http://existing.example/rss")
	assert.Contains(t, out.String(), "validated 2 feeds (0 errors)")
}

func TestValidateImportFeedsRejectsDuplicates(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	feeds := []ImportFeed{
		{URL: "http://dup.example/rss/", Categories: []string{"politics"}},
		{URL: "http://dup.example/rss", Categories: []string{"politics"}},
	}

	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, false)
	assert.True(t, hadErrors)
	assert.Contains(t, out.String(), "ERROR duplicate url http://dup.example/rss (also in row 1)")
	assert.Contains(t, out.String(), "validated 2 feeds (1 errors)")
}

func TestValidateImportFeedsRejectsInvalidCategory(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	feeds := []ImportFeed{{URL: "http://example.com/rss", Categories: []string{"politics", "not-a-category"}}}

	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, false)
	assert.True(t, hadErrors)
	assert.Contains(t, out.String(), "invalid category \"not-a-category\"")
	assert.Contains(t, out.String(), "validated 1 feeds (1 errors)")
}

func TestValidateImportFeedsRejectsInvalidLanguage(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	language := "eng"
	feeds := []ImportFeed{{URL: "http://example.com/rss", Language: &language, Categories: []string{"politics"}}}

	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, false)
	assert.True(t, hadErrors)
	assert.Contains(t, out.String(), "invalid language \"eng\"")
	assert.Contains(t, out.String(), "validated 1 feeds (1 errors)")
}

func TestValidateImportFeedsAllowsEmptyLanguage(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	empty := ""
	feeds := []ImportFeed{{URL: "http://example.com/rss", Language: &empty, Categories: []string{"politics"}}}

	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, false)
	assert.False(t, hadErrors)
	assert.NotContains(t, out.String(), "ERROR")
	assert.Contains(t, out.String(), "validated 1 feeds (0 errors)")
}

func TestReadImportFeedsRejectsInvalidJSON(t *testing.T) {
	path := t.TempDir() + "/invalid.json"
	assert.NoError(t, os.WriteFile(path, []byte("not-json"), 0o600))

	_, err := readImportFeeds(path)
	assert.Error(t, err)
}

func TestValidateImportFeedsFetchFollowsRedirectAndAcceptsXML(t *testing.T) {
	mock := NewMockRepo()
	repo = mock
	feedDownloader = downloader.NewDownloader(context.Background(), &config.Config{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/feed":
			assert.Contains(t, r.Header.Get("Accept"), "application/rss+xml")
			assert.Contains(t, r.Header.Get("User-Agent"), "Chrome/136.0.0.0")
			assert.Contains(t, r.Header.Get("Sec-CH-UA"), "Chromium")
			assert.Equal(t, "priority: u=0, i", r.Header.Get("Priority"))
			w.Header().Set("Content-Type", "application/rss+xml")
			_, _ = io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?><rss><channel><title>ok</title></channel></rss>`)
		case "/redirect":
			http.Redirect(w, r, "/feed", http.StatusFound)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	feeds := []ImportFeed{{URL: server.URL + "/redirect", Categories: []string{"politics"}}}
	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, true)
	assert.False(t, hadErrors)
	assert.NotContains(t, out.String(), "ERROR")
}

func TestValidateImportFeedsFetchRejectsNonXML(t *testing.T) {
	mock := NewMockRepo()
	repo = mock
	feedDownloader = downloader.NewDownloader(context.Background(), &config.Config{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, "<html><body>no feed</body></html>")
	}))
	defer server.Close()

	feeds := []ImportFeed{{URL: server.URL, Categories: []string{"politics"}}}
	var out bytes.Buffer
	hadErrors := validateImportFeeds(context.Background(), feeds, &out, true)
	assert.True(t, hadErrors)
	assert.Contains(t, out.String(), "ERROR not rss/atom")
	assert.True(t, strings.Contains(out.String(), "validated 1 feeds (1 errors)"))
}

func boolPtr(v bool) *bool {
	return &v
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
		addFeed(testURL, true, false, false, false, "", false, []string{}, "")
	})
	assert.Contains(t, out, "Added feed")
	assert.Contains(t, out, testURL)

	out = captureOutput(func() {
		setTags(testURL, "tag-one, tag-two")
	})
	assert.Contains(t, out, "Set tags to [tag-one tag-two]")

	out = captureOutput(func() {
		setCountry(testURL, "US")
	})
	assert.Contains(t, out, "Set country to US")

	out = captureOutput(func() {
		listFeeds(false, false)
	})
	assert.Contains(t, out, "Tags")
	assert.Contains(t, out, "tag-one,tag-two")
	assert.Contains(t, out, "US")

	out = captureOutput(func() {
		exportFeeds()
	})
	var exported []ImportFeed
	err = json.Unmarshal([]byte(out), &exported)
	assert.NoError(t, err)
	var foundExport *ImportFeed
	for i := range exported {
		if exported[i].URL == testURL {
			foundExport = &exported[i]
		}
	}
	assert.NotNil(t, foundExport)
	assert.Equal(t, []string{"tag-one", "tag-two"}, foundExport.Tags)
	assert.Equal(t, "US", *foundExport.Country)

	out = captureOutput(func() {
		setCountry(testURL, "")
	})
	assert.Contains(t, out, "Set country to ")

	out = captureOutput(func() {
		setTags(testURL, "")
	})
	assert.Contains(t, out, "Set tags to []")

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
	assert.Empty(t, feeds[0].Country)
	assert.Empty(t, feeds[0].Tags)
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
		addFeed(testURL2, false, true, false, false, "", false, []string{}, "") // Add disabled feed with polling=true
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
		addFeed(testURL3, true, false, false, true, "", false, []string{}, "")
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
		addFeed(testURL4, true, false, false, false, "", false, []string{}, "")
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
		addFeed(testURL5, true, false, false, false, "en", false, []string{}, "")
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
		addFeed(testURL6, true, false, false, false, "", true, []string{}, "")
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
		addFeed(testURL7, true, false, true, false, "", false, []string{}, "")
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
		addFeed(testURL8, true, false, false, false, "", false, []string{}, "")
	})

	out = captureOutput(func() {
		mineAllFeeds()
	})
	assert.Contains(t, out, "Triggered mining for url=http://example.com/rss7")
	assert.NotContains(t, out, "Triggered mining for url=http://example.com/rss8")

	// 21. Test sync on non-polling feed
	testURL9 := "http://example.com/rss9"
	captureOutput(func() {
		addFeed(testURL9, true, false, false, false, "", false, []string{}, "") // Polling is false
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
		addFeed(testURL10, true, false, false, false, "", false, []string{}, "") // Mining is false
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

func TestFeedErrorsCommand(t *testing.T) {
	mock := NewMockRepo()
	repo = mock

	rootDomain := "example.com"
	lastError := "boom"
	lastSyncedAt := time.Date(2026, time.May, 17, 12, 34, 0, 0, time.UTC)
	assert.NoError(t, mock.UpsertFeed(&database.Feed{
		URL:          "http://example.com/rss",
		RootDomain:   &rootDomain,
		LastError:    &lastError,
		LastSyncedAt: &lastSyncedAt,
	}))

	out := captureOutput(func() {
		listFeedErrors()
	})

	assert.Contains(t, out, "RootDomain")
	assert.Contains(t, out, "URL")
	assert.Contains(t, out, "LastSyncedAt")
	assert.Contains(t, out, "Error")
	assert.Contains(t, out, "example.com")
	assert.Contains(t, out, "http://example.com/rss")
	assert.Contains(t, out, "2026-05-17 12:34")
	assert.Contains(t, out, "boom")
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

func (m *MockRepo) UpsertItemWithTrendInvalidation(item *database.Item) error {
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

func (m *MockRepo) GetAllFeedErrors() ([]database.FeedError, error) {
	var feedErrors []database.FeedError
	for _, f := range m.feeds {
		if f.LastError == nil {
			continue
		}
		feedErrors = append(feedErrors, database.FeedError{
			RootDomain:   f.RootDomain,
			URL:          f.URL,
			Error:        *f.LastError,
			LastSyncedAt: f.LastSyncedAt,
		})
	}
	return feedErrors, nil
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

func (m *MockRepo) BeginThinkerBatch(limit int, since time.Time, minErrorCount int, maxErrorCount int, lockDuration time.Duration) ([]database.Item, error) {
	return nil, nil
}

func (m *MockRepo) BeginThinkerFixerBatch(limit int, since time.Time, minErrorCount int, maxErrorCount int, lockDuration time.Duration) ([]database.Item, error) {
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

func (m *MockRepo) FindAnalyzedItemsByRootDomain(rootDomain string, limit int) ([]database.AnalyzedItem, error) {
	return nil, nil
}

func (m *MockRepo) FindFirstAnalyzedItemByUrl(u *url.URL) (*database.AnalyzedItem, error) {
	return nil, nil
}

func (m *MockRepo) GetTopTrendByDomain(domain string, language string, date *time.Time, days int) ([]database.TrendMetric, error) {
	return nil, nil
}

func (m *MockRepo) GetContextByDomain(term string, domain string, language string, date *time.Time, days int) ([]database.TrendContext, error) {
	return nil, nil
}

func (m *MockRepo) GetLifecycleByDomain(term string, domain string, language string, date *time.Time, days int) ([]database.Lifecycle, error) {
	return nil, nil
}

func (m *MockRepo) GetDomainComparison(domainA string, domainB string, language string, date *time.Time, days int, utilityThreshold float64, outlierRatioThreshold float64, limit int) ([]database.DomainComparison, error) {
	return nil, nil
}

func (m *MockRepo) GetArticlesByTrend(term string, domain string, date *time.Time, days int, limit int, offset int) ([]database.AnalyzedArticle, error) {
	return nil, nil
}

func (m *MockRepo) GetSentimentsByTrend(term string, domain string, date *time.Time, days int) (*database.SentimentItem, error) {
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
