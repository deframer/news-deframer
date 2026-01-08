package database

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTime(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)

	got, err := repo.GetTime()
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), got, 1*time.Minute)
}

func TestFindFeedByUrl(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	t.Run("Found", func(t *testing.T) {
		// Create isolated test data
		rawURL := "http://test-find-feed-" + uuid.New().String() + ".test"
		feed := Feed{
			URL:     rawURL,
			Enabled: true,
		}
		assert.NoError(t, db.Create(&feed).Error)

		u, err := url.Parse(rawURL)
		assert.NoError(t, err)

		found, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		if assert.NotNil(t, found) {
			assert.Equal(t, rawURL, found.URL)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		u, err := url.Parse("http://does-not-exist.test/feed")
		assert.NoError(t, err)

		feed, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		assert.Nil(t, feed)
	})
}

func TestFindCachedFeedById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		expectedFound bool
	}{
		{
			name: "Found (Enabled)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: true,
		},
		{
			name: "NotFound (Disabled)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (Deleted)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)

				assert.NoError(t, db.Delete(&feed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (NonExistent)",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.setup()
			found, err := repo.FindCachedFeedById(id)
			assert.NoError(t, err)
			if tt.expectedFound {
				assert.NotNil(t, found)
				assert.Equal(t, id, found.ID)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}

func TestFindItemsByCachedFeedFeedId(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	// Helper to create valid SHA256 hash string
	makeHash := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	}

	tests := []struct {
		name          string
		setup         func() (uuid.UUID, string)
		expectedCount int
	}{
		{
			name: "Found (Enabled)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-enabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 1,
		},
		{
			name: "Empty (Disabled)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-disabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (Deleted)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-deleted-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)

				assert.NoError(t, db.Delete(&feed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (NonExistent)",
			setup: func() (uuid.UUID, string) {
				return uuid.New(), ""
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, hash := tt.setup()
			items, err := repo.FindItemsByCachedFeedFeedId(id)
			assert.NoError(t, err)
			assert.Len(t, items, tt.expectedCount)
			if tt.expectedCount > 0 {
				assert.Equal(t, hash, items[0].Hash)
			}
		})
	}
}

func TestFindItemsByUrl(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	// Helper to create valid SHA256 hash string
	makeHash := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	}

	t.Run("MatchingDomain", func(t *testing.T) {
		// Scenario: A standard news site where items belong to the site.
		feedURL := "http://nytimes.test/services/xml/rss/nyt/HomePage.xml"
		itemURL := "http://nytimes.test/2024/01/01/world/test-article.html"

		// EnforceFeedDomain not enforced by Database (!)
		feed := Feed{URL: feedURL, Enabled: true, EnforceFeedDomain: true}
		assert.NoError(t, db.Create(&feed).Error)

		// Test
		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Len(t, items, 0)
	})

	t.Run("Syndication_SameUrlInMultipleFeeds", func(t *testing.T) {
		// Scenario: The same article URL appears in two different feeds.
		// 1. The original source (Enforced)
		// 2. An aggregator (Not Enforced)
		sharedURL := "http://example.com/shared-story"

		// Create Source Feed & Item
		sourceFeed := Feed{URL: "http://example.com/rss", Enabled: true, EnforceFeedDomain: true}
		assert.NoError(t, db.Create(&sourceFeed).Error)

		// We use the same content hash to simulate exact syndication, though hashes can differ if content varies.
		contentHash := makeHash("shared-content-body")

		sourceItem := Item{Hash: contentHash, FeedID: sourceFeed.ID, URL: sharedURL, AIResult: JSONB{"src": "original"}}
		assert.NoError(t, db.Create(&sourceItem).Error)

		// Create Aggregator Feed & Item
		aggFeed := Feed{URL: "http://aggregator.test/rss", Enabled: true, EnforceFeedDomain: false}
		assert.NoError(t, db.Create(&aggFeed).Error)

		aggItem := Item{Hash: contentHash, FeedID: aggFeed.ID, URL: sharedURL, AIResult: JSONB{"src": "aggregator"}}
		assert.NoError(t, db.Create(&aggItem).Error)

		// Test: Should find BOTH items
		u, _ := url.Parse(sharedURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		// Verify we got one from each feed
		feedIDs := map[uuid.UUID]bool{items[0].FeedID: true, items[1].FeedID: true}
		assert.True(t, feedIDs[sourceFeed.ID])
		assert.True(t, feedIDs[aggFeed.ID])
	})

	t.Run("FeedDisabled", func(t *testing.T) {
		feedURL := "http://disabled-feed.test/rss"
		itemURL := "http://disabled-feed.test/article"

		feed := Feed{URL: feedURL, Enabled: false}
		assert.NoError(t, db.Create(&feed).Error)

		item := Item{
			Hash:     makeHash("disabled-content"),
			FeedID:   feed.ID,
			URL:      itemURL,
			AIResult: JSONB{"title": "Disabled"},
		}
		assert.NoError(t, db.Create(&item).Error)

		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("FeedDeleted", func(t *testing.T) {
		feedURL := "http://deleted-feed.test/rss"
		itemURL := "http://deleted-feed.test/article"

		feed := Feed{
			URL:     feedURL,
			Enabled: true,
		}
		assert.NoError(t, db.Create(&feed).Error)

		item := Item{
			Hash:     makeHash("deleted-content"),
			FeedID:   feed.ID,
			URL:      itemURL,
			AIResult: JSONB{"title": "Deleted"},
		}
		assert.NoError(t, db.Create(&item).Error)

		assert.NoError(t, db.Delete(&feed).Error)

		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}
