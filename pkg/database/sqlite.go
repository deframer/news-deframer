package database

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Item represents the rss items
type Item struct {
	gorm.Model
	Hash        string   `gorm:"type:text;uniqueIndex;not null"` // SHA-256 hash with unique index
	FeedUrl     string   `gorm:"type:text;not null"`
	Link        string   `gorm:"type:text;not null"`
	Guid        string   `gorm:"type:text;not null"`
	Title       string   `gorm:"type:text;not null"`
	Description string   `gorm:"type:text;not null"`
	Framing     *float32 `gorm:"type:real"` // Nullable
	TitleAI     *string  `gorm:"type:text"` // Nullable
	ReasonAI    *string  `gorm:"type:text"` // Nullable
}

// Cache represents the cached feed
type Cache struct {
	gorm.Model
	FeedUrl string `gorm:"type:text;uniqueIndex;not null"`
	Cache   string `gorm:"type:text;not null"`
}

// Database handles DB operations
type Database struct {
	db *gorm.DB
}

// NewDatabase initializes a new SQLite database
func NewDatabase(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate to create table with constraints
	err = db.AutoMigrate(&Item{}, &Cache{})
	if err != nil {
		return nil, err
	}

	// Explicitly ensure unique index on Hash
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_hash ON items(hash)").Error
	if err != nil {
		return nil, err
	}

	// Explicitly ensure unique index on FeedUrl
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_feed_url ON caches(feed_url)").Error
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// CreateItem inserts a new item, ignores if hash already exists
func (d *Database) CreateItem(item *Item) error {
	now := time.Now()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = now
	}

	return d.db.Exec(
		`INSERT OR IGNORE INTO items
		(hash, feed_url, link, guid, title, description, framing, title_ai, reason_ai, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.Hash, item.FeedUrl, item.Link, item.Guid, item.Title, item.Description,
		item.Framing, item.TitleAI, item.ReasonAI, item.CreatedAt, item.UpdatedAt,
	).Error
}

// FindItemByHash retrieves an item by its hash
func (d *Database) FindItemByHash(hash string) (*Item, error) {
	var item Item
	result := d.db.Where("hash = ?", hash).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

// CreateCache inserts or replaces a cache entry for the given FeedUrl.
func (d *Database) CreateCache(cache *Cache) error {
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "feed_url"}}, // unique constraint column
		UpdateAll: true,                                // update all fields on conflict
	}).Create(cache).Error
}

func (d *Database) FindCacheByFeedUrl(feedUrl string, maxAge time.Duration) (*Cache, error) {
	var cache Cache
	cutoff := time.Now().Add(-maxAge)
	err := d.db.
		Where("feed_url = ? AND updated_at >= ?", feedUrl, cutoff).
		First(&cache).Error
	if err != nil {
		return nil, err
	}
	return &cache, nil
}
