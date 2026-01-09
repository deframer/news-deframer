package database

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"testing"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFindFeedByUrlAndAvailability(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Found_Enabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		// Create isolated test data
		rawURL := "http://test-find-feed-" + uuid.New().String() + ".test"
		feed := Feed{
			URL:     rawURL,
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		u, err := url.Parse(rawURL)
		assert.NoError(t, err)

		// Test with onlyEnabled=true
		found, err := repo.FindFeedByUrlAndAvailability(u, true)
		assert.NoError(t, err)
		if assert.NotNil(t, found) {
			assert.Equal(t, rawURL, found.URL)
		}

		// Test with onlyEnabled=false
		found, err = repo.FindFeedByUrlAndAvailability(u, false)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	})

	t.Run("Found_Disabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		rawURL := "http://test-find-feed-disabled-" + uuid.New().String() + ".test"
		feed := Feed{
			URL:     rawURL,
			Enabled: false,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		u, err := url.Parse(rawURL)
		assert.NoError(t, err)

		// Test with onlyEnabled=true (Should NOT find)
		found, err := repo.FindFeedByUrlAndAvailability(u, true)
		assert.NoError(t, err)
		assert.Nil(t, found)

		// Test with onlyEnabled=false (Should find)
		found, err = repo.FindFeedByUrlAndAvailability(u, false)
		assert.NoError(t, err)
		assert.NotNil(t, found)
	})

	t.Run("NotFound_NonExistent", func(t *testing.T) {
		repo := baseRepo

		u, err := url.Parse("http://does-not-exist.test/feed")
		assert.NoError(t, err)

		feed, err := repo.FindFeedByUrlAndAvailability(u, true)
		assert.NoError(t, err)
		assert.Nil(t, feed)

		feed, err = repo.FindFeedByUrlAndAvailability(u, false)
		assert.NoError(t, err)
		assert.Nil(t, feed)
	})

	t.Run("FindFeedByUrl", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		rawURL := "http://test-find-feed-" + uuid.New().String() + ".test"
		feed := Feed{
			URL:     rawURL,
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		u, err := url.Parse(rawURL)
		assert.NoError(t, err)

		found, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		if assert.NotNil(t, found) {
			assert.Equal(t, rawURL, found.URL)
		}
	})
}

func TestFindFeedById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Found", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		id := uuid.New()
		feed := Feed{
			Base:    Base{ID: id},
			URL:     "http://test-find-feed-id-" + id.String() + ".test",
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		found, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		if assert.NotNil(t, found) {
			assert.Equal(t, feed.ID, found.ID)
			assert.Equal(t, feed.URL, found.URL)
		}
	})

	t.Run("Found_Disabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		id := uuid.New()
		feed := Feed{
			Base:    Base{ID: id},
			URL:     "http://test-find-feed-id-disabled-" + id.String() + ".test",
			Enabled: false,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		found, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, feed.ID, found.ID)
	})

	t.Run("NotFound_NonExistent", func(t *testing.T) {
		repo := baseRepo
		found, err := repo.FindFeedById(uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("Found_Deleted", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		id := uuid.New()
		feed := Feed{
			Base:    Base{ID: id},
			URL:     "http://test-find-feed-id-deleted-" + id.String() + ".test",
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed).Error)
		assert.NoError(t, repo.DeleteFeedById(id))

		found, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, feed.ID, found.ID)
	})
}

