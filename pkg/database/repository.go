package database

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	mylogger "github.com/egandro/news-deframer/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Repository interface {
	FindFeedByUrl(u *url.URL) (*Feed, error)
	FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*Feed, error)
	FindFeedById(feedID uuid.UUID) (*Feed, error)
	UpsertFeed(feed *Feed) error
	// FindItemsByUrl retrieves all items associated with a specific URL.
	// As per the specification (Syndication), a single URL can legitimately appear in multiple feeds.
	// Therefore, the system allows multiple Item records for the same URL, distinguished by their FeedID.
	// This means the same content (probably different Hashes) can exist multiple times if it is syndicated across different feeds.
	// Feed.EnforceFeedDomain = will enforce only items with the same base domain as the Feed URL
	FindItemsByUrl(u *url.URL) ([]Item, error)
	GetAllFeeds(deleted bool) ([]Feed, error)
	DeleteFeedById(id uuid.UUID) error
	EnqueueSync(id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error
	RemoveSync(id uuid.UUID) error
	BeginFeedUpdate(lockDuration time.Duration) (*Feed, error)
	EndFeedUpdate(id uuid.UUID, jobErr error, pollingInterval time.Duration) error
	GetPendingHashes(feedID uuid.UUID, hashes []string) (map[string]bool, error)
	GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]Item, error)
	UpsertItem(item *Item) error
	UpsertCachedFeed(cachedFeed *CachedFeed) error
	FindCachedFeedById(feedID uuid.UUID) (*CachedFeed, error)
	FindFeedScheduleById(feedID uuid.UUID) (*FeedSchedule, error)
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
		Logger: mylogger.NewSlogGormLogger(slog.Default(), lvl),
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

func (r *repository) FindFeedByUrl(u *url.URL) (*Feed, error) {
	return r.FindFeedByUrlAndAvailability(u, true)
}

func (r *repository) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*Feed, error) {
	var feed Feed
	query := r.db.Where("url = ?", u.String())
	if onlyEnabled {
		query = query.Where("enabled = ?", true)
	}
	if err := query.First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feed, nil
}

func (r *repository) FindFeedById(feedID uuid.UUID) (*Feed, error) {
	var feed Feed
	if err := r.db.Unscoped().Where("id = ?", feedID).First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &feed, nil
}

func (r *repository) UpsertFeed(feed *Feed) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		// Check for existing active feed with same URL
		// GORM automatically adds `deleted_at IS NULL` for models with DeletedAt
		query := tx.Model(&Feed{}).Where("url = ?", feed.URL)

		if feed.ID != uuid.Nil {
			query = query.Where("id != ?", feed.ID)
		}

		if err := query.Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			return fmt.Errorf("feed with url %q already exists", feed.URL)
		}

		return tx.Save(feed).Error
	})
}

func (r *repository) UpsertItem(item *Item) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if item.ID == uuid.Nil {
			var existing Item
			if err := tx.Select("id", "created_at").Where("feed_id = ? AND hash = ?", item.FeedID, item.Hash).First(&existing).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
			} else {
				item.ID = existing.ID
				item.CreatedAt = existing.CreatedAt
			}
		}
		return tx.Save(item).Error
	})
}

func (r *repository) EnqueueSync(id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.enqueueSyncTx(tx, id, pollingInterval, lockDuration)
	})
}

func (r *repository) enqueueSyncTx(tx *gorm.DB, id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error {
	var feed Feed
	// Check if feed exists, is enabled, and is not deleted.
	if err := tx.Where("id = ? AND enabled = ?", id, true).First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("feed not found or unavailable")
		}
		return err
	}

	if lockDuration > 0 {
		lockInterval := fmt.Sprintf("INTERVAL '%f seconds'", lockDuration.Seconds())

		// Check if the feed is locked for longer than the requested duration
		var count int64
		if err := tx.Model(&FeedSchedule{}).
			Where("id = ? AND locked_until > NOW() + "+lockInterval, id).
			Count(&count).Error; err != nil {
			return err
		}

		if count > 0 {
			// Lock is > now + duration, do nothing
			return nil
		}
	}

	runInterval := fmt.Sprintf("INTERVAL '%f seconds'", pollingInterval.Seconds())

	// Upsert FeedSchedule
	var schedule FeedSchedule
	if err := tx.Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Model(&FeedSchedule{}).Create(map[string]interface{}{
				"id":           id,
				"next_run_at":  gorm.Expr("NOW() + " + runInterval),
				"locked_until": nil,
				"last_error":   nil,
				"updated_at":   gorm.Expr("NOW()"),
			}).Error
		}
		return err
	}

	return tx.Model(&schedule).Updates(map[string]interface{}{
		"next_run_at":  gorm.Expr("NOW() + " + runInterval),
		"locked_until": nil,
		"last_error":   nil,
		"updated_at":   gorm.Expr("NOW()"),
	}).Error
}

func (r *repository) RemoveSync(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.removeSyncTx(tx, id)
	})
}

