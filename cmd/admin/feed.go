package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/syncer"
	"github.com/deframer/news-deframer/pkg/util/netutil"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	jsonOutput     bool
	showDeleted    bool
	feedEnabled    bool
	polling        bool
	mining         bool
	resolveItemUrl bool
	language       string
	categories     []string
	noRootDomain   bool
	purgeFeed      bool
	repo           database.Repository
	feedSyncer     *syncer.Syncer
)

func init() {
	feedCmd.AddCommand(addCmd)
	feedCmd.AddCommand(deleteCmd)
	feedCmd.AddCommand(enableCmd)
	feedCmd.AddCommand(disableCmd)
	feedCmd.AddCommand(listCmd)
	feedCmd.AddCommand(pollingCmd)
	feedCmd.AddCommand(miningCmd)
	feedCmd.AddCommand(resolveItemUrlCmd)
	feedCmd.AddCommand(syncCmd)
	feedCmd.AddCommand(syncAllCmd)
	feedCmd.AddCommand(mineCmd)
	feedCmd.AddCommand(mineAllCmd)

	languageCmd.AddCommand(setLanguageCmd)
	languageCmd.AddCommand(deleteLanguageCmd)
	feedCmd.AddCommand(languageCmd)

	categoriesCmd.AddCommand(setCategoriesCmd)
	categoriesCmd.AddCommand(deleteCategoriesCmd)
	feedCmd.AddCommand(categoriesCmd)

	addCmd.Flags().BoolVar(&feedEnabled, "enabled", DefaultFeedEnabled, "Enable the feed")
	addCmd.Flags().BoolVar(&polling, "polling", DefaultFeedPolling, "Enable polling")
	addCmd.Flags().BoolVar(&mining, "mining", DefaultFeedMining, "Enable mining")
	addCmd.Flags().BoolVar(&resolveItemUrl, "resolve-item-url", DefaultResolveItemUrl, "Enable item URL resolution")
	addCmd.Flags().StringVar(&language, "language", "", "Set a two-letter ISO 639-1 language code for the feed")
	addCmd.Flags().StringSliceVar(&categories, "categories", []string{}, "Set a comma-separated list of categories for the feed")
	addCmd.Flags().BoolVar(&noRootDomain, "no-root-domain", false, "Do not automatically populate root_domain")
	deleteCmd.Flags().BoolVar(&purgeFeed, "purge", false, "Purge the feed and all related data")
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
		feedSyncer, err = syncer.New(context.Background(), cfg, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create syncer: %v\n", err)
			os.Exit(1)
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Add a new feed URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addFeed(args[0], feedEnabled, polling, mining, noRootDomain, language, resolveItemUrl, categories)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete <uuid|url>",
	Short: "Delete a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteFeed(args[0], purgeFeed)
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

var pollingCmd = &cobra.Command{
	Use:   "polling <uuid|url> <true|false>",
	Short: "Set polling for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setPolling(args[0], args[1])
	},
}

var miningCmd = &cobra.Command{
	Use:   "mining <uuid|url> <true|false>",
	Short: "Set mining for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setMining(args[0], args[1])
	},
}

var resolveItemUrlCmd = &cobra.Command{
	Use:   "resolve-item-url <uuid|url> <true|false>",
	Short: "Set item URL resolution for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setResolveItemUrl(args[0], args[1])
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync <uuid|url>",
	Short: "Sync a feed immediately",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		syncFeed(args[0])
	},
}

var mineCmd = &cobra.Command{
	Use:   "mine <uuid|url>",
	Short: "Mine a feed immediately",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mineFeed(args[0])
	},
}

var syncAllCmd = &cobra.Command{
	Use:   "sync-all",
	Short: "Sync all active enabled feeds immediately",
	Run: func(cmd *cobra.Command, args []string) {
		syncAllFeeds()
	},
}

var mineAllCmd = &cobra.Command{
	Use:   "mine-all",
	Short: "Mine all active enabled feeds immediately",
	Run: func(cmd *cobra.Command, args []string) {
		mineAllFeeds()
	},
}

var languageCmd = &cobra.Command{
	Use:   "language",
	Short: "Manage feed language",
}

var setLanguageCmd = &cobra.Command{
	Use:   "set <uuid|url> <language-code>",
	Short: "Set the language for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setLanguage(args[0], args[1])
	},
}

