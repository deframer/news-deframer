package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/egandro/news-deframer/pkg/config"
	mylogger "github.com/egandro/news-deframer/pkg/logger"
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
	if useJSON {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
	} else {
		slog.SetDefault(mylogger.NewCommandLogger(lvl))
	}

	slog.Info("I am the admin tool", "level", lvl)
}
