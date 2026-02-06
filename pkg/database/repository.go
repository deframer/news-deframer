package database

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	mylogger "github.com/deframer/news-deframer/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

//go:embed sql/statement/top_trend_by_domain.sql
var topTrendByDomainQuery string

//go:embed sql/statement/context_by_domain.sql
var contextByDomainQuery string

//go:embed sql/statement/lifecycle_by_domain.sql
var lifecycleByDomainQuery string

type TrendMetric struct {
	TrendTopic   string    `gorm:"column:trend_topic" json:"trend_topic"`
	Frequency    int64     `gorm:"column:frequency" json:"frequency"`
	Utility      int64     `gorm:"column:utility" json:"utility"`
	OutlierRatio float64   `gorm:"column:outlier_ratio" json:"outlier_ratio"`
	TimeSlice    time.Time `gorm:"column:time_slice" json:"time_slice"`
}

type TrendContext struct {
	Context   string `gorm:"column:context_word" json:"context"`
	Frequency int64  `gorm:"column:frequency" json:"frequency"`
}

type Lifecycle struct {
	TimeSlice time.Time `gorm:"column:time_slice" json:"time_slice"`
	Frequency int64     `gorm:"column:frequency" json:"frequency"`
	Velocity  int64     `gorm:"column:velocity" json:"velocity"`
}

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
	FindItemsByRootDomain(rootDomain string, limit int) ([]Item, error)
	GetAllFeeds(deleted bool) ([]Feed, error)
	DeleteFeedById(id uuid.UUID) error
	PurgeFeedById(id uuid.UUID) error
	EnqueueSync(id uuid.UUID, pollingInterval time.Duration) error
	EnqueueMine(id uuid.UUID, miningInterval time.Duration) error
	RemoveMine(id uuid.UUID) error
	RemoveSync(id uuid.UUID) error
	BeginFeedUpdate(lockDuration time.Duration) (*Feed, error)
	EndFeedUpdate(id uuid.UUID, jobErr error, pollingInterval time.Duration) error
	GetPendingItems(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error)
	GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]Item, error)
	UpsertItem(item *Item) error
	UpsertCachedFeed(cachedFeed *CachedFeed) error
	FindCachedFeedById(feedID uuid.UUID) (*CachedFeed, error)
	FindFeedScheduleById(feedID uuid.UUID) (*FeedSchedule, error)
	CreateFeedSchedule(feedID uuid.UUID) error
	GetTopTrendByDomain(domain string, language string, daysInPast int) ([]TrendMetric, error)
	GetContextByDomain(term string, domain string, language string, daysInPast int) ([]TrendContext, error)
	GetLifecycleByDomain(term string, domain string, language string, daysInPast int) ([]Lifecycle, error)
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

	if err := Migrate(db, false); err != nil {
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
	if item.Categories == nil {
		item.Categories = []string{}
	}
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

		// Mitigation: check for URL conflict within the same feed (different hash)
		var urlConflict Item
		if err := tx.Where("feed_id = ? AND url = ? AND hash != ?", item.FeedID, item.URL, item.Hash).First(&urlConflict).Error; err == nil {
			slog.Warn("URL conflict detected, deleting old version to allow update", "url", item.URL, "old_hash", urlConflict.Hash, "new_hash", item.Hash)
			if err := tx.Where("item_id = ?", urlConflict.ID).Delete(&Trend{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&urlConflict).Error; err != nil {
				return err
			}
		}

		return tx.Save(item).Error
	})
}

func (r *repository) EnqueueSync(id uuid.UUID, pollingInterval time.Duration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.enqueueSyncTx(tx, id, pollingInterval)
	})
}

func (r *repository) EnqueueMine(id uuid.UUID, miningInterval time.Duration) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.enqueueMineTx(tx, id, miningInterval)
	})
}

func (r *repository) enqueueSyncTx(tx *gorm.DB, id uuid.UUID, pollingInterval time.Duration) error {
	var feed Feed
	// Check if feed exists, is enabled, and is not deleted.
	if err := tx.Where("id = ? AND enabled = ?", id, true).First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("feed not found or unavailable")
		}
		return err
	}

	runInterval := fmt.Sprintf("INTERVAL '%f seconds'", pollingInterval.Seconds())

	// Upsert FeedSchedule
	var schedule FeedSchedule
	if err := tx.Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Model(&FeedSchedule{}).Create(map[string]interface{}{
				"id":              id,
				"next_thinker_at": gorm.Expr("NOW() + " + runInterval),
				"updated_at":      gorm.Expr("NOW()"),
			}).Error
		}
		return err
	}

	return tx.Model(&schedule).Updates(map[string]interface{}{
		"next_thinker_at": gorm.Expr("NOW() + " + runInterval),
		"updated_at":      gorm.Expr("NOW()"),
	}).Error
}

