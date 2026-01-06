package feeds

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestParseFeed(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	rssContent := `
	<rss version="2.0">
		<channel>
			<title>Test Feed</title>
			<link>http://example.com</link>
			<description>Test Description</description>
			<item>
				<title>Test Item</title>
				<link>http://example.com/item</link>
				<description>Test Item Description</description>
			</item>
		</channel>
	</rss>`

	feed, err := f.ParseFeed(ctx, strings.NewReader(rssContent))
	assert.NoError(t, err)
	assert.NotNil(t, feed)
	assert.Equal(t, "Test Feed", feed.Title)
	assert.Equal(t, "Test Description", feed.Description)
	assert.Len(t, feed.Items, 1)
	assert.Equal(t, "Test Item", feed.Items[0].Title)
}

func TestParseFeed_Invalid(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	_, err := f.ParseFeed(ctx, strings.NewReader("invalid xml"))
	assert.Error(t, err)
}

func TestRenderFeed(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	feed := &gofeed.Feed{
		Title:       "Test Feed",
		Link:        "http://example.com",
		Description: "Test Description",
		Items: []*gofeed.Item{
			{
				Title:       "Test Item",
				Link:        "http://example.com/item",
				Description: "Test Item Description",
				Content:     "Full Content",
				Published:   "Mon, 02 Jan 2006 15:04:05 MST",
				GUID:        "12345",
			},
		},
	}

	xmlStr, err := f.RenderFeed(ctx, feed)
	assert.NoError(t, err)
	assert.Contains(t, xmlStr, "<rss version=\"2.0\"")
	assert.Contains(t, xmlStr, "<title>Test Feed</title>")
	assert.Contains(t, xmlStr, "<content:encoded>Full Content</content:encoded>")
}

func TestItemHashKey_Stable(t *testing.T) {
	now := time.Now()

	t.Run("Stability", func(t *testing.T) {
		item := &gofeed.Item{
			GUID:            "test-guid",
			Title:           "Test Title",
			Link:            "http://example.com",
			PublishedParsed: &now,
		}
		id1, err := ItemHashKey(item) // Test default (Stable)
		assert.NoError(t, err)
		id2, err := ItemHashKey(item, HashTypeStable)
		assert.NoError(t, err)
		assert.Equal(t, id1, id2)
	})

	t.Run("Priority GUID", func(t *testing.T) {
		item1 := &gofeed.Item{GUID: "id1", Link: "http://link1", Title: "Title1"}
		item2 := &gofeed.Item{GUID: "id1", Link: "http://link2", Title: "Title2"}
		// Should be same because GUID is same
		id1, err := ItemHashKey(item1)
		assert.NoError(t, err)
		id2, err := ItemHashKey(item2)
		assert.NoError(t, err)
		assert.Equal(t, id1, id2)
	})

	t.Run("Priority Link", func(t *testing.T) {
		item1 := &gofeed.Item{GUID: "", Link: "http://link1", Title: "Title1"}
		item2 := &gofeed.Item{GUID: "", Link: "http://link1", Title: "Title2"}
		// Should be same because Link is same and GUID is empty
		id1, err := ItemHashKey(item1)
		assert.NoError(t, err)
		id2, err := ItemHashKey(item2)
		assert.NoError(t, err)
		assert.Equal(t, id1, id2)
	})

	t.Run("Fallback Date", func(t *testing.T) {
		item := &gofeed.Item{Title: "Some Title", Link: "http://example.com"}
		_, err := ItemHashKey(item)
		assert.NoError(t, err)
	})

	t.Run("Empty Seed Error", func(t *testing.T) {
		item := &gofeed.Item{}
		_, err := ItemHashKey(item, HashTypeStable)
		assert.Error(t, err)
	})

	t.Run("Empty Link Error", func(t *testing.T) {
		item := &gofeed.Item{Title: "Title", GUID: "GUID"}
		_, err := ItemHashKey(item)
		assert.Error(t, err)
	})
}

func TestItemHashKey_Versioned(t *testing.T) {
	t.Run("Content Change", func(t *testing.T) {
		item1 := &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"}
		item2 := &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc B", Link: "http://example.com"}

		k1, err := ItemHashKey(item1, HashTypeVersioned)
		assert.NoError(t, err)
		k2, err := ItemHashKey(item2, HashTypeVersioned)
		assert.NoError(t, err)
		assert.NotEqual(t, k1, k2)
	})

	t.Run("Title Change", func(t *testing.T) {
		item1 := &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"}
		item2 := &gofeed.Item{GUID: "id1", Title: "Title B", Description: "Desc A", Link: "http://example.com"}

		k1, err := ItemHashKey(item1, HashTypeVersioned)
		assert.NoError(t, err)
		k2, err := ItemHashKey(item2, HashTypeVersioned)
		assert.NoError(t, err)
		assert.NotEqual(t, k1, k2)
	})

	t.Run("Same Content", func(t *testing.T) {
		item1 := &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"}
		item2 := &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"}

		k1, err := ItemHashKey(item1, HashTypeVersioned)
		assert.NoError(t, err)
		k2, err := ItemHashKey(item2, HashTypeVersioned)
		assert.NoError(t, err)
		assert.Equal(t, k1, k2)
	})
}
