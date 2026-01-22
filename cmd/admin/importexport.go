package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/util/netutil"
	"github.com/spf13/cobra"
)

var (
	importFile string
	exportFile string
)

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "File to import from (default: stdin)")
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "File to export to (default: stdout)")

	feedCmd.AddCommand(importCmd)
	feedCmd.AddCommand(exportCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import feeds from JSON",
	Run: func(cmd *cobra.Command, args []string) {
		importFeeds()
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export feeds to JSON",
	Run: func(cmd *cobra.Command, args []string) {
		exportFeeds()
	},
}

// ImportFeed struct matches the format in feeds.json
type ImportFeed struct {
	URL               string  `json:"url"`
	Language          *string `json:"language,omitempty"`
	RootDomain        *string `json:"root_domain,omitempty"`
	Enabled           *bool   `json:"enabled,omitempty"`
	EnforceFeedDomain *bool   `json:"enforce_feed_domain,omitempty"`
	Polling           *bool   `json:"polling,omitempty"`
}

func importFeeds() {
	var r io.Reader
	if importFile != "" {
		importFile = filepath.Clean(importFile)
		f, err := os.Open(importFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", importFile, err)
			}
		}()
		r = f
	} else {
		r = os.Stdin
	}

	var feeds []ImportFeed
	if err := json.NewDecoder(r).Decode(&feeds); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode JSON: %v\n", err)
		os.Exit(1)
	}

	for _, f := range feeds {
		u, err := parseAndNormalizeURL(f.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping invalid URL %s: %v\n", f.URL, err)
			continue
		}

		// Check if feed already exists
		existing, err := repo.FindFeedByUrlAndAvailability(u, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check feed %s: %v\n", u.String(), err)
			continue
		}

		if existing != nil {
			fmt.Printf("Feed already exists: %s\n", u.String())
			continue
		}

		// Calculate RootDomain if not provided
		var rootDomain *string
		if f.RootDomain != nil {
			rootDomain = f.RootDomain
		} else {
			d := netutil.GetRootDomain(u)
			rootDomain = &d
		}

		enabled := DefaultFeedEnabled
		if f.Enabled != nil {
			enabled = *f.Enabled
		}

		enforce := DefaultFeedEnforceDomain
		if f.EnforceFeedDomain != nil {
			enforce = *f.EnforceFeedDomain
		}

		polling := DefaultFeedPolling
		if f.Polling != nil {
			polling = *f.Polling
		}

		newFeed := &database.Feed{
			URL:               u.String(),
			RootDomain:        rootDomain,
			Language:          f.Language,
			Enabled:           enabled,
			EnforceFeedDomain: enforce,
			Polling:           polling,
		}

		if err := repo.UpsertFeed(newFeed); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create feed %s: %v\n", u.String(), err)
			continue
		}
		fmt.Printf("Imported feed: %s\n", u.String())
	}
}

func exportFeeds() {
	feeds, err := repo.GetAllFeeds(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	var exportFeeds []ImportFeed
	for _, f := range feeds {
		var enabled *bool
		if !f.Enabled {
			v := f.Enabled
			enabled = &v
		}

		var enforce *bool
		if !f.EnforceFeedDomain {
			v := f.EnforceFeedDomain
			enforce = &v
		}

		var polling *bool
		if f.Polling {
			v := f.Polling
			polling = &v
		}

		var exportRootDomain *string
		u, err := parseAndNormalizeURL(f.URL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Skipping invalid URL in export %s: %v\n", f.URL, err)
			continue
		}
		defaultRootDomain := netutil.GetRootDomain(u)
		if f.RootDomain != nil {
			if *f.RootDomain != defaultRootDomain {
				exportRootDomain = f.RootDomain
			}
		}

		exportFeeds = append(exportFeeds, ImportFeed{
			URL:               f.URL,
			Language:          f.Language,
			RootDomain:        exportRootDomain,
			Enabled:           enabled,
			EnforceFeedDomain: enforce,
			Polling:           polling,
		})
	}

	var w io.Writer
	if exportFile != "" {
		exportFile = filepath.Clean(exportFile)
		f, err := os.Create(exportFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", exportFile, err)
			}
		}()
		w = f
	} else {
		w = os.Stdout
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(exportFeeds); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode JSON: %v\n", err)
		os.Exit(1)
	}
}
