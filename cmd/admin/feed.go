package main

import (
	"fmt"

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
			listFeeds()
			return
		}

		if feedAdd != "" {
			addFeed(feedAdd)
			return
		}

		if feedDeleteUrl != "" {
			deleteFeedByUrl(feedDeleteUrl)
			return
		}

		if feedDelete != "" {
			deleteFeedByUUID(feedDelete)
			return
		}

		if feedEnable != "" {
			enableFeedByUUID(feedEnable)
			return
		}

		if feedEnableUrl != "" {
			enableFeedByUrl(feedEnableUrl)
			return
		}

		cmd.Help()
	},
}

func listFeeds() {
	fmt.Println("Listing feeds")
}

func addFeed(url string) {
	fmt.Printf("Adding feed url=%s\n", url)
}

func deleteFeedByUrl(url string) {
	fmt.Printf("Deleting feed url=%s\n", url)
}

func deleteFeedByUUID(uuid string) {
	fmt.Printf("Deleting feed uuid=%s\n", uuid)
}

func enableFeedByUUID(uuid string) {
	fmt.Printf("Enabling feed uuid=%s\n", uuid)
}

func enableFeedByUrl(url string) {
	fmt.Printf("Enabling feed url=%s\n", url)
}
