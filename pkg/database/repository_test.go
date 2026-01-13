package database

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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
			Polling: false,
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

	t.Run("Constraint_UniqueUrl", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		urlStr := "http://unique-constraint.test/" + uuid.New().String()

		// 1. Create first feed
		feed1 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed1))

		// 2. Create second feed with same URL (Should Fail)
		feed2 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		err := repo.UpsertFeed(feed2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// 3. Soft Delete feed1
		assert.NoError(t, repo.DeleteFeedById(feed1.ID))

		// 4. Create second feed again (Should Succeed now)
		err = repo.UpsertFeed(feed2)
		assert.NoError(t, err)

		// 5. Update feed2 to a new URL
		newUrl := "http://unique-constraint-new.test/" + uuid.New().String()
		feed2.URL = newUrl
		assert.NoError(t, repo.UpsertFeed(feed2))

		// 6. Create feed3 with the OLD url of feed2 (Should Succeed)
		feed3 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed3))

		// 7. Try to update feed3 to feed2's current URL (Should Fail)
		feed3.URL = newUrl
		err = repo.UpsertFeed(feed3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestUpsertItem(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("CreateNew", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://upsert-item-new.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		item := &Item{
			FeedID:      feed.ID,
			Hash:        "hash-new",
			URL:         "http://item-new",
			Content:     "content",
			ThinkResult: &ThinkResult{TitleCorrected: "bar"},
			ThinkRating: 0.5,
		}

		err := repo.UpsertItem(item)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, item.ID)

		var count int64
		tx.Model(&Item{}).Where("id = ?", item.ID).Count(&count)
		assert.Equal(t, int64(1), count)

		var stored Item
		tx.First(&stored, item.ID)
		assert.Equal(t, 0.5, stored.ThinkRating)
	})

	t.Run("UpdateExisting_ByHash", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://upsert-item-update.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Create initial item
		existing := Item{
			FeedID:      feed.ID,
			Hash:        "hash-update",
			URL:         "http://item-update",
			Content:     "content-old",
			ThinkResult: &ThinkResult{Framing: 0.1},
			ThinkRating: 0.1,
		}
		assert.NoError(t, tx.Create(&existing).Error)

		// Sleep to ensure UpdatedAt changes
		time.Sleep(10 * time.Millisecond)

		// Update with same hash/feedID but different content, and NO ID provided
		update := &Item{
			FeedID:      feed.ID,
			Hash:        "hash-update",
			URL:         "http://item-update",
			Content:     "content-new",
			ThinkResult: &ThinkResult{Framing: 0.2},
			ThinkRating: 0.2,
		}

		err := repo.UpsertItem(update)
		assert.NoError(t, err)
		assert.Equal(t, existing.ID, update.ID) // Should have picked up the ID

		var stored Item
		tx.First(&stored, existing.ID)
		assert.Equal(t, "content-new", stored.Content)
		// Verify JSONB update
		assert.Equal(t, 0.2, stored.ThinkResult.Framing)
		assert.Equal(t, 0.2, stored.ThinkRating)

		// Verify timestamps
		assert.Equal(t, existing.CreatedAt.UTC(), stored.CreatedAt.UTC(), "CreatedAt should be preserved")
		assert.True(t, stored.UpdatedAt.After(existing.UpdatedAt), "UpdatedAt should be updated")
	})
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
		feed := Feed{URL: feedURL, Enabled: true, EnforceFeedDomain: true, Polling: false}
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
		sourceFeed := Feed{URL: "http://example.com/rss-" + uid, Enabled: true, EnforceFeedDomain: true, Polling: false}
		assert.NoError(t, tx.Create(&sourceFeed).Error)

		// We use the same content hash to simulate exact syndication, though hashes can differ if content varies.
		contentHash := makeHash("shared-content-body")

		sourceItem := Item{Hash: contentHash, FeedID: sourceFeed.ID, URL: sharedURL, ThinkResult: &ThinkResult{TitleCorrected: "original"}, PubDate: time.Now()}
		assert.NoError(t, tx.Create(&sourceItem).Error)

		// Create Aggregator Feed & Item
		aggFeed := Feed{URL: "http://aggregator.test/rss-" + uid, Enabled: true, EnforceFeedDomain: false, Polling: false}
		assert.NoError(t, tx.Create(&aggFeed).Error)

		aggItem := Item{Hash: contentHash, FeedID: aggFeed.ID, URL: sharedURL, ThinkResult: &ThinkResult{TitleCorrected: "aggregator"}, PubDate: time.Now()}
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

		feed := Feed{URL: feedURL, Enabled: false, Polling: false}
		assert.NoError(t, tx.Create(&feed).Error)

		item := Item{
			Hash:        makeHash("disabled-content"),
			FeedID:      feed.ID,
			URL:         itemURL,
			ThinkResult: &ThinkResult{TitleCorrected: "Disabled"},
			PubDate:     time.Now(),
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
			Polling: false,
		}
		assert.NoError(t, tx.Create(&feed).Error)

		item := Item{
			Hash:        makeHash("deleted-content"),
			FeedID:      feed.ID,
			URL:         itemURL,
			ThinkResult: &ThinkResult{TitleCorrected: "Deleted"},
			PubDate:     time.Now(),
		}
		assert.NoError(t, tx.Create(&item).Error)

		assert.NoError(t, tx.Delete(&feed).Error)

		u, _ := url.Parse(itemURL)
		items, err := repo.FindItemsByUrl(u)
		assert.NoError(t, err)
		assert.Empty(t, items)
	})
}

