package feeds

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
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
	assert.Contains(t, xmlStr, "<content:encoded><![CDATA[Full Content]]></content:encoded>")
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
	xmlns:media="http://search.yahoo.com/mrss/"
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
		<item>
			<title>My Title</title>
			<link>http://example.com/data</link>
			<dc:creator><![CDATA[admin]]></dc:creator>
			<pubDate>Tue, 13 May 2025 03:49:00 +0000</pubDate>
			<category domain="http://example.com/foo"><![CDATA[Latest News]]></category>
			<category domain="http://example.com/foo"><![CDATA[Banner Posts]]></category>
			<guid isPermaLink="false">http://example.com/data</guid>
			<description><![CDATA[foo]]></description>
			<media:content height="1800" medium="image" url="https://cdn/picture.jpg" width="1800"/>
			<media:credit>Foobar/Wordpress</media:credit>
			<media:description>A nice picture.</media:description>
			<content:encoded>
				<![CDATA[ <p> something in <br/> html </p> ]]>
			</content:encoded>
		</item>
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
	assert.Contains(t, xmlStr, `<title>My Title</title>`)
	assert.Contains(t, xmlStr, `<dc:creator>admin</dc:creator>`)
	assert.Contains(t, xmlStr, `<category domain="http://example.com/foo">Latest News</category>`)
	assert.Contains(t, xmlStr, `<category domain="http://example.com/foo">Banner Posts</category>`)
	assert.Contains(t, xmlStr, `<media:content height="1800" medium="image" url="https://cdn/picture.jpg" width="1800"></media:content>`)
	assert.Contains(t, xmlStr, `<content:encoded><![CDATA[ <p> something in <br/> html </p> ]]></content:encoded>`)

	// go one more round - parse this and render it and test it again
	feed2, err := f.ParseFeed(ctx, strings.NewReader(xmlStr))
	assert.NoError(t, err)

	xmlStr2, err := f.RenderFeed(ctx, feed2)
	assert.NoError(t, err)
	assert.Contains(t, xmlStr2, `xmlns:content="http://purl.org/rss/1.0/modules/content/"`)
	assert.Contains(t, xmlStr2, `<sy:updatePeriod>hourly</sy:updatePeriod>`)
	assert.Contains(t, xmlStr2, `<title>My Title</title>`)
	assert.Contains(t, xmlStr2, `<dc:creator>admin</dc:creator>`)
	assert.Contains(t, xmlStr2, `<category domain="http://example.com/foo">Latest News</category>`)
	assert.Contains(t, xmlStr2, `<category domain="http://example.com/foo">Banner Posts</category>`)
	assert.Contains(t, xmlStr2, `<media:content height="1800" medium="image" url="https://cdn/picture.jpg" width="1800"></media:content>`)
	assert.Contains(t, xmlStr2, `<content:encoded><![CDATA[ <p> something in <br/> html </p> ]]></content:encoded>`)

	validateItemRenderConsistency(t, ctx, f, feed.Items[0])
}

func TestRenderItem(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	t.Run("Standard Fields", func(t *testing.T) {
		item := &gofeed.Item{
			Title:       "Test Title",
			Link:        "http://example.com",
			Description: "Test Description",
			GUID:        "test-guid",
			Published:   "Mon, 02 Jan 2006 15:04:05 MST",
			Content:     "<p>Full Content</p>",
		}

		xmlStr, err := f.RenderItem(ctx, item)
		assert.NoError(t, err)

		assert.Contains(t, xmlStr, "<title>Test Title</title>")
		assert.Contains(t, xmlStr, "<link>http://example.com</link>")
		assert.Contains(t, xmlStr, "<description>Test Description</description>")
		assert.Contains(t, xmlStr, "<guid>test-guid</guid>")
		assert.Contains(t, xmlStr, "<pubDate>Mon, 02 Jan 2006 15:04:05 MST</pubDate>")
		// Content mapped to content:encoded and wrapped in CDATA
		assert.Contains(t, xmlStr, "<content:encoded><![CDATA[<p>Full Content</p>]]></content:encoded>")

		validateItemRenderConsistency(t, ctx, f, item)
	})

	t.Run("With Extensions", func(t *testing.T) {
		item := &gofeed.Item{
			Title: "Extension Test",
			Extensions: ext.Extensions{
				"dc": {
					"creator": []ext.Extension{
						{Value: "John Doe"},
					},
				},
				"foo": {
					"bar": []ext.Extension{
						{
							Value: "baz",
							Attrs: map[string]string{"type": "test"},
						},
					},
				},
			},
		}

		xmlStr, err := f.RenderItem(ctx, item)
		assert.NoError(t, err)

		assert.Contains(t, xmlStr, "<dc:creator>John Doe</dc:creator>")
		// Check for attribute and value
		assert.Contains(t, xmlStr, `<foo:bar type="test">baz</foo:bar>`, "Should contain foo:bar with attribute")

		validateItemRenderConsistency(t, ctx, f, item)
	})

	t.Run("Nested Extensions", func(t *testing.T) {
		item := &gofeed.Item{
			Title: "Nested Test",
			Extensions: ext.Extensions{
				"media": {
					"group": []ext.Extension{
						{
							Children: map[string][]ext.Extension{
								"content": {
									{
										Attrs: map[string]string{"url": "http://example.com/img.jpg"},
									},
								},
							},
						},
					},
				},
			},
		}

		xmlStr, err := f.RenderItem(ctx, item)
		assert.NoError(t, err)

		assert.Contains(t, xmlStr, "<media:group>")
		assert.Contains(t, xmlStr, `content url="http://example.com/img.jpg"`)

		validateItemRenderConsistency(t, ctx, f, item)
	})
}

