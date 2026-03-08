package syncer

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/feeds"
	"github.com/deframer/news-deframer/pkg/think"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/stretchr/testify/assert"
)

type mockThinkEH struct {
	runFunc func(scope string, language string, req think.Request) (*database.ThinkResult, error)
}

func (m *mockThinkEH) Run(scope string, language string, req think.Request) (*database.ThinkResult, error) {
	if m.runFunc != nil {
		return m.runFunc(scope, language, req)
	}
	return &database.ThinkResult{}, nil
}

func TestProcessItem_Categories(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThink{}
	mockFeeds := &mockFeeds{}
	s.feeds = mockFeeds

	feed := &database.Feed{Base: database.Base{ID: uuid.New()}}
	hash := "test-hash"
	item := &gofeed.Item{
		Title: "Test item with categories",
		Link:  "http://example.com/item",
	}

	var capturedItem *database.Item
	repo.upsertItemFunc = func(dbItem *database.Item) error {
		capturedItem = dbItem
		return nil
	}

	mockFeeds.extractCategoriesFunc = func(item *gofeed.Item) []string {
		return []string{"cat1", "cat2", "cat3"}
	}

	s.processItem(feed, hash, item, "en", 0)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"cat1", "cat2", "cat3"}, capturedItem.Categories)
}

func TestProcessItem_Authors(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThink{}
	s.feeds = &mockFeeds{}

	feed := &database.Feed{Base: database.Base{ID: uuid.New()}}
	item := &gofeed.Item{
		Title: "Test item with authors",
		Link:  "http://example.com/item",
		Extensions: ext.Extensions{
			"dc": {
				"creator": []ext.Extension{{Value: "Alice & Bob / Carol"}},
			},
		},
	}

	var capturedItem *database.Item
	repo.upsertItemFunc = func(dbItem *database.Item) error {
		capturedItem = dbItem
		return nil
	}

	s.processItem(feed, "test-hash", item, "en", 0)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"Alice", "Bob", "Carol"}, capturedItem.Authors)
}

func TestProcessThinkerItem_Authors(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThink{}
	s.feeds = feeds.NewFeeds(context.Background(), cfg)

	dbItem := &database.Item{
		ID:     uuid.New(),
		FeedID: uuid.New(),
		URL:    "http://example.com/item",
		Content: strings.TrimSpace(`
			<item>
				<title>Thinker item</title>
				<link>http://example.com/item</link>
				<dc:creator>Alice | Example News</dc:creator>
				<description>hello</description>
			</item>`),
		Authors: database.StringArray{"stale"},
	}

	var capturedItem *database.Item
	repo.upsertItemFunc = func(updated *database.Item) error {
		capturedItem = updated
		return nil
	}

	s.processThinkerItem(dbItem)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"Alice"}, capturedItem.Authors)
}

func TestProcessItemErrorHandling(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	feedID := uuid.New()
	feed := &database.Feed{Base: database.Base{ID: feedID}}

	t.Run("IncrementErrorCountOnError", func(t *testing.T) {
		mockT := &mockThinkEH{
			runFunc: func(scope string, language string, req think.Request) (*database.ThinkResult, error) {
				return nil, errors.New("thinking failed")
			},
		}

		mockR := &mockRepo{
			upsertItemFunc: func(item *database.Item) error {
				assert.Equal(t, feedID, item.FeedID)
				assert.Equal(t, "hash1", item.Hash)
				assert.NotNil(t, item.ThinkError)
				assert.Equal(t, "thinking failed", *item.ThinkError)
				assert.Equal(t, 2, item.ThinkErrorCount)
				return nil
			},
		}

		s := &Syncer{
			ctx:    ctx,
			repo:   mockR,
			logger: logger,
			think:  mockT,
			feeds:  feeds.NewFeeds(ctx, &config.Config{}),
		}

		item := &gofeed.Item{
			Title:       "Test Item",
			Description: "Desc",
			Link:        "http://example.com/1",
		}

		s.processItem(feed, "hash1", item, "en", 1)
	})

	t.Run("ResetErrorCountOnSuccess", func(t *testing.T) {
		mockT := &mockThinkEH{
			runFunc: func(scope string, language string, req think.Request) (*database.ThinkResult, error) {
				return &database.ThinkResult{TitleCorrected: "Cool"}, nil
			},
		}

		mockR := &mockRepo{
			upsertItemFunc: func(item *database.Item) error {
				assert.Equal(t, feedID, item.FeedID)
				assert.Equal(t, "hash2", item.Hash)
				assert.Nil(t, item.ThinkError)
				assert.Equal(t, 0, item.ThinkErrorCount)
				return nil
			},
		}

		s := &Syncer{
			ctx:    ctx,
			repo:   mockR,
			logger: logger,
			think:  mockT,
			feeds:  feeds.NewFeeds(ctx, &config.Config{}),
		}

		item := &gofeed.Item{
			Title:       "Test Item 2",
			Description: "Desc",
			Link:        "http://example.com/2",
		}

		s.processItem(feed, "hash2", item, "en", 2)
	})
}