func TestFindItemsByRootDomain(t *testing.T) {
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

	t.Run("FilterByRootDomain", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		root := "example.com"
		feed := Feed{URL: "http://example.com/rss", Enabled: true, RootDomain: &root}
		assert.NoError(t, tx.Create(&feed).Error)

		h1 := makeHash("hash1")
		// Item 1
		item1 := Item{
			FeedID:      feed.ID,
			Hash:        h1,
			URL:         "http://example.com/1",
			Content:     "c1",
			ThinkResult: &ThinkResult{TitleCorrected: "bar"},
			PubDate:     time.Now(),
			ThinkRating: 0.8,
		}
		assert.NoError(t, tx.Create(&item1).Error)

		// Item 2 (Different root domain feed)
		otherRoot := "other.com"
		feed2 := Feed{URL: "http://other.com/rss", Enabled: true, RootDomain: &otherRoot}
		assert.NoError(t, tx.Create(&feed2).Error)
		h2 := makeHash("hash2")
		item2 := Item{FeedID: feed2.ID, Hash: h2, URL: "http://other.com/1", Content: "c2", PubDate: time.Now()}
		assert.NoError(t, tx.Create(&item2).Error)

		items, err := repo.FindItemsByRootDomain("example.com", 10)
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, h1, items[0].Hash)
		assert.Equal(t, 0.8, items[0].ThinkRating)
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
			Polling: false,
		}
		assert.NoError(t, tx.Create(&feed1).Error)

		// Create deleted feed
		feed2 := Feed{
			URL:     "http://deleted-feed.test/" + uuid.New().String(),
			Enabled: true,
			Polling: false,
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

func TestEnqueueSync(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("CheckFeedAvailability", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		// Case 1: Feed does not exist
		err := repo.EnqueueSync(uuid.New(), 0, time.Minute)
		assert.Error(t, err)
		assert.Equal(t, "feed not found or unavailable", err.Error())

		// Case 2: Feed is disabled
		feedDisabled := Feed{URL: "http://disabled.test", Enabled: false}
		assert.NoError(t, tx.Create(&feedDisabled).Error)
		err = repo.EnqueueSync(feedDisabled.ID, 0, time.Minute)
		assert.Error(t, err)
		assert.Equal(t, "feed not found or unavailable", err.Error())

		// Case 3: Feed is deleted
		feedDeleted := Feed{URL: "http://deleted.test", Enabled: true}
		assert.NoError(t, tx.Create(&feedDeleted).Error)
		assert.NoError(t, tx.Delete(&feedDeleted).Error)
		err = repo.EnqueueSync(feedDeleted.ID, 0, time.Minute)
		assert.Error(t, err)
		assert.Equal(t, "feed not found or unavailable", err.Error())
	})

	t.Run("RespectLongLock", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://long-lock.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Lock for 1 hour
		future := time.Now().Add(time.Hour)
		schedule := FeedSchedule{ID: feed.ID, LockedUntil: &future}
		assert.NoError(t, tx.Create(&schedule).Error)

		// Try to sync with 5 min lock
		err := repo.EnqueueSync(feed.ID, 0, 5*time.Minute)
		assert.NoError(t, err)

		// Verify lock was NOT updated (should still be ~1 hour from now)
		var updated FeedSchedule
		assert.NoError(t, tx.First(&updated, feed.ID).Error)
		assert.WithinDuration(t, future, *updated.LockedUntil, 5*time.Second)
	})

	t.Run("CreateSchedule", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://create-schedule.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		err := repo.EnqueueSync(feed.ID, 0, time.Minute)
		assert.NoError(t, err)

		var schedule FeedSchedule
		err = tx.Where("id = ?", feed.ID).First(&schedule).Error
		assert.NoError(t, err)
		assert.NotNil(t, schedule.NextRunAt)
		assert.WithinDuration(t, time.Now(), *schedule.NextRunAt, 5*time.Second)
		assert.Nil(t, schedule.LockedUntil)
		assert.Nil(t, schedule.LastError)
	})

	t.Run("UpdateSchedule_ForceUnlock", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://update-schedule.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Create existing schedule with a lock (short duration, so it passes the check)
		future := time.Now().Add(time.Minute)
		errMsg := "old error"
		schedule := FeedSchedule{
			ID:          feed.ID,
			LockedUntil: &future,
			LastError:   &errMsg,
		}
		assert.NoError(t, tx.Create(&schedule).Error)

		// EnqueueSync with 5 min lock duration (1 min < 5 min, so it proceeds)
		err := repo.EnqueueSync(feed.ID, 0, 5*time.Minute)
		assert.NoError(t, err)

		var updated FeedSchedule
		err = tx.Where("id = ?", feed.ID).First(&updated).Error
		assert.NoError(t, err)
		assert.NotNil(t, updated.NextRunAt)
		assert.WithinDuration(t, time.Now(), *updated.NextRunAt, 5*time.Second)
		assert.Nil(t, updated.LockedUntil) // Should be nilled out
		assert.Nil(t, updated.LastError)   // Should be nilled out
	})
}

