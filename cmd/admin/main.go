package main

import (
	"context"
	"fmt"
	"os"

	"github.com/deframer/news-deframer/pkg/config"
	applog "github.com/deframer/news-deframer/pkg/logger"
	"github.com/spf13/cobra"
	"goa.design/clue/log"
)

var (
	cfg *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin tool for news-deframer",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		cfg, err = config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		ctx := applog.NewLoggerContext(context.Background(), false)
		log.Print(ctx, log.KV{K: "component", V: "admin"})
		cmd.SetContext(ctx)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
