package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/egandro/news-deframer/pkg/config"
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

	slog.Info("I am the worker")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	slog.Info("Shutting down...")
}