func TestRemoveSync(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("RemoveExisting", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://remove-sync.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Create schedule
		now := time.Now()
		errMsg := "some error"
		schedule := FeedSchedule{
			ID:        feed.ID,
			NextRunAt: &now,
			LastError: &errMsg,
		}
		assert.NoError(t, tx.Create(&schedule).Error)

		err := repo.RemoveSync(feed.ID)
		assert.NoError(t, err)

		var updated FeedSchedule
		err = tx.Where("id = ?", feed.ID).First(&updated).Error
		assert.NoError(t, err)
		assert.Nil(t, updated.NextRunAt)
		assert.Nil(t, updated.LockedUntil)
		assert.Nil(t, updated.LastError)
	})

	t.Run("RemoveNonExistent", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		err := repo.RemoveSync(uuid.New())
		assert.NoError(t, err)
	})
}

func TestBeginFeedUpdate(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("NoDueFeeds", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		tx.Exec("DELETE FROM feed_schedules") // Clear seeded schedules
		repo := NewFromDB(tx)

		feed, err := repo.BeginFeedUpdate(15 * time.Minute)
		assert.NoError(t, err)
		assert.Nil(t, feed)
	})

	t.Run("DueFeed_Valid", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		tx.Exec("DELETE FROM feed_schedules")
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://due-valid.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Create schedule due now
		now := time.Now().Add(-1 * time.Minute)
		schedule := FeedSchedule{ID: feed.ID, NextRunAt: &now}
		assert.NoError(t, tx.Create(&schedule).Error)

		f, err := repo.BeginFeedUpdate(15 * time.Minute)
		assert.NoError(t, err)
		assert.NotNil(t, f)
		assert.Equal(t, feed.ID, f.ID)

		// Verify locked
		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		if assert.NotNil(t, s.LockedUntil) {
			assert.True(t, s.LockedUntil.After(time.Now()))
		}
	})

	t.Run("DueFeed_Disabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		tx.Exec("DELETE FROM feed_schedules")
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://due-disabled.test", Enabled: false}
		assert.NoError(t, tx.Create(&feed).Error)

		now := time.Now().Add(-1 * time.Minute)
		schedule := FeedSchedule{ID: feed.ID, NextRunAt: &now}
		assert.NoError(t, tx.Create(&schedule).Error)

		f, err := repo.BeginFeedUpdate(15 * time.Minute)
		assert.NoError(t, err)
		assert.Nil(t, f)

		// Verify schedule cleared (RemoveSync logic)
		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		assert.Nil(t, s.NextRunAt)
	})

	t.Run("DueFeed_Deleted", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		tx.Exec("DELETE FROM feed_schedules")
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://due-deleted.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)
		assert.NoError(t, tx.Delete(&feed).Error)

		now := time.Now().Add(-1 * time.Minute)
		schedule := FeedSchedule{ID: feed.ID, NextRunAt: &now}
		assert.NoError(t, tx.Create(&schedule).Error)

		f, err := repo.BeginFeedUpdate(15 * time.Minute)
		assert.NoError(t, err)
		assert.Nil(t, f)

		// Verify schedule cleared
		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		assert.Nil(t, s.NextRunAt)
	})
}

