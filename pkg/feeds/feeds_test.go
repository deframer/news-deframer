package feeds

import (
	"context"
	"strings"
	"testing"

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
