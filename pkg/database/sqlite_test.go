package database

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *Database {
	// Use in-memory SQLite for testing
	db, err := NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func TestNewDatabase(t *testing.T) {
	db := setupTestDB(t)
	assert.NotNil(t, db, "Database should be initialized")
}

func TestTableExists(t *testing.T) {
	db := setupTestDB(t)

	// Check if the 'items' table exists
	var tableExists bool
	err := db.db.Raw("SELECT count(*) > 0 FROM sqlite_master WHERE type='table' AND name='items'").Scan(&tableExists).Error
	assert.NoError(t, err, "Querying table existence should succeed")
	assert.True(t, tableExists, "Table 'items' should exist after migration")

	// Check if the 'caches' table exists
	err = db.db.Raw("SELECT count(*) > 0 FROM sqlite_master WHERE type='table' AND name='caches'").Scan(&tableExists).Error
	assert.NoError(t, err, "Querying table existence should succeed")
	assert.True(t, tableExists, "Table 'caches' should exist after migration")
}

func TestCreateItem(t *testing.T) {
	db := setupTestDB(t)

	// Create test item
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item := &Item{
		Hash:        hash,
		FeedUrl:     "dummy1",
		Link:        "dummy1",
		Guid:        "dummy1",
		Title:       "dummy1",
		Description: "dummy1",
		Framing:     new(float32),
		TitleAI:     new(string),
		ReasonAI:    new(string),
	}
	*item.Framing = 0.5
	*item.TitleAI = "Test Title"
	*item.ReasonAI = "Test Reason"

	// Insert item
	err := db.CreateItem(item)
	assert.NoError(t, err, "Item creation should succeed")

	// Verify item exists
	found, err := db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.NotNil(t, found, "Item should be found")
	assert.Equal(t, hash, found.Hash, "Hash should match")
	assert.Equal(t, *item.Framing, *found.Framing, "Framing should match")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "TitleAI should match")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "ReasonAI should match")

	// Test creating duplicate hash (should silently ignore)
	duplicate := &Item{
		Hash:        hash,
		FeedUrl:     "dummy2",
		Link:        "dummy2",
		Guid:        "dummy2",
		Title:       "dummy2",
		Description: "dummy2",
		Framing:     new(float32),
		TitleAI:     new(string),
		ReasonAI:    new(string),
	}
	*duplicate.Framing = 0.7
	*duplicate.TitleAI = "Different Title"
	*duplicate.ReasonAI = "Different Reason"
	err = db.CreateItem(duplicate)
	assert.NoError(t, err, "Creating item with duplicate hash should succeed (ignored)")

	// Verify original item is unchanged
	found, err = db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.Equal(t, *item.Framing, *found.Framing, "Original Framing should remain")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "Original TitleAI should remain")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "Original ReasonAI should remain")
}

func TestFindItemByHash(t *testing.T) {
	db := setupTestDB(t)

	// Create test item
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item := &Item{
		Hash:        hash,
		FeedUrl:     "dummy3",
		Link:        "dummy3",
		Guid:        "dummy3",
		Title:       "dummy3",
		Description: "dummy3",
		Framing:     new(float32),
		TitleAI:     new(string),
		ReasonAI:    new(string),
	}
	*item.Framing = 0.5
	*item.TitleAI = "Test Title"
	*item.ReasonAI = "Test Reason"

	// Insert item
	err := db.CreateItem(item)
	assert.NoError(t, err, "Item creation should succeed")

	// Test finding item by hash
	found, err := db.FindItemByHash(hash)
	assert.NoError(t, err, "FindItemByHash should succeed")
	assert.NotNil(t, found, "Item should be found")
	assert.Equal(t, hash, found.Hash, "Hash should match")
	assert.Equal(t, *item.Framing, *found.Framing, "Framing should match")
	assert.Equal(t, *item.TitleAI, *found.TitleAI, "TitleAI should match")
	assert.Equal(t, *item.ReasonAI, *found.ReasonAI, "ReasonAI should match")

	// Test non-existent hash
	_, err = db.FindItemByHash("nonexistent")
	assert.Error(t, err, "FindItemByHash should fail for non-existent hash")
}