var deleteLanguageCmd = &cobra.Command{
	Use:   "delete <uuid|url>",
	Short: "Delete the language for a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteLanguage(args[0])
	},
}

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "Manage feed categories",
}

var setCategoriesCmd = &cobra.Command{
	Use:   "set <uuid|url> <category1,category2,...>",
	Short: "Set the categories for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setCategories(args[0], args[1])
	},
}

var deleteCategoriesCmd = &cobra.Command{
	Use:   "delete <uuid|url>",
	Short: "Delete the categories for a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteCategories(args[0])
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
	if _, err := fmt.Fprintln(w, "Status\tPolling\tMining\tResolveItemUrl\tLanguage\tCategories\tID\tURL\tRootDomain\tEnforceDomain\tUpdated\tSync Status"); err != nil {
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

		state := "-"
		if f.FeedSchedule != nil {
			if f.FeedSchedule.LastError != nil {
				state = fmt.Sprintf("Error: %s", *f.FeedSchedule.LastError)
			} else if f.FeedSchedule.LockedUntil != nil {
				state = "Running"
			} else if f.FeedSchedule.NextRunAt != nil {
				state = "Next: " + f.FeedSchedule.NextRunAt.Format("2006-01-02 15:04")
			}
		}

		rootDomain := "-"
		if f.RootDomain != nil {
			rootDomain = *f.RootDomain
		}

		language := "-"
		if f.Language != nil {
			language = *f.Language
		}

		categories := "-"
		if len(f.Categories) > 0 {
			categories = strings.Join(f.Categories, ",")
		}

		if _, err := fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%s\t%s\t%s\t%s\t%s\t%v\t%s\t%s\n", status, f.Polling, f.Mining, f.ResolveItemUrl, language, categories, f.ID, f.URL, rootDomain, f.EnforceFeedDomain, f.UpdatedAt.Format("2006-01-02"), state); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush to stdout: %v\n", err)
		os.Exit(1)
	}
}

func addFeed(feedUrl string, enabled bool, polling bool, mining bool, noRootDomain bool, language string, resolveItemUrl bool, categories []string) {
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

	var rootDomain *string
	if !noRootDomain {
		d := netutil.GetRootDomain(u)
		rootDomain = &d
	}

	var languagePtr *string
	if language != "" {
		languagePtr = &language
	}

	newFeed := &database.Feed{
		URL:               u.String(),
		RootDomain:        rootDomain,
		Enabled:           enabled,
		EnforceFeedDomain: DefaultFeedEnforceDomain,
		Polling:           polling,
		Mining:            mining,
		Language:          languagePtr,
		ResolveItemUrl:    resolveItemUrl,
		Categories:        categories,
	}

	if err := repo.UpsertFeed(newFeed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create feed: %v\n", err)
		os.Exit(1)
	}

	if polling && enabled {
		if err := feedSyncer.SyncFeed(newFeed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to trigger initial sync: %v\n", err)
		}
	}

	rootDomainStr := "<nil>"
	if newFeed.RootDomain != nil {
		rootDomainStr = *newFeed.RootDomain
	}

	output := fmt.Sprintf("Added feed for url=%s with id=%s enabled=%v polling=%v mining=%v root_domain=%s", feedUrl, newFeed.ID, newFeed.Enabled, newFeed.Polling, newFeed.Mining, rootDomainStr)
	if newFeed.Language != nil {
		output += fmt.Sprintf(" language=%s", *newFeed.Language)
	}
	if newFeed.ResolveItemUrl {
		output += fmt.Sprintf(" resolve_item_url=%v", newFeed.ResolveItemUrl)
	}
	if len(newFeed.Categories) > 0 {
		output += fmt.Sprintf(" categories=%s", strings.Join(newFeed.Categories, ","))
	}
	fmt.Println(output)
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

func deleteFeed(input string, purge bool) {
	feed := resolveFeed(input, false)

	if err := feedSyncer.StopPolling(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to stop polling: %v\n", err)
	}

	if purge {
		if err := repo.PurgeFeedById(feed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to purge feed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Purged feed for url=%s with id=%s\n", feed.URL, feed.ID)
	} else {
		if err := repo.DeleteFeedById(feed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete feed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted feed for url=%s with id=%s\n", feed.URL, feed.ID)
	}
}

func enableFeed(input string) {
	feed := resolveFeed(input, false)

	feed.Enabled = true
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to enable feed: %v\n", err)
		os.Exit(1)
	}

	if feed.Polling {
		if err := feedSyncer.SyncFeed(feed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to sync feed: %v\n", err)
		}
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

	if err := feedSyncer.StopPolling(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to stop polling: %v\n", err)
	}

	fmt.Printf("Disabled feed for url=%s with id=%s\n", feed.URL, feed.ID)
}

func setPolling(input string, stateStr string) {
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid boolean value: %s\n", stateStr)
		os.Exit(1)
	}

	feed := resolveFeed(input, false)

	feed.Polling = state
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set polling: %v\n", err)
		os.Exit(1)
	}

	if !state {
		if err := feedSyncer.StopPolling(feed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to stop polling: %v\n", err)
		}
	} else {
		if err := feedSyncer.SyncFeed(feed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to sync feed: %v\n", err)
		}
	}

	fmt.Printf("Set polling to %v for url=%s with id=%s\n", feed.Polling, feed.URL, feed.ID)
}

func setMining(input string, stateStr string) {
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid boolean value: %s\n", stateStr)
		os.Exit(1)
	}

	feed := resolveFeed(input, false)

	feed.Mining = state
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set mining: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Set mining to %v for url=%s with id=%s\n", feed.Mining, feed.URL, feed.ID)
}