func (r *repository) removeSyncTx(tx *gorm.DB, id uuid.UUID) error {
	var count int64
	if err := tx.Model(&FeedSchedule{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	return tx.Model(&FeedSchedule{}).Where("id = ?", id).Updates(map[string]interface{}{
		"next_run_at":  nil,
		"locked_until": nil,
		"last_error":   nil,
		"updated_at":   gorm.Expr("NOW()"),
	}).Error
}

func (r *repository) BeginFeedUpdate(lockDuration time.Duration) (*Feed, error) {
	var feed *Feed
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var schedule FeedSchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("next_run_at <= NOW()").
			Where("locked_until IS NULL OR locked_until < NOW()"). // another worker / thread
			Order("next_run_at ASC").                              // the oldest not scheduled
			First(&schedule).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		var f Feed
		if err := tx.Unscoped().First(&f, schedule.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return r.removeSyncTx(tx, schedule.ID)
			}
			return err
		}

		if !f.Enabled || f.DeletedAt.Valid {
			// the feed might was disabled / delete
			return r.removeSyncTx(tx, schedule.ID)
		}

		interval := fmt.Sprintf("INTERVAL '%f seconds'", lockDuration.Seconds())
		// Lock the feed
		if err := tx.Model(&schedule).Updates(map[string]interface{}{
			"locked_until": gorm.Expr("NOW() + " + interval), // lock
			"updated_at":   gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}

		feed = &f
		return nil
	})
	return feed, err
}

func (r *repository) EndFeedUpdate(id uuid.UUID, jobErr error, pollingInterval time.Duration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var feed Feed
		if err := tx.Where("id = ?", id).First(&feed).Error; err != nil {
			return nil
		}

		updates := map[string]interface{}{
			"locked_until": nil,
			"updated_at":   gorm.Expr("NOW()"),
		}

		if jobErr != nil {
			updates["last_error"] = jobErr.Error()
			updates["next_run_at"] = nil
		} else {
			updates["last_error"] = nil
			if !feed.Polling {
				updates["next_run_at"] = nil
			}
		}

		if err := tx.Model(&FeedSchedule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if jobErr == nil && feed.Enabled && feed.Polling {
			return r.enqueueSyncTx(tx, id, pollingInterval, 0)
		}

		return nil
	})
}

func (r *repository) GetPendingHashes(feedID uuid.UUID, hashes []string) (map[string]bool, error) {
	processedHashes := make(map[string]bool)
	const batchSize = 1000

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hashes); i += batchSize {
			end := i + batchSize
			if end > len(hashes) {
				end = len(hashes)
			}
			batch := hashes[i:end]

			var foundHashes []string
			// We are looking for items that are already processed (AnalyzerResult is NOT NULL)
			// to exclude them from the pending list.
			if err := tx.Model(&Item{}).
				Where("feed_id = ?", feedID).
				Where("hash IN ?", batch).
				Where("think_result IS NOT NULL").
				Pluck("hash", &foundHashes).Error; err != nil {
				return err
			}

			for _, h := range foundHashes {
				processedHashes[h] = true
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, h := range hashes {
		if !processedHashes[h] {
			result[h] = true
		}
	}
	return result, nil
}

func (r *repository) GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]Item, error) {
	var items []Item
	const batchSize = 1000

	if len(hashes) == 0 {
		return items, nil
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hashes); i += batchSize {
			end := i + batchSize
			if end > len(hashes) {
				end = len(hashes)
			}
			batch := hashes[i:end]

			var batchItems []Item
			if err := tx.Where("feed_id = ?", feedID).
				Where("hash IN ?", batch).
				Find(&batchItems).Error; err != nil {
				return err
			}
			items = append(items, batchItems...)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repository) UpsertCachedFeed(cachedFeed *CachedFeed) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing CachedFeed
		if err := tx.Select("created_at").Where("id = ?", cachedFeed.ID).First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			cachedFeed.CreatedAt = existing.CreatedAt
		}
		return tx.Save(cachedFeed).Error
	})
}

func (r *repository) FindCachedFeedById(feedID uuid.UUID) (*CachedFeed, error) {
	var cached CachedFeed
	if err := r.db.Where("id = ?", feedID).First(&cached).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cached, nil
}

func (r *repository) FindFeedScheduleById(feedID uuid.UUID) (*FeedSchedule, error) {
	var schedule FeedSchedule
	if err := r.db.Where("id = ?", feedID).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *repository) DeleteFeedById(id uuid.UUID) error {
	return r.db.Delete(&Feed{Base: Base{ID: id}}).Error
}

func (r *repository) GetAllFeeds(deleted bool) ([]Feed, error) {
	var feeds []Feed
	query := r.db.Preload("FeedSchedule").Order("url ASC")
	if deleted {
		query = query.Unscoped()
	}
	if err := query.Find(&feeds).Error; err != nil {
		return nil, err
	}
	return feeds, nil
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
