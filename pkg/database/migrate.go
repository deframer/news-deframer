package database

import (
	"fmt"
	"log/slog"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
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

	idEnforced := uuid.New()
	idOpen := uuid.New()

	feeds := []Feed{
		{URL: "http://rssbridge/?action=display&bridge=TheVerge&format=Atom", Enabled: true, EnforceFeedDomain: true, AutoPolling: false},
		{Base: Base{ID: idEnforced}, URL: "http://dummy-enforced", Enabled: true, EnforceFeedDomain: true, AutoPolling: false},
		{Base: Base{ID: idOpen}, URL: "http://dummy-open", Enabled: true, EnforceFeedDomain: false, AutoPolling: false},
		{URL: "http://wordpress/feed", Enabled: true, EnforceFeedDomain: true, AutoPolling: true},
		{URL: "http://localhost:8003/feed", Enabled: true, EnforceFeedDomain: true, AutoPolling: true},
	}

	if err := db.Create(&feeds).Error; err != nil {
		return err
	}

	slog.Info("Seeding dummy items...")
	items := []Item{
		// Enforced Feed: Item pointing to itself
		{Hash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", FeedID: idEnforced, URL: "http://dummy-enforced/item-1", AnalyzerResult: JSONB{"title": "Enforced Item 1"}},
		// Open Feed: Item pointing to itself
		{Hash: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", FeedID: idOpen, URL: "http://dummy-open/item-2", AnalyzerResult: JSONB{"title": "Open Item 2"}},
		// Open Feed: Item pointing to Enforced Domain (Syndication)
		// Same content hash as item 1, same URL as item 1, but different FeedID
		{Hash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", FeedID: idOpen, URL: "http://dummy-enforced/item-1", AnalyzerResult: JSONB{"title": "Syndicated Item 1"}},
	}

	if err := db.Create(&items).Error; err != nil {
		return err
	}

	// add a dummy cached read for  idEnforced and idOpen
	cachedFeeds := []CachedFeed{
		{
			ID:        idEnforced,
			XMLHeader: `<rss version="2.0"><channel><title>Dummy Enforced</title></channel></rss>`,
			ItemRefs:  StringArray{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		},
		{
			ID:        idOpen,
			XMLHeader: `<rss version="2.0"><channel><title>Dummy Open</title></channel></rss>`,
			ItemRefs:  StringArray{"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		},
	}

	return db.Create(&cachedFeeds).Error
}
