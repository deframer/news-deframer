package database

import (
	"fmt"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	GetTime() (time.Time, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository initializes a new repository with a database connection from config.
func NewRepository(cfg *config.Config) (Repository, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// if err := Migrate(db); err != nil {
	// 	return nil, fmt.Errorf("failed to migrate database: %w", err)
	// }

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
