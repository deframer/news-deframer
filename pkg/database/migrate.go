package database

import (
	"fmt"
	"log/slog"

	"github.com/egandro/news-deframer/pkg/config"
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
	// 1. AutoMigrate the schema
	if err := db.AutoMigrate(&Feed{}, &CachedFeed{}, &Item{}); err != nil {
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
	END $$;`)

	// 3. Seed Dummy Data
	return seed(db)
}

func seed(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Feed{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	slog.Info("Seeding dummy feeds...")
	feeds := []Feed{
		{URL: "http://rssbridge/?action=display&bridge=TheVerge&format=Atom", Enabled: true, EnforceFeedDomain: true},
		{URL: "http://dummy", Enabled: true, EnforceFeedDomain: true},
		{URL: "http://wordpress/feed", Enabled: true, EnforceFeedDomain: true},
		{URL: "http://localhost:8003/feed", Enabled: true, EnforceFeedDomain: true},
	}

	return db.Create(&feeds).Error
}
