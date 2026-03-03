package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/syncer"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	jsonLog := flag.Bool("json-log", false, "Enable JSON logging")
	mode := flag.String("mode", string(syncer.ModeWorker), "Run mode: worker or think-fixer")
	about := flag.Bool("about", false, "Show build and mode information")
	flag.Usage = func() {
		// #nosec G705: usage string is escaped before printing
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", html.EscapeString(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if *about {
		fmt.Printf("news-deframer worker\n")
		fmt.Printf("version: %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("date: %s\n", date)
		fmt.Printf("modes: %s, %s\n", syncer.ModeWorker, syncer.ModeThinkFixer)
		return
	}

	selectedMode := syncer.Mode(*mode)
	switch selectedMode {
	case syncer.ModeWorker, syncer.ModeThinkFixer:
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s (expected %s or %s)\n", *mode, syncer.ModeWorker, syncer.ModeThinkFixer)
		os.Exit(2)
	}

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
	s.Poll(selectedMode)
	slog.Info("Shutting down...")
}