func TestFilterItems(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{}
	f := NewFeeds(ctx, cfg)

	tests := []struct {
		name           string
		feed           *gofeed.Feed
		domains        []string
		expectedCount  int
		expectedTitles []string
	}{
		{
			name: "No domains filter",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://example.com/1"},
					{Title: "Item 2", Link: "http://other.com/1"},
				},
			},
			domains:        nil,
			expectedCount:  2,
			expectedTitles: []string{"Item 1", "Item 2"},
		},
		{
			name: "Filter by domain",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://example.com/1"},
					{Title: "Item 2", Link: "http://other.com/1"},
				},
			},
			domains:        []string{"example.com"},
			expectedCount:  1,
			expectedTitles: []string{"Item 1"},
		},
		{
			name: "Filter by subdomain",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://sub.example.com/1"},
					{Title: "Item 2", Link: "http://example.com/1"},
					{Title: "Item 3", Link: "http://other.com/1"},
				},
			},
			domains:        []string{"example.com"},
			expectedCount:  2,
			expectedTitles: []string{"Item 1", "Item 2"},
		},
		{
			name: "Invalid Link",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: ""},
				},
			},
			domains:        nil,
			expectedCount:  0,
			expectedTitles: []string{},
		},
		{
			name:           "Nil Feed",
			feed:           nil,
			domains:        nil,
			expectedCount:  0,
			expectedTitles: []string{},
		},
		{
			name: "Multiple allowed domains",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://a.com/1"},
					{Title: "Item 2", Link: "http://b.com/1"},
					{Title: "Item 3", Link: "http://c.com/1"},
				},
			},
			domains:        []string{"a.com", "b.com"},
			expectedCount:  2,
			expectedTitles: []string{"Item 1", "Item 2"},
		},
		{
			name: "Public Suffix Domain (example.co.uk)",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://sub.example.co.uk/1"},
					{Title: "Item 2", Link: "http://example.co.uk/1"},
					{Title: "Item 3", Link: "http://other.co.uk/1"},
				},
			},
			domains:        []string{"example.co.uk"},
			expectedCount:  2,
			expectedTitles: []string{"Item 1", "Item 2"},
		},
		{
			name: "Duplicate Items",
			feed: &gofeed.Feed{
				Items: []*gofeed.Item{
					{Title: "Item 1", Link: "http://example.com/1", GUID: "1"},
					{Title: "Item 1 Duplicate", Link: "http://example.com/1", GUID: "1"},
				},
			},
			domains:        nil,
			expectedCount:  1,
			expectedTitles: []string{"Item 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pairs := f.FilterItems(ctx, tt.feed, tt.domains)

			if tt.feed == nil {
				assert.Nil(t, pairs)
				return
			}

			assert.Equal(t, tt.expectedCount, len(pairs))

			for i, pair := range pairs {
				assert.Equal(t, tt.expectedTitles[i], pair.Item.Title)
				assert.NotEmpty(t, pair.Hash)
			}
		})
	}
}

func validateItemRenderConsistency(t *testing.T, ctx context.Context, f Feeds, item *gofeed.Item) {
	itemXML, err := f.RenderItem(ctx, item)
	assert.NoError(t, err)

	// Wrap in a feed to get RenderFeed output
	feed := &gofeed.Feed{
		Title: "Wrapper",
		Items: []*gofeed.Item{item},
	}

	feedXML, err := f.RenderFeed(ctx, feed)
	assert.NoError(t, err)

	// Simple normalization: remove all whitespace to ignore indentation differences
	clean := func(s string) string {
		return strings.Map(func(r rune) rune {
			if r == ' ' || r == '\n' || r == '\t' {
				return -1
			}
			return r
		}, s)
	}

	assert.Contains(t, clean(feedXML), clean(itemXML), "RenderFeed output should contain RenderItem output")
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

func TestSetExtension(t *testing.T) {
	item := &gofeed.Item{}
	SetExtension(item, "testprefix", "testname", "testvalue")

	assert.NotNil(t, item.Extensions)
	assert.Contains(t, item.Extensions, "testprefix")
	assert.Contains(t, item.Extensions["testprefix"], "testname")
	assert.Equal(t, "testvalue", item.Extensions["testprefix"]["testname"][0].Value)
	assert.Equal(t, "testname", item.Extensions["testprefix"]["testname"][0].Name)
}

func TestAddNamespace(t *testing.T) {
	feed := &gofeed.Feed{}
	AddNamespace(feed, "xmlns:test", "http://test.com")

	assert.NotNil(t, feed.Custom)
	assert.Equal(t, "http://test.com", feed.Custom["xmlns:test"])
}
