package feeds

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/mmcdole/gofeed"
)

type Feeds interface {
	ParseFeed(ctx context.Context, content io.Reader) (*gofeed.Feed, error)
	RenderFeed(ctx context.Context, feed *gofeed.Feed) (string, error)
}

type feeds struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
}

// NewFeeds initializes a new feeds service
func NewFeeds(ctx context.Context, cfg *config.Config) Feeds {
	return &feeds{
		ctx:    ctx,
		cfg:    cfg,
		logger: slog.With("component", "feeds"),
	}
}

// ParseFeed parses the raw feed content
func (f *feeds) ParseFeed(ctx context.Context, content io.Reader) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}
	return feed, nil
}

// RenderFeed renders a gofeed.Feed to an RSS 2.0 XML string
func (f *feeds) RenderFeed(ctx context.Context, feed *gofeed.Feed) (string, error) {
	if feed == nil {
		return "", fmt.Errorf("feed cannot be nil")
	}

	rss := &rss2{
		Version:    "2.0",
		ContentNS:  "http://purl.org/rss/1.0/modules/content/",
		DeframerNS: "https://github.com/egandro/news-deframer/",
		Channel: rss2Channel{
			Title:       feed.Title,
			Link:        feed.Link,
			Description: feed.Description,
		},
	}

	for _, item := range feed.Items {
		rss.Channel.Items = append(rss.Channel.Items, f.toRSS2Item(item))
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal feed: %w", err)
	}

	return xml.Header + string(output), nil
}

func (f *feeds) toRSS2Item(item *gofeed.Item) rss2Item {
	var custom string
	if item.Custom != nil {
		custom = item.Custom["custom_field"]
	}

	return rss2Item{
		Title:       item.Title,
		Link:        item.Link,
		Description: item.Description,
		Content:     item.Content,
		PubDate:     item.Published,
		GUID:        item.GUID,
		CustomField: custom,
	}
}

// Internal structs for XML marshalling
type rss2 struct {
	XMLName    xml.Name    `xml:"rss"`
	Version    string      `xml:"version,attr"`
	ContentNS  string      `xml:"xmlns:content,attr"`
	DeframerNS string      `xml:"xmlns:deframer,attr"`
	Channel    rss2Channel `xml:"channel"`
}

type rss2Channel struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Items       []rss2Item `xml:"item"`
}

type rss2Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title,omitempty"`
	Link        string   `xml:"link,omitempty"`
	Description string   `xml:"description,omitempty"`
	Content     string   `xml:"content:encoded,omitempty"`
	PubDate     string   `xml:"pubDate,omitempty"`
	GUID        string   `xml:"guid,omitempty"`
	CustomField string   `xml:"deframer:custom,omitempty"`
}
