package feeds

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"strings"

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

// HashType defines the strategy for generating the item hash
type HashType int

const (
	// HashTypeDefault generates a hash based only on the item's identity (GUID/Link/Title)
	HashTypeDefault HashType = iota
	// HashTypeVersioned generates a hash based on identity + content (Title + Description)
	HashTypeVersioned
)

// ItemHashKey creates a deterministic ID based on the item's attributes.
// By default, it uses HashTypeDefault. You can pass a specific HashType as a second argument.
func ItemHashKey(item *gofeed.Item, ht ...HashType) (string, error) {
	mode := HashTypeDefault
	if len(ht) > 0 {
		mode = ht[0]
	}

	link := strings.TrimSpace(item.Link)
	if link == "" {
		return "", fmt.Errorf("item link is empty")
	}

	seed := strings.TrimSpace(item.GUID)
	if seed == "" {
		seed = link
	}

	if mode == HashTypeVersioned {
		cleanTitle := strings.TrimSpace(item.Title)
		cleanDesc := strings.TrimSpace(item.Description)
		// Combine identity with content to detect changes
		seed = fmt.Sprintf("%s|%s|%s", seed, cleanTitle, cleanDesc)
	}

	if seed == "" {
		return "", fmt.Errorf("seed is empty")
	}

	hash := sha256.Sum256([]byte(seed))
	key := hex.EncodeToString(hash[:])

	return key, nil
}
