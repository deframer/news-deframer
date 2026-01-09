package database

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	pkglogger "github.com/egandro/news-deframer/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	GetTime() (time.Time, error)
	FindFeedByUrl(u *url.URL) (*Feed, error)
	FindCachedFeedById(feedID uuid.UUID) (*CachedFeed, error)
	FindItemsByCachedFeedFeedId(feedID uuid.UUID) ([]Item, error)
	// FindItemsByUrl retrieves all items associated with a specific URL.
	// As per the specification (Syndication), a single URL can legitimately appear in multiple feeds.
	// Therefore, the system allows multiple Item records for the same URL, distinguished by their FeedID.
	// This means the same content (probably different Hashes) can exist multiple times if it is syndicated across different feeds.
	// Feed.EnforceFeedDomain = will enforce only items with the same base domain as the Feed URL
	FindItemsByUrl(u *url.URL) ([]Item, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository initializes a new repository with a database connection from config.
func NewRepository(cfg *config.Config) (Repository, error) {
	lvl := logger.Error // maybe Silent is better
	if cfg.LogDatabase {
		lvl = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: pkglogger.NewSlogGormLogger(slog.Default(), lvl),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := Migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &repository{db: db}, nil
}

// NewFromDB creates a repository from an existing GORM connection.
func NewFromDB(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetTime() (time.Time, error) {
	var t time.Time
	// Simple query to get DB time
	result := r.db.Raw("SELECT NOW()").Scan(&t)
	return t, result.Error
}

func (r *repository) FindFeedByUrl(u *url.URL) (*Feed, error) {
	var feed Feed
	if err := r.db.Where("url = ? AND enabled = ?", u.String(), true).First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feed, nil
}

func (r *repository) FindCachedFeedById(feedID uuid.UUID) (*CachedFeed, error) {
	var cachedFeed CachedFeed
	if err := r.db.Joins("JOIN feeds ON feeds.id = cached_feeds.id").
		Where("cached_feeds.id = ? AND feeds.enabled = ? AND feeds.deleted_at IS NULL", feedID, true).
		First(&cachedFeed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cachedFeed, nil
}

func (r *repository) FindItemsByCachedFeedFeedId(feedID uuid.UUID) ([]Item, error) {
	cachedFeed, err := r.FindCachedFeedById(feedID)
	if err != nil {
		return nil, err
	}
	if cachedFeed == nil || len(cachedFeed.ItemRefs) == 0 {
		return []Item{}, nil
	}

	var items []Item
	if err := r.db.Where("feed_id = ? AND hash IN ?", feedID, []string(cachedFeed.ItemRefs)).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) FindItemsByUrl(u *url.URL) ([]Item, error) {
	var items []Item
	// We query by URL. Since URL is not unique (unique only per Feed), this returns a list.
	if err := r.db.Joins("JOIN feeds ON feeds.id = items.feed_id").
		Where("items.url = ? AND feeds.enabled = ? AND feeds.deleted_at IS NULL", u.String(), true).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
