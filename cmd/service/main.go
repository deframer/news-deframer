package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	"github.com/deframer/news-deframer/pkg/server"
)

func main() {
	jsonLog := flag.Bool("json-log", false, "Enable JSON logging")
	disableETag := flag.Bool("disable-etag", false, "Disable ETag caching")
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

	if *disableETag {
		cfg.DisableETag = true
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
	slog.Info("Service", "level", lvl, "hostname", hostname, "etag_disabled", cfg.DisableETag)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := database.NewRepository(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	f := facade.New(ctx, cfg, repo)

	srv := server.New(ctx, cfg, f)

	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	slog.Info("Shutting down...")

	// Shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
	slog.Info("Server stopped")
}
