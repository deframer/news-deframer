package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/util/netutil"
	"github.com/spf13/cobra"
)

var (
	importFile    string
	exportFile    string
	validateFile  string
	validateFetch bool
)

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "File to import from (default: stdin)")
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "File to export to (default: stdout)")
	validateCmd.Flags().StringVarP(&validateFile, "file", "f", "", "File to validate from (default: stdin)")
	validateCmd.Flags().BoolVar(&validateFetch, "fetch", false, "Fetch URLs and verify feed XML")

	feedCmd.AddCommand(importCmd)
	feedCmd.AddCommand(exportCmd)
	feedCmd.AddCommand(validateCmd)
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

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate feeds from JSON",
	Run: func(cmd *cobra.Command, args []string) {
		validateFeeds()
	},
}

// ImportFeed struct matches the format in feeds.json
type ImportFeed struct {
	URL               string   `json:"url"`
	Language          *string  `json:"language,omitempty"`
	Country           *string  `json:"country,omitempty"`
	Categories        []string `json:"categories,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	RootDomain        *string  `json:"root_domain,omitempty"`
	PortalUrl         *string  `json:"portal_url,omitempty"`
	Enabled           *bool    `json:"enabled,omitempty"`
	EnforceFeedDomain *bool    `json:"enforce_feed_domain,omitempty"`
	Polling           *bool    `json:"polling,omitempty"`
	Mining            *bool    `json:"mining,omitempty"`
	ResolveItemUrl    *bool    `json:"resolve_item_url,omitempty"`
}

func importFeeds() {
	feeds, err := readImportFeeds(importFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read JSON: %v\n", err)
		os.Exit(1)
	}

	for _, f := range feeds {
		f.Categories = validateCategoryList(f.Categories)
		if f.Tags == nil {
			f.Tags = []string{}
		}

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
			fmt.Printf("Updating existing feed: %s\n", u.String())
			updateExistingFeed(existing, f)
			if err := repo.UpsertFeed(existing); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update feed %s: %v\n", u.String(), err)
			}
			if err := repo.CreateFeedSchedule(existing.ID); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to ensure schedule for feed %s: %v\n", u.String(), err)
			}

			// Handle Polling schedule
			if existing.Enabled && existing.Polling {
				if err := repo.EnqueueSync(existing.ID, 0); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to enqueue sync for feed %s: %v\n", u.String(), err)
				}
			} else {
				if err := repo.RemoveSync(existing.ID); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove sync for feed %s: %v\n", u.String(), err)
				}
			}

			// Handle Mining schedule
			if existing.Enabled && existing.Mining {
				if err := repo.EnqueueMine(existing.ID, 0); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to enqueue mine for feed %s: %v\n", u.String(), err)
				}
			} else {
				if err := repo.RemoveMine(existing.ID); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove mine for feed %s: %v\n", u.String(), err)
				}
			}
			continue
		}

		fmt.Printf("Importing new feed: %s\n", u.String())
		newFeed := createNewFeed(u, f)
		if err := repo.UpsertFeed(newFeed); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create feed %s: %v\n", u.String(), err)
			continue
		}
		if err := repo.CreateFeedSchedule(newFeed.ID); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create schedule for feed %s: %v\n", u.String(), err)
		}
		if newFeed.Enabled {
			if newFeed.Polling {
				if err := repo.EnqueueSync(newFeed.ID, 0); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to enqueue sync for feed %s: %v\n", u.String(), err)
				}
			}
			if newFeed.Mining {
				if err := repo.EnqueueMine(newFeed.ID, 0); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to enqueue mine for feed %s: %v\n", u.String(), err)
				}
			}
		}
	}
}

func validateFeeds() {
	validateFeedsWithContext(context.Background())
}

func validateFeedsWithContext(ctx context.Context) {
	feeds, err := readImportFeeds(validateFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read JSON: %v\n", err)
		os.Exit(1)
	}

	if hadErrors := validateImportFeeds(ctx, feeds, os.Stdout, validateFetch); hadErrors {
		os.Exit(1)
	}
}

func readImportFeeds(filePath string) ([]ImportFeed, error) {
	var r io.Reader
	if filePath != "" {
		filePath = filepath.Clean(filePath)
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", filePath, err)
			}
		}()
		r = f
	} else {
		r = os.Stdin
	}

	var feeds []ImportFeed
	if err := json.NewDecoder(r).Decode(&feeds); err != nil {
		return nil, err
	}
	return feeds, nil
}

func validateImportFeeds(ctx context.Context, feeds []ImportFeed, w io.Writer, fetch bool) bool {
	hadErrors := false
	errorRows := 0
	seen := make(map[string]int, len(feeds))
	for idx, f := range feeds {
		line := idx + 1
		rawURL := strings.TrimSpace(f.URL)
		if rawURL == "" {
			_, _ = fmt.Fprintf(w, "%d ERROR missing url\n", line)
			hadErrors = true
			errorRows++
			continue
		}

		u, err := parseAndNormalizeURL(rawURL)
		if err != nil {
			_, _ = fmt.Fprintf(w, "%d ERROR invalid url %q: %v\n", line, rawURL, err)
			hadErrors = true
			errorRows++
			continue
		}

		normalizedURL := u.String()
		if firstLine, ok := seen[normalizedURL]; ok {
			_, _ = fmt.Fprintf(w, "%d ERROR duplicate url %s (also in row %d)\n", line, normalizedURL, firstLine)
			hadErrors = true
			errorRows++
			continue
		}
		seen[normalizedURL] = line

		_, categoryErrors := validateCategoryListStrict(f.Categories)
		if len(categoryErrors) > 0 {
			for _, categoryError := range categoryErrors {
				_, _ = fmt.Fprintf(w, "%d ERROR %s\n", line, categoryError)
			}
			hadErrors = true
			errorRows++
			continue
		}

		if f.Language != nil {
			language := strings.TrimSpace(*f.Language)
			if language != "" && len(language) != 2 {
				_, _ = fmt.Fprintf(w, "%d ERROR invalid language %q; expected a two-letter ISO 639-1 code\n", line, language)
				hadErrors = true
				errorRows++
				continue
			}
		}

		existing, err := repo.FindFeedByUrlAndAvailability(u, false)
		if err != nil {
			_, _ = fmt.Fprintf(w, "%d ERROR failed to check existing feed %s: %v\n", line, normalizedURL, err)
			hadErrors = true
			errorRows++
			continue
		}

		status := "would insert"
		if existing != nil {
			status = "would update"
		}

		row := fmt.Sprintf("%d OK %s url=%s", line, status, normalizedURL)
		if f.Language != nil {
			row += " language=" + strings.TrimSpace(*f.Language)
		}
		if f.Country != nil && strings.TrimSpace(*f.Country) != "" {
			row += " country=" + strings.TrimSpace(*f.Country)
		}
		if f.RootDomain != nil {
			row += " root_domain=" + strings.TrimSpace(*f.RootDomain)
		}
		if f.PortalUrl != nil {
			row += " portal_url=" + strings.TrimSpace(*f.PortalUrl)
		}
		if f.Enabled != nil {
			row += fmt.Sprintf(" enabled=%v", *f.Enabled)
		}
		if f.Polling != nil {
			row += fmt.Sprintf(" polling=%v", *f.Polling)
		}
		if f.Mining != nil {
			row += fmt.Sprintf(" mining=%v", *f.Mining)
		}
		if f.ResolveItemUrl != nil {
			row += fmt.Sprintf(" resolve_item_url=%v", *f.ResolveItemUrl)
		}
		if len(f.Categories) > 0 {
			row += " categories=" + strings.Join(f.Categories, ",")
		}
		if len(f.Tags) > 0 {
			row += " tags=" + strings.Join(f.Tags, ",")
		}
		_, _ = fmt.Fprintln(w, row)

		if fetch {
			if err := probeFeedXML(ctx, u); err != nil {
				_, _ = fmt.Fprintf(w, "%d ERROR %s\n", line, err)
				hadErrors = true
				errorRows++
			}
		}
	}

	_, _ = fmt.Fprintf(w, "validated %d feeds (%d errors)\n", len(feeds), errorRows)
	return hadErrors
}

func probeFeedXML(ctx context.Context, u *url.URL) error {
	if feedDownloader == nil {
		return fmt.Errorf("not reachable: downloader not initialized")
	}

	body, err := feedDownloader.DownloadRSSFeed(ctx, u)
	if err != nil {
		return fmt.Errorf("not reachable: %w", err)
	}
	defer func() { _ = body.Close() }()

	decoder := xml.NewDecoder(io.LimitReader(body, 64*1024))
	for {
		tok, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("not rss/atom: %w", err)
		}
		if start, ok := tok.(xml.StartElement); ok {
			root := strings.ToLower(start.Name.Local)
			if root == "rss" || root == "feed" || root == "rdf" {
				return nil
			}
			return fmt.Errorf("not rss/atom: root element %q", start.Name.Local)
		}
	}
}

func createNewFeed(u *url.URL, f ImportFeed) *database.Feed {
	if f.Tags == nil {
		f.Tags = []string{}
	}

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

	mining := DefaultFeedMining
	if f.Mining != nil {
		mining = *f.Mining
	}

	resolveItemUrl := DefaultResolveItemUrl
	if f.ResolveItemUrl != nil {
		resolveItemUrl = *f.ResolveItemUrl
	}

	return &database.Feed{
		URL:        u.String(),
		RootDomain: rootDomain,
		PortalUrl:  f.PortalUrl,
		Language:   normalizeImportLanguage(f.Language),
		Country: func() string {
			if f.Country != nil {
				return *f.Country
			}
			return ""
		}(),
		Categories:        f.Categories,
		Tags:              f.Tags,
		Enabled:           enabled,
		EnforceFeedDomain: enforce,
		Polling:           polling,
		Mining:            mining,
		ResolveItemUrl:    resolveItemUrl,
	}
}

func updateExistingFeed(existing *database.Feed, f ImportFeed) {
	if f.Language != nil {
		existing.Language = normalizeImportLanguage(f.Language)
	}
	if f.Country != nil {
		existing.Country = *f.Country
	}
	existing.Categories = f.Categories
	existing.Tags = f.Tags
	if f.RootDomain != nil {
		existing.RootDomain = f.RootDomain
	}
	if f.PortalUrl != nil {
		existing.PortalUrl = f.PortalUrl
	}
	if f.Enabled != nil {
		existing.Enabled = *f.Enabled
	}
	if f.EnforceFeedDomain != nil {
		existing.EnforceFeedDomain = *f.EnforceFeedDomain
	}
	if f.Polling != nil {
		existing.Polling = *f.Polling
	}
	if f.Mining != nil {
		existing.Mining = *f.Mining
	}
	if f.ResolveItemUrl != nil {
		existing.ResolveItemUrl = *f.ResolveItemUrl
	}
}

func normalizeImportLanguage(language *string) *string {
	if language == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*language)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func exportFeeds() {
	feeds, err := repo.GetAllFeeds(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list feeds: %v\n", err)
		os.Exit(1)
	}

	var exportFeeds []ImportFeed
	for _, f := range feeds {
		tags := f.Tags
		if tags == nil {
			tags = []string{}
		}
		var country *string
		if f.Country != "" {
			country = &f.Country
		}
		exportFeeds = append(exportFeeds, ImportFeed{
			URL:               f.URL,
			Language:          f.Language,
			Country:           country,
			Categories:        f.Categories,
			Tags:              tags,
			RootDomain:        f.RootDomain,
			PortalUrl:         f.PortalUrl,
			Enabled:           &f.Enabled,
			EnforceFeedDomain: &f.EnforceFeedDomain,
			Polling:           &f.Polling,
			Mining:            &f.Mining,
			ResolveItemUrl:    &f.ResolveItemUrl,
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
