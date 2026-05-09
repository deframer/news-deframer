package main

import (
	"context"
	"os"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"goa.design/clue/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(context.Background(), err, "Failed to load config")
		os.Exit(1)
	}

	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	cfg.DatabaseLogging = true

	// Print the DSN as requested
	log.Print(ctx, log.KV{K: "dsn", V: cfg.DSN})

	if err := database.RunMigrations(cfg, true); err != nil {
		log.Fatalf(ctx, err, "Failed to run migrations")
		os.Exit(1)
	}

	log.Print(ctx, log.KV{K: "message", V: "Migrations completed successfully"})
}
