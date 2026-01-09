package main

import (
	"log/slog"
	"os"

	"github.com/egandro/news-deframer/pkg/config"
	mylogger "github.com/egandro/news-deframer/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	jsonLog bool
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin tool for news-deframer",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg, err = config.Load()
		if err != nil {
			slog.Error("Failed to load config", "error", err)
			os.Exit(1)
		}

		var lvl slog.Level
		if err := lvl.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
			lvl = slog.LevelInfo
		}

		if jsonLog {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
		} else {
			slog.SetDefault(mylogger.NewCommandLogger(lvl))
		}

		slog.Info("I am the admin tool", "level", lvl)
	},
}

func main() {
	rootCmd.PersistentFlags().BoolVar(&jsonLog, "json-log", false, "Enable JSON logging")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
