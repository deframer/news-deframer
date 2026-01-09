package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/spf13/cobra"
)

var (
	targetUrl  string
	jsonOutput bool
)

func init() {
	feedCmd.AddCommand(addCmd)
	feedCmd.AddCommand(deleteCmd)
	feedCmd.AddCommand(enableCmd)
	feedCmd.AddCommand(listCmd)

	deleteCmd.Flags().StringVar(&targetUrl, "url", "", "Delete a feed by URL")
	enableCmd.Flags().StringVar(&targetUrl, "url", "", "Enable a feed by URL")
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

	rootCmd.AddCommand(feedCmd)
}

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage feeds",
}

var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a new feed URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addFeed(args[0])
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [uuid]",
	Short: "Delete a feed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if targetUrl != "" {
			deleteFeedByUrl(targetUrl)
			return
		}
		if len(args) > 0 {
			deleteFeedByUUID(args[0])
			return
		}
		_ = cmd.Help()
		os.Exit(1)
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable [uuid]",
	Short: "Enable a feed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if targetUrl != "" {
			enableFeedByUrl(targetUrl)
			return
		}
		if len(args) > 0 {
			enableFeedByUUID(args[0])
			return
		}
		_ = cmd.Help()
		os.Exit(1)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all feeds",
	Run: func(cmd *cobra.Command, args []string) {
		listFeeds(jsonOutput)
	},
}

func listFeeds(asJson bool) {
	repo, err := database.NewRepository(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	feeds, err := repo.GetAllFeeds(true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	if asJson {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(feeds); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode feeds to JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Status\tID\tURL\tEnforceDomain\tUpdated")
	for _, f := range feeds {
		status := "disabled"
		if f.DeletedAt.Valid {
			status = "deleted"
		} else if f.Enabled {
			status = "enabled"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", status, f.ID, f.URL, f.EnforceFeedDomain, f.UpdatedAt.Format("2006-01-02"))
	}
	w.Flush()
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
