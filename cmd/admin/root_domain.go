package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rootDomainCmd)

	rootDomainCmd.AddCommand(rootDomainListCmd)
	rootDomainListCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output as JSON")
	rootDomainCmd.AddCommand(rootDomainSetCmd)
	rootDomainCmd.AddCommand(rootDomainClearCmd)
}

var rootDomainCmd = &cobra.Command{
	Use:   "root-domain",
	Short: "Manage root domains",
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

var rootDomainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all unique root domains",
	Run: func(cmd *cobra.Command, args []string) {
		listRootDomains()
	},
}

var rootDomainSetCmd = &cobra.Command{
	Use:   "set <uuid|url> <root_domain>",
	Short: "Set root domain for a feed",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setRootDomain(args[0], args[1])
	},
}

var rootDomainClearCmd = &cobra.Command{
	Use:   "clear <uuid|url>",
	Short: "Clear root domain for a feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clearRootDomain(args[0])
	},
}

func listRootDomains() {
	feeds, err := repo.GetAllFeeds(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	type FeedInfo struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}

	type DomainGroup struct {
		RootDomain string     `json:"root_domain"`
		Count      int        `json:"count"`
		Feeds      []FeedInfo `json:"feeds"`
	}

	domains := make(map[string][]FeedInfo)
	for _, f := range feeds {
		if !f.Enabled {
			continue
		}
		if f.RootDomain != nil && *f.RootDomain != "" {
			domains[*f.RootDomain] = append(domains[*f.RootDomain], FeedInfo{ID: f.ID.String(), URL: f.URL})
		}
	}

	if jsonOutput {
		var output []DomainGroup
		for d, infos := range domains {
			output = append(output, DomainGroup{
				RootDomain: d,
				Count:      len(infos),
				Feeds:      infos,
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(output); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode to JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "RootDomain\tCount"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
		os.Exit(1)
	}
	for d, infos := range domains {
		if _, err := fmt.Fprintf(w, "%s\t%d\n", d, len(infos)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush to stdout: %v\n", err)
		os.Exit(1)
	}
}

func setRootDomain(input string, domain string) {
	feed := resolveFeed(input, false)
	if feed.DeletedAt.Valid {
		fmt.Fprintf(os.Stderr, "Feed is deleted: %s\n", feed.ID)
		os.Exit(1)
	}
	feed.RootDomain = &domain
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Set root_domain=%s for feed %s\n", domain, feed.ID)
}

func clearRootDomain(input string) {
	feed := resolveFeed(input, false)
	if feed.DeletedAt.Valid {
		fmt.Fprintf(os.Stderr, "Feed is deleted: %s\n", feed.ID)
		os.Exit(1)
	}
	feed.RootDomain = nil
	if err := repo.UpsertFeed(feed); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update feed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cleared root_domain for feed %s\n", feed.ID)
}
