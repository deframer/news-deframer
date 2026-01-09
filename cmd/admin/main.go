package main

import (
	"fmt"
	"log/slog"
	"os"

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
	slog.Info("I am the admin tool")
}