func TestUpsertFeed(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Create", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := &Feed{
			URL:     "http://upsert-create.test/" + uuid.New().String(),
			Enabled: true,
		}
		err := repo.UpsertFeed(feed)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, feed.ID)

		found, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, feed.URL, found.URL)
	})

	t.Run("Update", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		// Create first
		feed := &Feed{
			URL:     "http://upsert-update.test/" + uuid.New().String(),
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed))

		// Modify
		feed.Enabled = false
		assert.NoError(t, repo.UpsertFeed(feed))

		found, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.False(t, found.Enabled)
	})

	t.Run("Create_Disabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := &Feed{
			URL:     "http://upsert-disabled.test/" + uuid.New().String(),
			Enabled: false,
		}
		assert.NoError(t, repo.UpsertFeed(feed))

		// FindById should find it (it ignores enabled status)
		foundId, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundId)
		assert.Equal(t, feed.ID, foundId.ID)
		assert.False(t, foundId.Enabled)

		// FindByUrl should NOT find it (it requires enabled=true)
		u, _ := url.Parse(feed.URL)
		foundUrl, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		assert.Nil(t, foundUrl)
	})

	t.Run("Create_Then_Delete", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := &Feed{
			URL:     "http://upsert-delete.test/" + uuid.New().String(),
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed))

		// Delete it
		assert.NoError(t, repo.DeleteFeedById(feed.ID))

		// FindById should find it (Unscoped)
		foundId, err := repo.FindFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundId)
		assert.True(t, foundId.DeletedAt.Valid)

		// FindByUrl should NOT find it (Soft deleted)
		u, _ := url.Parse(feed.URL)
		foundUrl, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		assert.Nil(t, foundUrl)
	})
}

