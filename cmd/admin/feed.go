package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/valkey"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	jsonOutput  bool
	showDeleted bool
	feedEnabled bool
	repo        database.Repository
)

func init() {
	feedCmd.AddCommand(addCmd)
	feedCmd.AddCommand(deleteCmd)
	feedCmd.AddCommand(enableCmd)
	feedCmd.AddCommand(disableCmd)
	feedCmd.AddCommand(listCmd)

	addCmd.Flags().BoolVar(&feedEnabled, "enabled", true, "Enable the feed")
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	listCmd.Flags().BoolVar(&showDeleted, "deleted", false, "Show deleted feeds")

	rootCmd.AddCommand(feedCmd)
}

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Manage feeds",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if rootCmd.PersistentPreRun != nil {
			rootCmd.PersistentPreRun(cmd, args)
		}
		var err error
		repo, err = database.NewRepository(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a new feed URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addFeed(args[0], feedEnabled)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <uuid|url>",
	Short: "Delete a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteFeed(args[0])
	},
}

var enableCmd = &cobra.Command{
	Use:   "enable <uuid|url>",
	Short: "Enable a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enableFeed(args[0])
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable <uuid|url>",
	Short: "Disable a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		disableFeed(args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all feeds",
	Run: func(cmd *cobra.Command, args []string) {
		listFeeds(jsonOutput, showDeleted)
	},
}

func listFeeds(asJson bool, showDeleted bool) {
	feeds, err := repo.GetAllFeeds(showDeleted)
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
	if _, err := fmt.Fprintln(w, "Status\tID\tURL\tEnforceDomain\tUpdated"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
		os.Exit(1)
	}
	for _, f := range feeds {
		status := "disabled"
		if f.DeletedAt.Valid {
			status = "deleted"
		} else if f.Enabled {
			status = "enabled"
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", status, f.ID, f.URL, f.EnforceFeedDomain, f.UpdatedAt.Format("2006-01-02")); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush to stdout: %v\n", err)
		os.Exit(1)
	}
}

func addFeed(feedUrl string, enabled bool) {
	u, err := parseAndNormalizeURL(feedUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
		os.Exit(1)
	}

	feed, err := repo.FindFeedByUrl(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find feed: %v\n", err)
		os.Exit(1)
	}

	if feed != nil {
		fmt.Fprintf(os.Stderr, "Feed already exists: %s\n", feed.ID)
		os.Exit(1)
	}

	newFeed := &database.Feed{
		URL:               u.String(),
		Enabled:           enabled,
		EnforceFeedDomain: true,
	}

	if err := repo.UpsertFeed(newFeed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added feed for url=%s with id=%s enabled=%v\n", feedUrl, newFeed.ID, newFeed.Enabled)
}

func resolveFeed(input string, onlyEnabled bool) *database.Feed {
	var feed *database.Feed
	var err error

	// Try to parse as UUID first
	if id, err := uuid.Parse(input); err == nil {
		feed, err = repo.FindFeedById(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to find feed: %v\n", err)
			os.Exit(1)
		}
		if feed == nil {
			fmt.Fprintf(os.Stderr, "Feed does not exist with UUID: %s\n", id)
			os.Exit(1)
		}
		return feed
	}

	// If not UUID, treat as URL
	u, err := parseAndNormalizeURL(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid input (neither UUID nor URL): %v\n", err)
		os.Exit(1)
	}

	feed, err = repo.FindFeedByUrlAndAvailability(u, onlyEnabled)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find feed: %v\n", err)
		os.Exit(1)
	}

	if feed == nil {
		fmt.Fprintf(os.Stderr, "Feed does not exist with URL: %s\n", u)
		os.Exit(1)
	}

	return feed
}

func deleteFeed(input string) {
	feed := resolveFeed(input, true)

	if err := repo.DeleteFeedById(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete feed: %v\n", err)
		os.Exit(1)
	}

	drainFeed(feed.ID)
	fmt.Printf("Deleted feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func enableFeed(input string) {
	feed := resolveFeed(input, false)

	feed.Enabled = true
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to enable feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Enabled feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func disableFeed(input string) {
	feed := resolveFeed(input, true)

	feed.Enabled = false
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to disable feed: %v\n", err)
		os.Exit(1)
	}

	drainFeed(feed.ID)
	fmt.Printf("Disabled feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func parseAndNormalizeURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimSuffix(rawURL, "/")
	return url.ParseRequestURI(rawURL)
}

func drainFeed(feedID uuid.UUID) {
	ctx := context.Background()
	if vk, err := valkey.New(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to valkey: %v\n", err)
	} else {
		defer vk.Close()
		if err := vk.DrainFeed(feedID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to drain feed: %v\n", err)
		}
	}
}
