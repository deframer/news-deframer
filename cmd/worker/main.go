package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/egandro/news-deframer/pkg/config"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	slog.Info("Shutting down...")
}