func (r *repository) enqueueMineTx(tx *gorm.DB, id uuid.UUID, miningInterval time.Duration) error {
	var feed Feed
	// Check if feed exists, is enabled, and is not deleted.
	if err := tx.Where("id = ? AND enabled = ?", id, true).First(&feed).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("feed not found or unavailable")
		}
		return err
	}

	runInterval := fmt.Sprintf("INTERVAL '%f seconds'", miningInterval.Seconds())

	// Upsert FeedSchedule
	var schedule FeedSchedule
	if err := tx.Where("id = ?", id).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Model(&FeedSchedule{}).Create(map[string]interface{}{
				"id":             id,
				"next_mining_at": gorm.Expr("NOW() + " + runInterval),
				"updated_at":     gorm.Expr("NOW()"),
			}).Error
		}
		return err
	}

	return tx.Model(&schedule).Updates(map[string]interface{}{
		"next_mining_at": gorm.Expr("NOW() + " + runInterval),
		"updated_at":     gorm.Expr("NOW()"),
	}).Error
}

func (r *repository) RemoveMine(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.removeMineTx(tx, id)
	})
}

func (r *repository) removeMineTx(tx *gorm.DB, id uuid.UUID) error {
	var count int64
	if err := tx.Model(&FeedSchedule{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	return tx.Model(&FeedSchedule{}).Where("id = ?", id).Updates(map[string]interface{}{
		"next_mining_at":      nil,
		"mining_locked_until": nil,
		"updated_at":          gorm.Expr("NOW()"),
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
		"next_thinker_at":      nil,
		"thinker_locked_until": nil,
		"updated_at":           gorm.Expr("NOW()"),
	}).Error
}

func (r *repository) BeginFeedUpdate(lockDuration time.Duration) (*Feed, error) {
	var feed *Feed
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var schedule FeedSchedule
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("next_thinker_at <= NOW()").
			Where("thinker_locked_until IS NULL OR thinker_locked_until < NOW()"). // another worker / thread
			Order("next_thinker_at ASC").                                          // the oldest not scheduled
			First(&schedule).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		var f Feed
		if err := tx.Unscoped().First(&f, schedule.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// soft deleted
				slog.Warn("Feed schedule exists but feed not found, removing sync", "feed_id", schedule.ID)
				return r.removeSyncTx(tx, schedule.ID)
			}
			return err
		}

		if !f.Enabled || f.DeletedAt.Valid {
			// the feed might was disabled / delete
			slog.Warn("Feed schedule exists but feed is disabled or deleted, removing sync", "feed_id", schedule.ID)
			return r.removeSyncTx(tx, schedule.ID)
		}

		interval := fmt.Sprintf("INTERVAL '%f seconds'", lockDuration.Seconds())
		// Lock the feed
		if err := tx.Model(&schedule).Updates(map[string]interface{}{
			"thinker_locked_until": gorm.Expr("NOW() + " + interval), // lock
			"updated_at":           gorm.Expr("NOW()"),
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
			"thinker_locked_until": nil,
			"updated_at":           gorm.Expr("NOW()"),
		}

		if jobErr != nil {
			updates["next_thinker_at"] = nil
		} else {
			if !feed.Polling {
				updates["next_thinker_at"] = nil
			}
		}

		if err := tx.Model(&FeedSchedule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if jobErr == nil && feed.Enabled && feed.Polling {
			return r.enqueueSyncTx(tx, id, pollingInterval)
		}

		return nil
	})
}

func (r *repository) GetPendingItems(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error) {
	// Initialize map with all hashes having count 0 (assuming they are new)
	pendingItems := make(map[string]int)
	for _, h := range hashes {
		pendingItems[h] = 0
	}

	const batchSize = 1000

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(hashes); i += batchSize {
			end := i + batchSize
			if end > len(hashes) {
				end = len(hashes)
			}
			batch := hashes[i:end]

			var foundItems []struct {
				Hash            string
				ThinkResult     *ThinkResult
				ThinkErrorCount int
			}

			// We fetch items that exist to check their status
			if err := tx.Model(&Item{}).
				Select("hash, think_result, think_error_count").
				Where("feed_id = ?", feedID).
				Where("hash IN ?", batch).
				Scan(&foundItems).Error; err != nil {
				return err
			}

			for _, item := range foundItems {
				// If processed (ThinkResult not nil) OR failed too many times, remove from pending
				if item.ThinkResult != nil || item.ThinkErrorCount > maxRetries {
					delete(pendingItems, item.Hash)
				} else {
					// It exists but needs retry (or initial processing if error count is small)
					// Update the error count so the caller knows the current state
					pendingItems[item.Hash] = item.ThinkErrorCount
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pendingItems, nil
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

func (r *repository) CreateFeedSchedule(feedID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&FeedSchedule{}).Where("id = ?", feedID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		return tx.Create(&FeedSchedule{ID: feedID}).Error
	})
}

func (r *repository) DeleteFeedById(id uuid.UUID) error {
	return r.db.Delete(&Feed{Base: Base{ID: id}}).Error
}

func (r *repository) PurgeFeedById(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Delete CachedFeed (1:1 with Feed ID)
		if err := tx.Where("id = ?", id).Delete(&CachedFeed{}).Error; err != nil {
			return fmt.Errorf("failed to delete cached feed: %w", err)
		}
		// 2. Delete FeedSchedule (1:1 with Feed ID)
		if err := tx.Where("id = ?", id).Delete(&FeedSchedule{}).Error; err != nil {
			return fmt.Errorf("failed to delete feed schedule: %w", err)
		}
		// Delete Trends (N:1 with Feed ID)
		if err := tx.Where("feed_id = ?", id).Delete(&Trend{}).Error; err != nil {
			return fmt.Errorf("failed to delete trends: %w", err)
		}
		// 3. Delete Items (N:1 with Feed ID)
		if err := tx.Where("feed_id = ?", id).Delete(&Item{}).Error; err != nil {
			return fmt.Errorf("failed to delete items: %w", err)
		}
		// 4. Hard Delete Feed
		if err := tx.Unscoped().Delete(&Feed{}, "id = ?", id).Error; err != nil {
			return fmt.Errorf("failed to delete feed: %w", err)
		}
		return nil
	})
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

func (r *repository) GetTopTrendByDomain(domain string, language string, daysInPast int) ([]TrendMetric, error) {
	var metrics []TrendMetric

	if daysInPast < 1 {
		daysInPast = 1
	}
	if daysInPast > 365 {
		daysInPast = 365
	}

	if err := r.db.Raw(topTrendByDomainQuery,
		sql.Named("language", language),
		sql.Named("domain", domain),
		sql.Named("days_in_past", daysInPast),
	).Scan(&metrics).Error; err != nil {
		return nil, err
	}
	return metrics, nil
}

func (r *repository) GetContextByDomain(term string, domain string, language string, daysInPast int) ([]TrendContext, error) {
	var contexts []TrendContext

	if daysInPast < 1 {
		daysInPast = 1
	}
	if daysInPast > 365 {
		daysInPast = 365
	}

	if err := r.db.Raw(contextByDomainQuery,
		sql.Named("term", term),
		sql.Named("domain", domain),
		sql.Named("language", language),
		sql.Named("days_in_past", daysInPast),
	).Scan(&contexts).Error; err != nil {
		return nil, err
	}
	return contexts, nil
}

func (r *repository) GetLifecycleByDomain(term string, domain string, language string, daysInPast int) ([]Lifecycle, error) {
	var items []Lifecycle

	if daysInPast < 1 {
		daysInPast = 1
	}
	if daysInPast > 365 {
		daysInPast = 365
	}

	if err := r.db.Raw(lifecycleByDomainQuery,
		sql.Named("term", term),
		sql.Named("domain", domain),
		sql.Named("language", language),
		sql.Named("days_in_past", daysInPast),
	).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindItemsByRootDomain retrieves the most recent items from all feeds belonging to the given root domain.
// It enforces the domain boundary via a JOIN on the feeds table and filters by the root_domain column.
// It ensures that items with the same URL are deduplicated, keeping the most recent one.
// Note: This aggregates items from all matching feeds and sorts them globally by publication date.
func (r *repository) FindItemsByRootDomain(rootDomain string, limit int) ([]Item, error) {
	var items []Item
	subQuery := r.db.Select("DISTINCT ON (items.url) items.*").
		Table("items").
		Joins("JOIN feeds ON feeds.id = items.feed_id").
		Where("feeds.root_domain = ? AND feeds.enabled = ? AND feeds.deleted_at IS NULL", rootDomain, true).
		Where("items.think_result IS NOT NULL AND items.think_error IS NULL AND items.think_error_count = 0").
		Order("items.url, items.pub_date DESC")

	if err := r.db.Table("(?) as unique_items", subQuery).
		Order("unique_items.pub_date DESC").
		Limit(limit).
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
