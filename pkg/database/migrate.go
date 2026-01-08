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

	// 2. Seed Dummy Data
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
		{URL: "http://rssbridge/?action=display&bridge=TheVerge&format=Atom", Enabled: true},
		{URL: "http://dummy", Enabled: true},
		{URL: "http://wordpress/feed", Enabled: true},
		{URL: "http://localhost:8003/feed", Enabled: true},
	}

	return db.Create(&feeds).Error
}
