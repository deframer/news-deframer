package syncer

import (
	"context"
	"errors"
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
	runFunc func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error)
}

func (m *mockThinkEH) Run(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
	if m.runFunc != nil {
		return m.runFunc(scope, language, req, ignoreCategoryErrors)
	}
	return &database.ThinkResult{}, nil
}

func TestProcessItem_Categories(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThinkEH{runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
		t.Fatal("syncItem must not call think")
		return nil, nil
	}}
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

	s.syncItem(feed, hash, item, "en")

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"cat1", "cat2", "cat3"}, capturedItem.Categories)
}

func TestProcessItem_Authors(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThinkEH{runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
		t.Fatal("syncItem must not call think")
		return nil, nil
	}}
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

	s.syncItem(feed, "test-hash", item, "en")

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"Alice", "Bob", "Carol"}, capturedItem.Authors)
}

func TestProcessItem_EmptyAuthors(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThinkEH{runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
		t.Fatal("syncItem must not call think")
		return nil, nil
	}}
	s.feeds = &mockFeeds{}

	feed := &database.Feed{Base: database.Base{ID: uuid.New()}}
	item := &gofeed.Item{
		Title: "Test item without authors",
		Link:  "http://example.com/item",
	}

	var capturedItem *database.Item
	repo.upsertItemFunc = func(dbItem *database.Item) error {
		capturedItem = dbItem
		return nil
	}

	s.syncItem(feed, "test-hash", item, "en")

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{}, capturedItem.Authors)
}

func TestProcessItem_FromRSSXML(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThinkEH{runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
		t.Fatal("syncItem must not call think")
		return nil, nil
	}}
	s.feeds = &mockFeeds{}

	rssContent := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <item>
      <title>Example title</title>
      <link>https://example.test/story#ref=rss</link>
      <description>Example description</description>
      <content:encoded>Example description</content:encoded>
      <category>Example category</category>
      <enclosure type="image/jpeg" url="https://example.test/image.jpg" length="0"/>
      <guid isPermaLink="false">example-guid</guid>
      <pubDate>Sat, 16 May 2026 16:38:00 +0200</pubDate>
    </item>
  </channel>
</rss>`

	feed, err := gofeed.NewParser().Parse(strings.NewReader(rssContent))
	assert.NoError(t, err)
	if assert.Len(t, feed.Items, 1) {
		var capturedItem *database.Item
		repo.upsertItemFunc = func(dbItem *database.Item) error {
			capturedItem = dbItem
			return nil
		}

		baseFeed := &database.Feed{Base: database.Base{ID: uuid.New()}}
		s.syncItem(baseFeed, "test-hash", feed.Items[0], "de")

		if assert.NotNil(t, capturedItem) {
			assert.Equal(t, baseFeed.ID, capturedItem.FeedID)
			assert.Equal(t, "test-hash", capturedItem.Hash)
			assert.Equal(t, "https://example.test/story", capturedItem.URL)
			assert.NotNil(t, capturedItem.MediaContent)
			if assert.NotNil(t, capturedItem.MediaContent) {
				assert.Equal(t, "https://example.test/image.jpg", capturedItem.MediaContent.URL)
				assert.Equal(t, "image/jpeg", capturedItem.MediaContent.Type)
				assert.Equal(t, "image", capturedItem.MediaContent.Medium)
				assert.Equal(t, 1920, capturedItem.MediaContent.Width)
				assert.Equal(t, 1080, capturedItem.MediaContent.Height)
				assert.False(t, capturedItem.MediaContent.HasExplicitDims)
			}
		}
	}
}

func TestThinkItem_PreservesMediaContentFromSyncedItem(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	s.think = &mockThinkEH{runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
		return &database.ThinkResult{
			TitleCorrected:       "Corrected title",
			DescriptionCorrected: "Corrected description",
			OverallReason:        "Reason",
		}, nil
	}}
	s.feeds = feeds.NewFeeds(context.Background(), cfg)

	rssContent := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <item>
      <title>Example title</title>
      <link>https://example.test/story#ref=rss</link>
      <description>Example description</description>
      <content:encoded>Example description</content:encoded>
      <category>Example category</category>
      <enclosure type="image/jpeg" url="https://example.test/image.jpg" length="0"/>
      <guid isPermaLink="false">example-guid</guid>
      <pubDate>Sat, 16 May 2026 16:38:00 +0200</pubDate>
    </item>
  </channel>
</rss>`

	feed, err := gofeed.NewParser().Parse(strings.NewReader(rssContent))
	assert.NoError(t, err)
	if assert.Len(t, feed.Items, 1) {
		feedMeta := &database.Feed{Base: database.Base{ID: uuid.New()}}

		var syncedItem *database.Item
		repo.upsertItemFunc = func(dbItem *database.Item) error {
			syncedItem = dbItem
			return nil
		}

		s.syncItem(feedMeta, "test-hash", feed.Items[0], "de")

		if assert.NotNil(t, syncedItem) {
			assert.Equal(t, "https://example.test/story", syncedItem.URL)
			assert.NotNil(t, syncedItem.MediaContent)
			assert.Equal(t, "https://example.test/image.jpg", syncedItem.MediaContent.URL)
			assert.False(t, syncedItem.MediaContent.HasExplicitDims)
			assert.NotContains(t, syncedItem.Content, "enclosure")
		}

		var finalItem *database.Item
		repo.upsertItemInvalidateFunc = func(updated *database.Item) error {
			finalItem = updated
			return nil
		}

		s.thinkItem(syncedItem)

		if assert.NotNil(t, finalItem) {
			assert.NotNil(t, finalItem.MediaContent)
			assert.Equal(t, "https://example.test/image.jpg", finalItem.MediaContent.URL)
			assert.Equal(t, "image/jpeg", finalItem.MediaContent.Type)
			assert.Equal(t, "image", finalItem.MediaContent.Medium)
			assert.Equal(t, 1920, finalItem.MediaContent.Width)
			assert.Equal(t, 1080, finalItem.MediaContent.Height)
			assert.False(t, finalItem.MediaContent.HasExplicitDims)
		}
	}
}

