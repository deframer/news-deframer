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

func TestRenderFeed_Complex(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	rssContent := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:content="http://purl.org/rss/1.0/modules/content/"
	xmlns:wfw="http://wellformedweb.org/CommentAPI/"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:atom="http://www.w3.org/2005/Atom"
	xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
	xmlns:slash="http://purl.org/rss/1.0/modules/slash/">
	<channel>
		<title>Global News Agency</title>
		<atom:link href="http://localhost:8003/feed/" rel="self" type="application/rss+xml" />
		<link>http://localhost:8003</link>
		<description></description>
		<lastBuildDate>Tue, 13 May 2025 03:49:00 +0000</lastBuildDate>
		<language>en-US</language>
		<sy:updatePeriod>hourly</sy:updatePeriod>
		<sy:updateFrequency>1</sy:updateFrequency>
		<generator>https://wordpress.org/?v=6.9</generator>
	</channel>
</rss>`

	feed, err := f.ParseFeed(ctx, strings.NewReader(rssContent))
	assert.NoError(t, err)

	xmlStr, err := f.RenderFeed(ctx, feed)
	assert.NoError(t, err)

	assert.Contains(t, xmlStr, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`)
	assert.Contains(t, xmlStr, `xmlns:wfw="http://wellformedweb.org/CommentAPI/"`)
	assert.Contains(t, xmlStr, `xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"`)
	assert.Contains(t, xmlStr, `<sy:updatePeriod>hourly</sy:updatePeriod>`)
	assert.Contains(t, xmlStr, `<atom:link href="http://localhost:8003/feed/" rel="self" type="application/rss+xml"`)

	// go one more round - parse this and render it and test it again
	feed2, err := f.ParseFeed(ctx, strings.NewReader(xmlStr))
	assert.NoError(t, err)

	xmlStr2, err := f.RenderFeed(ctx, feed2)
	assert.NoError(t, err)
	assert.Contains(t, xmlStr2, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`)
	assert.Contains(t, xmlStr2, `<sy:updatePeriod>hourly</sy:updatePeriod>`)
}

func TestGetValidItems(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	feed := &gofeed.Feed{
		Title: "Test Feed",
		Items: []*gofeed.Item{
			{Title: "Valid Item", Link: "http://example.com/1"},
			{Title: "Invalid Item", Link: ""}, // Missing link -> error in ItemHashKey
		},
	}

	processedFeed, pairs := f.GetValidItems(ctx, feed)

	assert.NotNil(t, processedFeed)
	assert.Equal(t, 1, len(processedFeed.Items))
	assert.Equal(t, "Valid Item", processedFeed.Items[0].Title)

	assert.Equal(t, 1, len(pairs))
	assert.Equal(t, "Valid Item", pairs[0].Item.Title)
	assert.NotEmpty(t, pairs[0].Key)
}

func TestItemHashKey(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		item1     *gofeed.Item
		type1     HashType
		item2     *gofeed.Item
		type2     HashType
		wantErr   bool
		wantEqual bool
	}{
		{
			name:      "Stability",
			item1:     &gofeed.Item{GUID: "test-guid", Title: "Test Title", Link: "http://example.com", PublishedParsed: &now},
			type1:     HashTypeDefault,
			item2:     &gofeed.Item{GUID: "test-guid", Title: "Test Title", Link: "http://example.com", PublishedParsed: &now},
			type2:     HashTypeDefault,
			wantEqual: true,
		},
		{
			name:      "Priority GUID",
			item1:     &gofeed.Item{GUID: "id1", Link: "http://link1", Title: "Title1"},
			type1:     HashTypeDefault,
			item2:     &gofeed.Item{GUID: "id1", Link: "http://link2", Title: "Title2"},
			type2:     HashTypeDefault,
			wantEqual: true,
		},
		{
			name:      "Priority Link",
			item1:     &gofeed.Item{GUID: "", Link: "http://link1", Title: "Title1"},
			type1:     HashTypeDefault,
			item2:     &gofeed.Item{GUID: "", Link: "http://link1", Title: "Title2"},
			type2:     HashTypeDefault,
			wantEqual: true,
		},
		{
			name:      "Whitespace Handling",
			item1:     &gofeed.Item{GUID: " id1 ", Link: " http://link1 ", Title: "Title1"},
			type1:     HashTypeDefault,
			item2:     &gofeed.Item{GUID: "id1", Link: "http://link1", Title: "Title1"},
			type2:     HashTypeDefault,
			wantEqual: true,
		},
		{
			name:  "Valid Item (Link only)",
			item1: &gofeed.Item{Title: "Some Title", Link: "http://example.com"},
			type1: HashTypeDefault,
		},
		{
			name:    "Empty Seed Error",
			item1:   &gofeed.Item{},
			type1:   HashTypeDefault,
			wantErr: true,
		},
		{
			name:    "Empty Link Error",
			item1:   &gofeed.Item{Title: "Title", GUID: "GUID"},
			type1:   HashTypeDefault,
			wantErr: true,
		},
		{
			name:      "Content Change",
			item1:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type1:     HashTypeVersioned,
			item2:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc B", Link: "http://example.com"},
			type2:     HashTypeVersioned,
			wantEqual: false,
		},
		{
			name:      "Title Change",
			item1:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type1:     HashTypeVersioned,
			item2:     &gofeed.Item{GUID: "id1", Title: "Title B", Description: "Desc A", Link: "http://example.com"},
			type2:     HashTypeVersioned,
			wantEqual: false,
		},
		{
			name:      "Same Content",
			item1:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type1:     HashTypeVersioned,
			item2:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type2:     HashTypeVersioned,
			wantEqual: true,
		},
		{
			name:      "Whitespace Content",
			item1:     &gofeed.Item{GUID: "id1", Title: " Title A ", Description: " Desc A ", Link: "http://example.com"},
			type1:     HashTypeVersioned,
			item2:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type2:     HashTypeVersioned,
			wantEqual: true,
		},
		{
			name:      "Stable vs Versioned",
			item1:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type1:     HashTypeDefault,
			item2:     &gofeed.Item{GUID: "id1", Title: "Title A", Description: "Desc A", Link: "http://example.com"},
			type2:     HashTypeVersioned,
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k1, err1 := ItemHashKey(tt.item1, tt.type1)
			if tt.wantErr {
				assert.Error(t, err1)
				if err1 != nil {
					assert.Contains(t, err1.Error(), "item link is empty")
				}
				return
			}
			assert.NoError(t, err1)

			if tt.item2 != nil {
				k2, err2 := ItemHashKey(tt.item2, tt.type2)
				assert.NoError(t, err2)

				if tt.wantEqual {
					assert.Equal(t, k1, k2)
				} else {
					assert.NotEqual(t, k1, k2)
				}
			}
		})
	}
}
