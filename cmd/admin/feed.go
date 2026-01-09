package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	targetUrl   string
	jsonOutput  bool
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
	deleteCmd.Flags().StringVar(&targetUrl, "url", "", "Delete a feed by URL")
	enableCmd.Flags().StringVar(&targetUrl, "url", "", "Enable a feed by URL")
	disableCmd.Flags().StringVar(&targetUrl, "url", "", "Disable a feed by URL")
	listCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")

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
	Use:   "delete [uuid]",
	Short: "Delete a feed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if targetUrl != "" {
			deleteFeed(targetUrl, true)
			return
		}
		if len(args) > 0 {
			deleteFeed(args[0], false)
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
			enableFeed(targetUrl, true)
			return
		}
		if len(args) > 0 {
			enableFeed(args[0], false)
			return
		}
		_ = cmd.Help()
		os.Exit(1)
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable [uuid]",
	Short: "Disable a feed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if targetUrl != "" {
			disableFeed(targetUrl, true)
			return
		}
		if len(args) > 0 {
			disableFeed(args[0], false)
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

	fmt.Printf("Added feed for url=%s with id=%s\n", feedUrl, newFeed.ID)
}

func resolveFeed(input string, isUrl bool, onlyEnabled bool) *database.Feed {
	var feed *database.Feed
	var err error

	if isUrl {
		var u *url.URL
		u, err = parseAndNormalizeURL(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
			os.Exit(1)
		}
		feed, err = repo.FindFeedByUrlAndAvailability(u, onlyEnabled)
	} else {
		var id uuid.UUID
		id, err = uuid.Parse(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid UUID: %v\n", err)
			os.Exit(1)
		}
		feed, err = repo.FindFeedById(id)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to find feed: %v\n", err)
		os.Exit(1)
	}

	if feed == nil {
		fmt.Fprintf(os.Stderr, "Feed does not exist: %s\n", input)
		os.Exit(1)
	}

	return feed
}

func deleteFeed(input string, isUrl bool) {
	feed := resolveFeed(input, isUrl, true)

	if err := repo.DeleteFeedById(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete feed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func enableFeed(input string, isUrl bool) {
	feed := resolveFeed(input, isUrl, false)

	feed.Enabled = true
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to enable feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Enabled feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func disableFeed(input string, isUrl bool) {
	feed := resolveFeed(input, isUrl, true)

	feed.Enabled = false
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to disable feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Disabled feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func parseAndNormalizeURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimSuffix(rawURL, "/")
	return url.ParseRequestURI(rawURL)
}
