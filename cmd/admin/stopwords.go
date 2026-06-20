package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	stopWordsFile string
	stopWordsJSON bool
)

func init() {
	stopWordsImportCmd.Flags().StringVarP(&stopWordsFile, "file", "f", "", "File to import from (default: stdin)")
	stopWordsExportCmd.Flags().StringVarP(&stopWordsFile, "file", "f", "", "File to export to (default: stdout)")
	stopWordsListCmd.Flags().BoolVar(&stopWordsJSON, "json", false, "Output as JSON")

	stopWordsCmd.AddCommand(stopWordsImportCmd)
	stopWordsCmd.AddCommand(stopWordsExportCmd)
	stopWordsCmd.AddCommand(stopWordsListCmd)
	stopWordsCmd.AddCommand(stopWordsDeleteCmd)
	stopWordsCmd.AddCommand(stopWordsDeleteFeedCmd)

	rootCmd.AddCommand(stopWordsCmd)
}

var stopWordsCmd = &cobra.Command{
	Use:   "stopwords",
	Short: "Manage stopwords",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if rootCmd.PersistentPreRun != nil {
			rootCmd.PersistentPreRun(cmd, args)
		}
		var err error
		repo, err = database.NewRepository(cmd.Context(), cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
	},
}

var stopWordsImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import stopwords from JSON",
	Run: func(cmd *cobra.Command, args []string) {
		importStopWords()
	},
}

var stopWordsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export stopwords to JSON",
	Run: func(cmd *cobra.Command, args []string) {
		exportStopWords()
	},
}

var stopWordsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stopwords",
	Run: func(cmd *cobra.Command, args []string) {
		listStopWords(stopWordsJSON)
	},
}

var stopWordsDeleteCmd = &cobra.Command{
	Use:   "delete <language>",
	Short: "Delete stopwords by language",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteStopWords(args[0])
	},
}

var stopWordsDeleteFeedCmd = &cobra.Command{
	Use:   "delete-feed <uuid|url>",
	Short: "Delete stopwords by feed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteStopWordsByFeed(args[0])
	},
}

type importStopWordsEntry struct {
	Language   string     `json:"language"`
	RootDomain *string    `json:"root_domain,omitempty"`
	FeedID     *uuid.UUID `json:"feed_id,omitempty"`
	FeedURL    *string    `json:"feed_url,omitempty"`
	NounStems  []string   `json:"noun_stems"`
}

func importStopWords() {
	entries, err := readImportStopWords(stopWordsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read JSON: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		language := strings.ToLower(strings.TrimSpace(entry.Language))
		if language == "" {
			fmt.Fprintf(os.Stderr, "Skipping stopwords entry with empty language\n")
			continue
		}

		var stopWords database.StopWords
		stopWords.Language = language
		stopWords.NounStems = normalizeStopWordStems(entry.NounStems)

		if entry.FeedID != nil {
			feed, err := repo.FindFeedById(*entry.FeedID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to resolve feed by id for stopwords import: %v\n", err)
				os.Exit(1)
			}
			if feed == nil {
				fmt.Fprintf(os.Stderr, "Feed does not exist with UUID: %s\n", entry.FeedID.String())
				os.Exit(1)
			}
			stopWords.FeedID = &feed.ID
		} else if entry.FeedURL != nil && strings.TrimSpace(*entry.FeedURL) != "" {
			feed := resolveFeed(*entry.FeedURL, false)
			stopWords.FeedID = &feed.ID
		}

		if err := repo.UpsertStopWords(&stopWords); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import stopwords for language=%s: %v\n", language, err)
			os.Exit(1)
		}
	}
	fmt.Printf("Imported %d stopword entries\n", len(entries))
}

func exportStopWords() {
	stopWords, err := repo.ListStopWords()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list stopwords: %v\n", err)
		os.Exit(1)
	}

	entries := buildStopWordsEntries(stopWords, false)

	var w io.Writer = os.Stdout
	if stopWordsFile != "" {
		cleanPath := filepath.Clean(stopWordsFile)
		f, err := os.Create(cleanPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create file %s: %v\n", cleanPath, err)
			os.Exit(1)
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", cleanPath, err)
			}
		}()
		w = f
	}

	writeStopWordsJSON(w, entries)
}

func buildStopWordsEntries(stopWords []database.StopWords, includeFeedSpecific bool) []importStopWordsEntry {

	entries := make([]importStopWordsEntry, 0, len(stopWords))
	for i := range stopWords {
		if len(stopWords[i].NounStems) == 0 {
			continue
		}
		if !includeFeedSpecific && stopWords[i].FeedID != nil {
			continue
		}
		entry := importStopWordsEntry{
			Language:  stopWords[i].Language,
			NounStems: []string(stopWords[i].NounStems),
		}
		if stopWords[i].FeedID != nil {
			feed, err := repo.FindFeedById(*stopWords[i].FeedID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to resolve feed for stopwords export: %v\n", err)
				os.Exit(1)
			}
			if feed != nil {
				entry.FeedID = &feed.ID
				entry.FeedURL = &feed.URL
				entry.RootDomain = feed.RootDomain
			}
		}
		entries = append(entries, entry)
	}

	return entries
}

func writeStopWordsJSON(w io.Writer, entries []importStopWordsEntry) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode stopwords to JSON: %v\n", err)
		os.Exit(1)
	}
}

func listStopWords(asJSON bool) {
	stopWords, err := repo.ListStopWords()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list stopwords: %v\n", err)
		os.Exit(1)
	}

	entries := buildStopWordsEntries(stopWords, true)
	if asJSON {
		writeStopWordsJSON(os.Stdout, entries)
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "Language\tRootDomain\tFeedID\tFeedURL\tNounStems"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
		os.Exit(1)
	}
	for i := range entries {
		rootDomain := ""
		if entries[i].RootDomain != nil {
			rootDomain = *entries[i].RootDomain
		}
		feedID := ""
		if entries[i].FeedID != nil {
			feedID = entries[i].FeedID.String()
		}
		feedURL := ""
		if entries[i].FeedURL != nil {
			feedURL = *entries[i].FeedURL
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", entries[i].Language, rootDomain, feedID, feedURL, strings.Join(entries[i].NounStems, ",")); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to stdout: %v\n", err)
			os.Exit(1)
		}
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to flush to stdout: %v\n", err)
		os.Exit(1)
	}
}

func deleteStopWords(language string) {
	language = strings.ToLower(strings.TrimSpace(language))
	if language == "" {
		fmt.Fprintln(os.Stderr, "Language must not be empty")
		os.Exit(1)
	}

	if err := repo.DeleteStopWordsByLanguage(language); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete stopwords: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted stopwords for language=%s\n", language)
}

func deleteStopWordsByFeed(input string) {
	feed := resolveFeed(input, false)

	if err := repo.DeleteStopWordsByFeedID(feed.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete stopwords: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleted stopwords for url=%s with id=%s\n", feed.URL, feed.ID)
}

func readImportStopWords(filePath string) ([]importStopWordsEntry, error) {
	var r io.Reader
	if filePath != "" {
		cleanPath := filepath.Clean(filePath)
		f, err := os.Open(cleanPath)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", cleanPath, err)
			}
		}()
		r = f
	} else {
		r = os.Stdin
	}

	var entries []importStopWordsEntry
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func normalizeStopWordStems(values []string) database.StringArray {
	seen := make(map[string]struct{}, len(values))
	result := make(database.StringArray, 0, len(values))
	for i := range values {
		value := strings.TrimSpace(values[i])
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
