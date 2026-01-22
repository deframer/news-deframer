package database

import (
	"fmt"

	"github.com/deframer/news-deframer/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connects to the database and runs the migration
func RunMigrations(cfg *config.Config) error {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return Migrate(db)
}

// Migrate updates the database schema and seeds initial data.
func Migrate(db *gorm.DB) error {
	// Acquire an advisory lock to prevent concurrent migrations from causing race conditions
	// (like "duplicate key value violates unique constraint pg_type_typname_nsp_index").
	// 123456789 is an arbitrary lock ID.
	if err := db.Exec("SELECT pg_advisory_lock(123456789)").Error; err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	defer db.Exec("SELECT pg_advisory_unlock(123456789)")

	// Ensure the Extension for UUIDs exists
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return fmt.Errorf("failed to create extension uuid-ossp: %w", err)
	}

	// 1. AutoMigrate the schema
	if err := db.AutoMigrate(&Feed{}, &CachedFeed{}, &Item{}, &FeedSchedule{}); err != nil {
		return err
	}

	// 2. Manually add FK for CachedFeed -> Feed (Shared PK)
	// Workaround for GORM AutoMigrate issue:
	// GORM struggles with Shared PK + Belongs To direction inference,
	// causing circular dependency errors during table creation if defined in the struct.
	_ = db.Exec(`DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_cached_feeds_feeds') THEN
			ALTER TABLE cached_feeds ADD CONSTRAINT fk_cached_feeds_feeds FOREIGN KEY (id) REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_feed_schedules_feeds') THEN
			ALTER TABLE feed_schedules ADD CONSTRAINT fk_feed_schedules_feeds FOREIGN KEY (id) REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE;
		END IF;
	END $$;`)

	return nil
	// // 3. Seed Dummy Data
	// return seed(db)
}

// func seed(db *gorm.DB) error {
// 	var count int64
// 	if err := db.Model(&Feed{}).Count(&count).Error; err != nil {
// 		return err
// 	}

// 	if count > 0 {
// 		return nil
// 	}

// 	slog.Info("Seeding dummy feeds...")
// 	idWordPressInternal := uuid.New()

// 	feeds := []Feed{
// 		{Base: Base{ID: idWordPressInternal}, URL: "http://wordpress/feed", Enabled: true, EnforceFeedDomain: true, Polling: true},
// 	}

// 	if err := db.Create(&feeds).Error; err != nil {
// 		return err
// 	}

// 	// Seed FeedSchedules
// 	var schedules []FeedSchedule
// 	now := time.Now()
// 	for _, f := range feeds {
// 		fs := FeedSchedule{
// 			ID: f.ID,
// 		}

// 		// limit to the internal wordpress for debugging
// 		if f.ID != idWordPressInternal {
// 			continue
// 		}

// 		if f.Enabled && f.Polling && !f.DeletedAt.Valid {
// 			fs.NextRunAt = &now
// 		}
// 		schedules = append(schedules, fs)
// 	}

// 	return db.Create(&schedules).Error
// }
