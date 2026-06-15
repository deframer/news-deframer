package main

import (
	"context"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/joho/godotenv"
	"goa.design/clue/log"
)

// bootstrap our own services
func bootstrap(ctx context.Context, httpPortF *string, outDbgF *bool) (outHttpPortF *string) {
	outHttpPortF = httpPortF

	_ = godotenv.Load() // load .env file - if exist
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf(ctx, err, "can't initialize config")
	}
	cfg.ApplicationName = cfg.ApplicationName + " (Browser Extension Service)"
	if err := database.Connect(ctx, cfg); err != nil {
		log.Fatalf(ctx, err, "can't connect to database")
	}

	log.Printf(ctx, "starting: %v", cfg.ApplicationName)

	if cfg.Port != "" {
		outHttpPortF = &cfg.Port
	}
	if outDbgF != nil {
		*outDbgF = *outDbgF || cfg.DebugLog
	}

	return
}
