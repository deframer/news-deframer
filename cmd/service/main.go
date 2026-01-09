package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/facade"
	"github.com/egandro/news-deframer/pkg/server"
	"github.com/egandro/news-deframer/pkg/valkey"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		lvl = slog.LevelInfo
	}
	fmt.Printf("Log level: %v\n", lvl)
	hostname, _ := os.Hostname()
	fmt.Printf("Hostname: %v\n", hostname)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo, err := database.NewRepository(cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	valkeyClient, err := valkey.New(ctx, cfg)
	if err != nil {
		slog.Error("Failed to connect to valkey", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := valkeyClient.Close(); err != nil {
			slog.Error("Failed to close valkey client", "error", err)
		}
	}()

	f := facade.New(ctx, cfg, valkeyClient, repo)

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
