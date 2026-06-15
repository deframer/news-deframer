package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"os"
	"os/signal"
	"syscall"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	applog "github.com/deframer/news-deframer/pkg/logger"
	"github.com/deframer/news-deframer/pkg/syncer"
	"goa.design/clue/log"
)

func main() {
	mode := flag.String("mode", string(syncer.ModeIngester), "Run mode: ingester, thinker, thinker-fixer, or thinker-update-llm-model")
	flag.Usage = func() {
		// #nosec G705: usage string is escaped before printing
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", html.EscapeString(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	selectedMode := syncer.Mode(*mode)
	switch selectedMode {
	case syncer.ModeIngester, syncer.ModeThinker, syncer.ModeThinkerFixer, syncer.ModeThinkerUpdateLLMModel:
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s (expected %s, %s, %s, or %s)\n", *mode, syncer.ModeIngester, syncer.ModeThinker, syncer.ModeThinkerFixer, syncer.ModeThinkerUpdateLLMModel)
		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(context.Background(), err, "Failed to load config")
		os.Exit(1)
	}
	cfg.ApplicationName = fmt.Sprintf("%s (%s)", cfg.ApplicationName, selectedMode)

	logCtx := applog.NewLoggerContext(context.Background(), cfg.DebugLog)

	hostname, _ := os.Hostname()
	log.Print(logCtx, log.KV{K: "component", V: "worker"}, log.KV{K: "hostname", V: hostname})

	repo, err := database.NewRepository(logCtx, cfg)
	if err != nil {
		log.Fatalf(logCtx, err, "Failed to connect to database")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(logCtx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s, err := syncer.New(ctx, cfg, repo)
	if err != nil {
		log.Fatalf(logCtx, err, "Failed to create syncer")
		os.Exit(1)
	}

	// start syncer poll
	s.Poll(selectedMode)
	log.Print(logCtx, log.KV{K: "message", V: "Shutting down..."})
}