func TestFindCachedFeedById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	tests := []struct {
		name          string
		setup         func(tx *gorm.DB) uuid.UUID
		expectedFound bool
	}{
		{
			name: "Found (Enabled)",
			setup: func(tx *gorm.DB) uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: true,
		},
		{
			name: "NotFound (Disabled)",
			setup: func(tx *gorm.DB) uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (Deleted)",
			setup: func(tx *gorm.DB) uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)

				assert.NoError(t, tx.Delete(&feed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (NonExistent)",
			setup: func(tx *gorm.DB) uuid.UUID {
				return uuid.New()
			},
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := baseDB.Begin()
			defer tx.Rollback()
			repo := NewFromDB(tx)

			id := tt.setup(tx)
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

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	// Helper to create valid SHA256 hash string
	makeHash := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	}

	tests := []struct {
		name          string
		setup         func(tx *gorm.DB) (uuid.UUID, string)
		expectedCount int
	}{
		{
			name: "Found (Enabled)",
			setup: func(tx *gorm.DB) (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				hash := makeHash("item-enabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AnalyzerResult: JSONB{"a": 1}}
				assert.NoError(t, tx.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 1,
		},
		{
			name: "Empty (Disabled)",
			setup: func(tx *gorm.DB) (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				hash := makeHash("item-disabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AnalyzerResult: JSONB{"a": 1}}
				assert.NoError(t, tx.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (Deleted)",
			setup: func(tx *gorm.DB) (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, tx.Create(&feed).Error)

				hash := makeHash("item-deleted-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AnalyzerResult: JSONB{"a": 1}}
				assert.NoError(t, tx.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, tx.Create(&cachedFeed).Error)

				assert.NoError(t, tx.Delete(&feed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (NonExistent)",
			setup: func(tx *gorm.DB) (uuid.UUID, string) {
				return uuid.New(), ""
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := baseDB.Begin()
			defer tx.Rollback()
			repo := NewFromDB(tx)

			id, hash := tt.setup(tx)
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

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	// Helper to create valid SHA256 hash string
	makeHash := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	}

	t.Run("MatchingDomain", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		uid := uuid.New().String()
		// Scenario: A standard news site where items belong to the site.
		feedURL := "http://nytimes.test/services/xml/rss/nyt/HomePage-" + uid + ".xml"
		itemURL := "http://nytimes.test/2024/01/01/world/test-article-" + uid + ".html"

		// EnforceFeedDomain not enforced by Database (!)
		feed := Feed{URL: feedURL, Enabled: true, EnforceFeedDomain: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Test
		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Len(t, items, 0)
	})

	t.Run("Syndication_SameUrlInMultipleFeeds", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		uid := uuid.New().String()
		// Scenario: The same article URL appears in two different feeds.
		// 1. The original source (Enforced)
		// 2. An aggregator (Not Enforced)
		sharedURL := "http://example.com/shared-story-" + uid

		// Create Source Feed & Item
		sourceFeed := Feed{URL: "http://example.com/rss-" + uid, Enabled: true, EnforceFeedDomain: true}
		assert.NoError(t, tx.Create(&sourceFeed).Error)

		// We use the same content hash to simulate exact syndication, though hashes can differ if content varies.
		contentHash := makeHash("shared-content-body")

		sourceItem := Item{Hash: contentHash, FeedID: sourceFeed.ID, URL: sharedURL, AnalyzerResult: JSONB{"src": "original"}}
		assert.NoError(t, tx.Create(&sourceItem).Error)

		// Create Aggregator Feed & Item
		aggFeed := Feed{URL: "http://aggregator.test/rss-" + uid, Enabled: true, EnforceFeedDomain: false}
		assert.NoError(t, tx.Create(&aggFeed).Error)

		aggItem := Item{Hash: contentHash, FeedID: aggFeed.ID, URL: sharedURL, AnalyzerResult: JSONB{"src": "aggregator"}}
		assert.NoError(t, tx.Create(&aggItem).Error)

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
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		uid := uuid.New().String()
		feedURL := "http://disabled-feed.test/rss-" + uid
		itemURL := "http://disabled-feed.test/article-" + uid

		feed := Feed{URL: feedURL, Enabled: false}
		assert.NoError(t, tx.Create(&feed).Error)

		item := Item{
			Hash:           makeHash("disabled-content"),
			FeedID:         feed.ID,
			URL:            itemURL,
			AnalyzerResult: JSONB{"title": "Disabled"},
		}
		assert.NoError(t, tx.Create(&item).Error)

		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})

	t.Run("FeedDeleted", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		uid := uuid.New().String()
		feedURL := "http://deleted-feed.test/rss-" + uid
		itemURL := "http://deleted-feed.test/article-" + uid

		feed := Feed{
			URL:     feedURL,
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		item := Item{
			Hash:           makeHash("deleted-content"),
			FeedID:         feed.ID,
			URL:            itemURL,
			AnalyzerResult: JSONB{"title": "Deleted"},
		}
		assert.NoError(t, tx.Create(&item).Error)

		assert.NoError(t, tx.Delete(&feed).Error)

		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestGetAllFeeds(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("ListFeeds", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		// Create active feed
		feed1 := Feed{
			URL:     "http://active-feed.test/" + uuid.New().String(),
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed1).Error)

		// Create deleted feed
		feed2 := Feed{
			URL:     "http://deleted-feed.test/" + uuid.New().String(),
			Enabled: true,
		}
		assert.NoError(t, tx.Create(&feed2).Error)
		assert.NoError(t, tx.Delete(&feed2).Error)

		// Test GetAllFeeds(false) - should only return active feeds
		feeds, err := repo.GetAllFeeds(false)
		assert.NoError(t, err)

		foundFeed1 := false
		foundFeed2 := false
		for _, f := range feeds {
			if f.ID == feed1.ID {
				foundFeed1 = true
			}
			if f.ID == feed2.ID {
				foundFeed2 = true
			}
		}
		assert.True(t, foundFeed1, "Active feed should be found")
		assert.False(t, foundFeed2, "Deleted feed should not be found when deleted=false")

		// Test GetAllFeeds(true) - should return both
		feedsAll, err := repo.GetAllFeeds(true)
		assert.NoError(t, err)

		foundFeed1 = false
		foundFeed2 = false
		for _, f := range feedsAll {
			if f.ID == feed1.ID {
				foundFeed1 = true
			}
			if f.ID == feed2.ID {
				foundFeed2 = true
			}
		}
		assert.True(t, foundFeed1, "Active feed should be found")
		assert.True(t, foundFeed2, "Deleted feed should be found when deleted=true")
	})
}
