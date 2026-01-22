package syncer

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/feeds"
	"github.com/deframer/news-deframer/pkg/think"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

type mockThink struct {
	runFunc func(scope string, language string, req think.Request) (*database.ThinkResult, error)
}

func (m *mockThink) Run(scope string, language string, req think.Request) (*database.ThinkResult, error) {
	if m.runFunc != nil {
		return m.runFunc(scope, language, req)
	}
	return &database.ThinkResult{}, nil
}

func TestProcessItemErrorHandling(t *testing.T) {
	ctx := context.Background()
	logger := slog.Default()

	feedID := uuid.New()
	feed := &database.Feed{Base: database.Base{ID: feedID}}

	t.Run("IncrementErrorCountOnError", func(t *testing.T) {
		mockT := &mockThink{
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
				assert.Equal(t, 2, item.ThinkErrorCount) // Started with 1, should be 2
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

		// Simulate existing error count of 1
		s.processItem(feed, "hash1", item, "en", 1)
	})

	t.Run("ResetErrorCountOnSuccess", func(t *testing.T) {
		mockT := &mockThink{
			runFunc: func(scope string, language string, req think.Request) (*database.ThinkResult, error) {
				return &database.ThinkResult{TitleCorrected: "Cool"}, nil
			},
		}

		mockR := &mockRepo{
			upsertItemFunc: func(item *database.Item) error {
				assert.Equal(t, feedID, item.FeedID)
				assert.Equal(t, "hash2", item.Hash)
				assert.Nil(t, item.ThinkError)
				assert.Equal(t, 0, item.ThinkErrorCount) // Should be reset to 0
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

		// Simulate existing error count of 2
		s.processItem(feed, "hash2", item, "en", 2)
	})
}

// Helper struct for test (extending mockRepo in syncer_test.go which is limited)
// In a real scenario I would update the main mockRepo, but for this specific test file I'll inject the behavior if I can.
// However, Syncer uses database.Repository interface.
// I will just add the UpsertItemFunc to the mockRepo in syncer_test.go if possible or define a local one if the interface matches.
// Since mockRepo is in the same package (syncer_test.go), I can modify it there.