func TestThinkItem_Authors(t *testing.T) {
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
	repo.upsertItemInvalidateFunc = func(updated *database.Item) error {
		capturedItem = updated
		return nil
	}

	s.thinkItem(dbItem)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"Alice"}, capturedItem.Authors)
}

func TestThinkItem_EmptyAuthors(t *testing.T) {
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
				<description>hello</description>
			</item>`),
		Authors: database.StringArray{"stale"},
	}

	var capturedItem *database.Item
	repo.upsertItemInvalidateFunc = func(updated *database.Item) error {
		capturedItem = updated
		return nil
	}

	s.thinkItem(dbItem)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{}, capturedItem.Authors)
}

func TestThinkItem_UsesTrendInvalidatingUpsert(t *testing.T) {
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
				<description>hello</description>
			</item>`),
	}

	plainCalled := false
	invalidateCalled := false
	repo.upsertItemFunc = func(updated *database.Item) error {
		plainCalled = true
		return nil
	}
	repo.upsertItemInvalidateFunc = func(updated *database.Item) error {
		invalidateCalled = true
		return nil
	}

	s.thinkItem(dbItem)

	assert.False(t, plainCalled)
	assert.True(t, invalidateCalled)
}

func TestThinkItemErrorHandling(t *testing.T) {
	ctx := context.Background()
	strPtr := func(s string) *string { return &s }

	feedID := uuid.New()

	t.Run("IncrementErrorCountOnError", func(t *testing.T) {
		mockT := &mockThinkEH{
			runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
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
			ctx:   ctx,
			repo:  mockR,
			think: mockT,
			feeds: feeds.NewFeeds(ctx, &config.Config{}),
		}

		item := &gofeed.Item{
			Title:       "Test Item",
			Description: "Desc",
			Link:        "http://example.com/1",
		}

		s.thinkItem(&database.Item{ID: uuid.New(), Hash: "hash1", FeedID: feedID, URL: item.Link, Content: "<item><title>Test Item</title><description>Desc</description></item>", ThinkErrorCount: 1, Language: strPtr("en")})
	})

	t.Run("ResetErrorCountOnSuccess", func(t *testing.T) {
		mockT := &mockThinkEH{
			runFunc: func(scope string, language string, req think.Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
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
			ctx:   ctx,
			repo:  mockR,
			think: mockT,
			feeds: feeds.NewFeeds(ctx, &config.Config{}),
		}

		item := &gofeed.Item{
			Title:       "Test Item 2",
			Description: "Desc",
			Link:        "http://example.com/2",
		}

		s.thinkItem(&database.Item{ID: uuid.New(), Hash: "hash2", FeedID: feedID, URL: item.Link, Content: "<item><title>Test Item 2</title><description>Desc</description></item>", ThinkErrorCount: 2, Language: strPtr("en")})
	})
}