func TestEndFeedUpdate(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Polling_Enabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://end-polling.test", Enabled: true, Polling: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// Create locked schedule
		future := time.Now().Add(time.Hour)
		schedule := FeedSchedule{ID: feed.ID, LockedUntil: &future}
		assert.NoError(t, tx.Create(&schedule).Error)

		pollingInterval := 10 * time.Minute
		err := repo.EndFeedUpdate(feed.ID, nil, pollingInterval)
		assert.NoError(t, err)

		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		assert.Nil(t, s.LockedUntil)
		assert.Nil(t, s.LastError)
		assert.NotNil(t, s.NextRunAt)
		// Should be scheduled for future (now + pollingInterval)
		assert.WithinDuration(t, time.Now().Add(pollingInterval), *s.NextRunAt, 5*time.Second)
	})

	t.Run("With_Error", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://end-error.test", Enabled: true, Polling: true}
		assert.NoError(t, tx.Create(&feed).Error)

		schedule := FeedSchedule{ID: feed.ID}
		assert.NoError(t, tx.Create(&schedule).Error)

		jobErr := errors.New("something went wrong")
		err := repo.EndFeedUpdate(feed.ID, jobErr, time.Minute)
		assert.NoError(t, err)

		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		if assert.NotNil(t, s.LastError) {
			assert.Equal(t, "something went wrong", *s.LastError)
		}
	})

	t.Run("Polling_Disabled", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://end-no-polling.test", Enabled: true, Polling: false}
		assert.NoError(t, tx.Create(&feed).Error)

		// Schedule was due now
		now := time.Now()
		schedule := FeedSchedule{ID: feed.ID, NextRunAt: &now}
		assert.NoError(t, tx.Create(&schedule).Error)

		err := repo.EndFeedUpdate(feed.ID, nil, time.Minute)
		assert.NoError(t, err)

		var s FeedSchedule
		assert.NoError(t, tx.First(&s, feed.ID).Error)
		assert.Nil(t, s.LockedUntil)
		assert.Nil(t, s.NextRunAt)
	})
}