func TestItemConstraints(t *testing.T) {
	db := setupTestDB(t)

	// Test unique constraint on Hash
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte("test")))
	item1 := &Item{
		Hash:        hash,
		FeedUrl:     "dummy4",
		Link:        "dummy4",
		Guid:        "dummy4",
		Title:       "dummy4",
		Description: "dummy4",
		Framing:     new(float32),
		TitleAI:     new(string),
		ReasonAI:    new(string),
	}
	*item1.Framing = 0.5
	*item1.TitleAI = "Title1"
	*item1.ReasonAI = "Reason1"
	err := db.CreateItem(item1)
	assert.NoError(t, err, "First item creation should succeed")

	item2 := &Item{
		Hash:        hash,
		FeedUrl:     "dummy5",
		Link:        "dummy5",
		Guid:        "dummy5",
		Title:       "dummy5",
		Description: "dummy5",
		Framing:     new(float32),
		TitleAI:     new(string),
		ReasonAI:    new(string),
	}
	*item2.Framing = 0.7
	*item2.TitleAI = "Title2"
	*item2.ReasonAI = "Reason2"
	err = db.CreateItem(item2)
	assert.NoError(t, err, "Creating item with duplicate Hash should succeed (ignored)")
}

func TestCreateCache(t *testing.T) {
	d := setupTestDB(t)

	cache := &Cache{
		FeedUrl: "https://example.com/rss",
		Cache:   "<rss>initial content</rss>",
	}

	// First insert
	err := d.CreateCache(cache)
	assert.NoError(t, err)

	var result Cache
	err = d.db.First(&result, "feed_url = ?", cache.FeedUrl).Error
	assert.NoError(t, err)
	assert.Equal(t, cache.FeedUrl, result.FeedUrl)
	assert.Equal(t, cache.Cache, result.Cache)

	// Update cache content and test upsert behavior
	newContent := "<rss>updated content</rss>"
	cache.Cache = newContent
	time.Sleep(time.Millisecond) // To ensure updated_at is different
	err = d.CreateCache(cache)
	assert.NoError(t, err)

	var updated Cache
	err = d.db.First(&updated, "feed_url = ?", cache.FeedUrl).Error
	assert.NoError(t, err)
	assert.Equal(t, newContent, updated.Cache)
	assert.True(t, updated.UpdatedAt.After(result.UpdatedAt), "UpdatedAt should be refreshed")
}

func TestFindCacheByFeedUrl(t *testing.T) {
	d := setupTestDB(t)

	// Insert a cache entry with UpdatedAt = now (done by gorm)
	cache := &Cache{
		FeedUrl: "https://example.com/rss",
		Cache:   "<rss>some content</rss>",
	}
	err := d.CreateCache(cache)
	assert.NoError(t, err)

	// Refresh cache from DB to get accurate UpdatedAt timestamp
	var inserted Cache
	err = d.db.First(&inserted, "feed_url = ?", cache.FeedUrl).Error
	assert.NoError(t, err)

	// Case 1: maxAge large enough to include the entry
	found, err := d.FindCacheByFeedUrl(cache.FeedUrl, time.Minute*5)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, cache.FeedUrl, found.FeedUrl)

	// Case 2: maxAge too small, entry should NOT be found
	found, err = d.FindCacheByFeedUrl(cache.FeedUrl, time.Nanosecond)
	assert.Error(t, err)
	assert.Nil(t, found)

	// Case 3: Non-existent FeedUrl returns error
	found, err = d.FindCacheByFeedUrl("nonexistent", time.Minute*5)
	assert.Error(t, err)
	assert.Nil(t, found)
}
