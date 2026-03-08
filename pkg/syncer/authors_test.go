package syncer

import (
	"context"
	"log/slog"
	"testing"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/stretchr/testify/assert"
)

func TestExtractAuthors(t *testing.T) {
	s := &Syncer{}

	t.Run("extracts from dc creator", func(t *testing.T) {
		item := &gofeed.Item{
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "John Doe"}},
				},
			},
		}

		assert.Equal(t, database.StringArray{"John Doe"}, s.extractAndNormalizeAuthors(item, "en"))
	})

	t.Run("splits combined authors like sql", func(t *testing.T) {
		item := &gofeed.Item{
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "Alice &amp; Bob / Carol and Dora"}},
				},
			},
		}

		assert.Equal(t, database.StringArray{"Alice", "Bob", "Carol", "Dora"}, s.extractAndNormalizeAuthors(item, "en"))
	})

	t.Run("splits language specific joiners", func(t *testing.T) {
		testCases := []struct {
			name     string
			language string
			value    string
			expected database.StringArray
		}{
			{name: "german", language: "de", value: "Alice und Bob", expected: database.StringArray{"Alice", "Bob"}},
			{name: "french", language: "fr", value: "Alice et Bob", expected: database.StringArray{"Alice", "Bob"}},
			{name: "spanish", language: "es", value: "Alice y Bob", expected: database.StringArray{"Alice", "Bob"}},
			{name: "portuguese", language: "pt", value: "Alice e Bob", expected: database.StringArray{"Alice", "Bob"}},
			{name: "swedish", language: "sv", value: "Alice och Bob", expected: database.StringArray{"Alice", "Bob"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				item := &gofeed.Item{
					Extensions: ext.Extensions{
						"dc": {
							"creator": []ext.Extension{{Value: tc.value}},
						},
					},
				}

				assert.Equal(t, tc.expected, s.extractAndNormalizeAuthors(item, tc.language))
			})
		}
	})

	t.Run("falls back to default joiners for unknown language", func(t *testing.T) {
		item := &gofeed.Item{
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "Alice and Bob"}},
				},
			},
		}

		assert.Equal(t, database.StringArray{"Alice", "Bob"}, s.extractAndNormalizeAuthors(item, "zz"))
	})

	t.Run("removes publisher suffix", func(t *testing.T) {
		item := &gofeed.Item{
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "Jane Doe | Example News"}},
				},
			},
		}

		assert.Equal(t, database.StringArray{"Jane Doe"}, s.extractAndNormalizeAuthors(item, "en"))
	})

	t.Run("prefers parsed gofeed authors", func(t *testing.T) {
		item := &gofeed.Item{
			Authors: []*gofeed.Person{{Name: "Parsed Author"}},
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "Fallback Author"}},
				},
			},
		}

		assert.Equal(t, database.StringArray{"Parsed Author"}, s.extractAndNormalizeAuthors(item, "en"))
	})
}

func TestGuessAuthorsByTag(t *testing.T) {
	s := &Syncer{}

	t.Run("prefers parsed gofeed authors over extensions", func(t *testing.T) {
		item := &gofeed.Item{
			Authors: []*gofeed.Person{{Name: "Parsed Author"}},
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{{Value: "Fallback Author"}},
				},
			},
		}

		assert.Equal(t, "Parsed Author", s.guessAuthorsByTag(item))
	})

	t.Run("extracts nested atom author name", func(t *testing.T) {
		item := &gofeed.Item{
			Extensions: ext.Extensions{
				"atom": {
					"author": []ext.Extension{{
						Children: map[string][]ext.Extension{
							"name": {{Value: "Atom Author"}},
						},
					}},
				},
			},
		}

		assert.Equal(t, "Atom Author", s.guessAuthorsByTag(item))
	})

	t.Run("returns empty string when no author tags exist", func(t *testing.T) {
		assert.Equal(t, "", s.guessAuthorsByTag(&gofeed.Item{}))
	})
}

func TestThinkRenderAndExtract_Authors(t *testing.T) {
	s := &Syncer{
		ctx:    context.Background(),
		logger: slog.Default(),
		think:  &mockThink{},
		feeds:  &mockFeeds{},
	}

	item := &gofeed.Item{
		Title:       "Author test",
		Description: "Desc",
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{{Value: "Alice & Bob"}},
			},
		},
	}

	result, err := s.thinkRenderAndExtract(item, "en", 0)
	assert.NoError(t, err)
	assert.Equal(t, database.StringArray{"Alice", "Bob"}, result.authors)
}
