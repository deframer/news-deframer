package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/syncer"
)

func main() {
	jsonLog := flag.Bool("json-log", false, "Enable JSON logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		lvl = slog.LevelInfo
	}

	useJSON := *jsonLog
	if !useJSON {
		// Check if stdout is a terminal (interactive)
		if fi, err := os.Stdout.Stat(); err == nil {
			if (fi.Mode() & os.ModeCharDevice) == 0 {
				useJSON = true
			}
		}
	}

	if useJSON {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
	} else {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
	}

	hostname, _ := os.Hostname()
	slog.Info("Worker", "level", lvl, "hostname", hostname)

	repo, err := database.NewRepository(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	s, err := syncer.New(ctx, cfg, repo)
	if err != nil {
		slog.Error("Failed to create syncer", "error", err)
		os.Exit(1)
	}

	// start syncer poll
	s.Poll()
	slog.Info("Shutting down...")
}