func setResolveItemUrl(input string, stateStr string) {
	state, err := strconv.ParseBool(stateStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid boolean value: %s\n", stateStr)
		os.Exit(1)
	}

	feed := resolveFeed(input, false)

	feed.ResolveItemUrl = state
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set resolve_item_url: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Set resolve_item_url to %v for url=%s with id=%s\n", feed.ResolveItemUrl, feed.URL, feed.ID)
}

func syncFeed(input string) {
	feed := resolveFeed(input, false)

	if err := feedSyncer.SyncFeed(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sync feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Triggered sync for url=%s with id=%s\n", feed.URL, feed.ID)
}

func mineFeed(input string) {
	feed := resolveFeed(input, false)

	if err := repo.EnqueueMine(feed.ID, 0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to enqueue mining for feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Triggered mining for url=%s with id=%s\n", feed.URL, feed.ID)
}

func syncAllFeeds() {
	feeds, err := repo.GetAllFeeds(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, f := range feeds {
		if !f.Enabled {
			continue
		}
		if err := feedSyncer.SyncFeed(f.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to sync feed %s (%s): %v\n", f.URL, f.ID, err)
		} else {
			fmt.Printf("Triggered sync for url=%s with id=%s\n", f.URL, f.ID)
			count++
		}
	}
	fmt.Printf("Triggered sync for %d feeds\n", count)
}

func mineAllFeeds() {
	feeds, err := repo.GetAllFeeds(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, f := range feeds {
		if !f.Enabled || !f.Mining {
			continue
		}
		if err := repo.EnqueueMine(f.ID, 0, 0); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to enqueue mining for feed %s (%s): %v\n", f.URL, f.ID, err)
		} else {
			fmt.Printf("Triggered mining for url=%s with id=%s\n", f.URL, f.ID)
			count++
		}
	}
	fmt.Printf("Triggered mining for %d feeds\n", count)
}

func setLanguage(input string, language string) {
	feed := resolveFeed(input, false)

	feed.Language = &language
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set language: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Set language to %s for url=%s with id=%s\n", language, feed.URL, feed.ID)
}

func deleteLanguage(input string) {
	feed := resolveFeed(input, false)

	feed.Language = nil
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete language: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted language for url=%s with id=%s\n", feed.URL, feed.ID)
}

func parseAndNormalizeURL(rawURL string) (*url.URL, error) {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.TrimSuffix(rawURL, "/")
	return url.ParseRequestURI(rawURL)
}

func setCategories(input string, categoriesStr string) {
	feed := resolveFeed(input, false)

	categories := strings.Split(categoriesStr, ",")
	feed.Categories = categories
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set categories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Set categories to %v for url=%s with id=%s\n", categories, feed.URL, feed.ID)
}

func deleteCategories(input string) {
	feed := resolveFeed(input, false)

	feed.Categories = []string{}
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete categories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted categories for url=%s with id=%s\n", feed.URL, feed.ID)
}