func TestGetPendingHashes(t *testing.T) {
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

	t.Run("FilterProcessed", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://pending.test", Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// 1. Processed Item
		h1 := makeHash("item1")
		item1 := Item{
			FeedID:      feed.ID,
			Hash:        h1,
			URL:         "http://item1",
			Content:     "c1",
			ThinkResult: &ThinkResult{TitleCorrected: "bar"},
		}
		assert.NoError(t, tx.Create(&item1).Error)

		// 2. Unprocessed Item (AnalyzerResult is nil)
		h2 := makeHash("item2")
		item2 := Item{
			FeedID:      feed.ID,
			Hash:        h2,
			URL:         "http://item2",
			Content:     "c2",
			ThinkResult: nil,
		}
		assert.NoError(t, tx.Create(&item2).Error)

		// 3. New Item (Not in DB)
		h3 := makeHash("item3")

		input := []string{h1, h2, h3}
		pendingMap, err := repo.GetPendingHashes(feed.ID, input)
		assert.NoError(t, err)

		assert.Len(t, pendingMap, 2)
		assert.True(t, pendingMap[h2])
		assert.True(t, pendingMap[h3])
		assert.False(t, pendingMap[h1])
	})
}

func TestGetItemsByHashes(t *testing.T) {
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

	t.Run("FetchItems", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://items-by-hash.test/" + uuid.New().String(), Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		h1 := makeHash("item1")
		item1 := Item{
			FeedID:      feed.ID,
			Hash:        h1,
			URL:         "http://item1/" + uuid.New().String(),
			Content:     "c1",
			ThinkResult: &ThinkResult{TitleCorrected: "bar"},
		}
		assert.NoError(t, tx.Create(&item1).Error)

		h2 := makeHash("item2")
		item2 := Item{
			FeedID:      feed.ID,
			Hash:        h2,
			URL:         "http://item2/" + uuid.New().String(),
			Content:     "c2",
			ThinkResult: nil,
		}
		assert.NoError(t, tx.Create(&item2).Error)

		// Item in another feed
		feed2 := Feed{URL: "http://other-feed.test/" + uuid.New().String(), Enabled: true}
		assert.NoError(t, tx.Create(&feed2).Error)
		h3 := makeHash("item3")
		item3 := Item{
			FeedID:  feed2.ID,
			Hash:    h3,
			URL:     "http://item3/" + uuid.New().String(),
			Content: "c3",
		}
		assert.NoError(t, tx.Create(&item3).Error)

		// Test 1: Get specific items from feed 1
		// h3 is in another feed, so it should not be returned even if requested
		items, err := repo.GetItemsByHashes(feed.ID, []string{h1, h2, h3, "non-existent"})
		assert.NoError(t, err)
		assert.Len(t, items, 2)

		// Verify contents
		foundH1 := false
		foundH2 := false
		for _, item := range items {
			if item.Hash == h1 {
				foundH1 = true
				assert.Equal(t, "c1", item.Content)
			}
			if item.Hash == h2 {
				foundH2 = true
				assert.Equal(t, "c2", item.Content)
			}
		}
		assert.True(t, foundH1, "Should find item1")
		assert.True(t, foundH2, "Should find item2")

		// Test 2: Empty list
		itemsEmpty, err := repo.GetItemsByHashes(feed.ID, []string{})
		assert.NoError(t, err)
		assert.Empty(t, itemsEmpty)
	})
}

