package main

import (
	"log/slog"
	"os"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// force this
	cfg.LogLevel = "debug"
	cfg.LogDatabase = true

	// Print the DSN as requested
	slog.Info("Running migrations", "dsn", cfg.DSN)

	if err := database.RunMigrations(cfg, true); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("Migrations completed successfully")
}
