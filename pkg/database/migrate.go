package database

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/deframer/news-deframer/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed sql
var migrationFS embed.FS

// Connects to the database and runs the migration
func RunMigrations(cfg *config.Config, forced bool) error {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	return Migrate(db, forced)
}

// Migrate updates the database schema and seeds initial data.
func Migrate(db *gorm.DB, forced bool) error {
	// If not forced, check if the main table exists. If so, assume migrated.
	if !forced && db.Migrator().HasTable(&Feed{}) {
		return nil
	}

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

	// Ensure the Extension for pg_duckdb exists
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pg_duckdb";`).Error; err != nil {
		return fmt.Errorf("failed to create extension pg_duckdb: %w. Are you using the pgduckdb/pgduckdb docker image or a postgres build with pg_duckdb? Check https://github.com/duckdb/pg_duckdb and https://pgxman.com/x/pg_duckdb", err)
	}

	// AutoMigrate the schema
	if err := db.AutoMigrate(&Feed{}, &CachedFeed{}, &Item{}, &FeedSchedule{}, &Trend{}); err != nil {
		return err
	}

	// Manually add FK for CachedFeed -> Feed (Shared PK)
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
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_trends_items') THEN
		END IF;
		IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_trends_feeds') THEN
			ALTER TABLE trends ADD CONSTRAINT fk_trends_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE ON UPDATE CASCADE;
		END IF;
	END $$;`)

	// Run embedded SQL files
	if err := migrateViews(db); err != nil {
		return err
	}

	return nil
}

func migrateViews(db *gorm.DB) error {
	entries, err := migrationFS.ReadDir("sql/views")
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		content, err := migrationFS.ReadFile("sql/views/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration file %s: %w", entry.Name(), err)
		}
	}

	return nil
}
