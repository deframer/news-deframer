package main

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	feedAdd       string
	feedDeleteUrl string
	feedDelete    string
	feedEnable    string
	feedEnableUrl string
	feedList      bool
)

func init() {
	feedCmd.Flags().StringVar(&feedAdd, "add", "", "Add a new feed URL")
	feedCmd.Flags().StringVar(&feedDeleteUrl, "delete-url", "", "Delete a feed by URL")
	feedCmd.Flags().StringVar(&feedDelete, "delete", "", "Delete a feed by UUID")
	feedCmd.Flags().StringVar(&feedEnable, "enable", "", "Enable a feed by UUID")
	feedCmd.Flags().StringVar(&feedEnableUrl, "enable-url", "", "Enable a feed by URL")
	feedCmd.Flags().BoolVar(&feedList, "list", false, "List all feeds")

	rootCmd.AddCommand(feedCmd)
}

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage feeds",
	Run: func(cmd *cobra.Command, args []string) {
		if feedList {
			slog.Info("Listing feeds")
			return
		}

		if feedAdd != "" {
			slog.Info("Adding feed", "url", feedAdd)
			return
		}

		if feedDeleteUrl != "" {
			slog.Info("Deleting feed", "url", feedDeleteUrl)
			return
		}

		if feedDelete != "" {
			slog.Info("Deleting feed", "uuid", feedDelete)
			return
		}

		if feedEnable != "" {
			slog.Info("Enabling feed", "uuid", feedEnable)
			return
		}

		if feedEnableUrl != "" {
			slog.Info("Enabling feed", "url", feedEnableUrl)
			return
		}

		cmd.Help()
	},
}