func TestUpsertCachedFeed(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("CreateAndUpdate", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		// 1. Create Feed
		feed := Feed{URL: "http://cached-feed.test/" + uuid.New().String(), Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		// 2. Create CachedFeed
		xmlContent := "<rss>...</rss>"
		cachedFeed := &CachedFeed{
			ID:         feed.ID,
			XMLContent: &xmlContent,
			ItemRefs:   StringArray{"hash1", "hash2"},
		}

		// 3. Upsert (Create)
		err := repo.UpsertCachedFeed(cachedFeed)
		assert.NoError(t, err)

		// Verify
		var stored CachedFeed
		err = tx.First(&stored, feed.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, feed.ID, stored.ID)
		assert.Equal(t, xmlContent, *stored.XMLContent)
		assert.Equal(t, StringArray{"hash1", "hash2"}, stored.ItemRefs)

		// 4. Upsert (Update)
		newXML := "<rss>updated</rss>"
		cachedFeed.XMLContent = &newXML
		cachedFeed.ItemRefs = StringArray{"hash3"}

		err = repo.UpsertCachedFeed(cachedFeed)
		assert.NoError(t, err)

		// Verify Update
		var updated CachedFeed
		err = tx.First(&updated, feed.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, newXML, *updated.XMLContent)
		assert.Equal(t, StringArray{"hash3"}, updated.ItemRefs)
	})
}

func TestFindCachedFeedById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Found", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://cached-find.test/" + uuid.New().String(), Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		xmlContent := "<rss>found</rss>"
		cached := CachedFeed{
			ID:         feed.ID,
			XMLContent: &xmlContent,
			ItemRefs:   StringArray{"h1"},
		}
		assert.NoError(t, tx.Create(&cached).Error)

		found, err := repo.FindCachedFeedById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, feed.ID, found.ID)
		assert.Equal(t, xmlContent, *found.XMLContent)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := baseRepo
		found, err := repo.FindCachedFeedById(uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("Constraint_UniqueUrl", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		urlStr := "http://unique-constraint.test/" + uuid.New().String()

		// 1. Create first feed
		feed1 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed1))

		// 2. Create second feed with same URL (Should Fail)
		feed2 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		err := repo.UpsertFeed(feed2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// 3. Soft Delete feed1
		assert.NoError(t, repo.DeleteFeedById(feed1.ID))

		// 4. Create second feed again (Should Succeed now)
		err = repo.UpsertFeed(feed2)
		assert.NoError(t, err)

		// 5. Update feed2 to a new URL
		newUrl := "http://unique-constraint-new.test/" + uuid.New().String()
		feed2.URL = newUrl
		assert.NoError(t, repo.UpsertFeed(feed2))

		// 6. Create feed3 with the OLD url of feed2 (Should Succeed)
		feed3 := &Feed{
			URL:     urlStr,
			Enabled: true,
		}
		assert.NoError(t, repo.UpsertFeed(feed3))

		// 7. Try to update feed3 to feed2's current URL (Should Fail)
		feed3.URL = newUrl
		err = repo.UpsertFeed(feed3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestFindFeedScheduleById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	baseRepo, err := NewRepository(cfg)
	assert.NoError(t, err)
	baseDB := baseRepo.(*repository).db

	t.Run("Found", func(t *testing.T) {
		tx := baseDB.Begin()
		defer tx.Rollback()
		repo := NewFromDB(tx)

		feed := Feed{URL: "http://schedule-find.test/" + uuid.New().String(), Enabled: true}
		assert.NoError(t, tx.Create(&feed).Error)

		now := time.Now()
		schedule := FeedSchedule{
			ID:        feed.ID,
			NextRunAt: &now,
		}
		assert.NoError(t, tx.Create(&schedule).Error)

		found, err := repo.FindFeedScheduleById(feed.ID)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, feed.ID, found.ID)
		assert.WithinDuration(t, now, *found.NextRunAt, time.Second)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := baseRepo
		found, err := repo.FindFeedScheduleById(uuid.New())
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}
