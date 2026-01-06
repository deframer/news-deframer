package feeds

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strings"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

type Feeds interface {
	ParseFeed(ctx context.Context, content io.Reader) (*gofeed.Feed, error)
	RenderFeed(ctx context.Context, feed *gofeed.Feed) (string, error)
	GetValidItems(ctx context.Context, feed *gofeed.Feed) (*gofeed.Feed, []ItemKeyPair)
}

type feeds struct {
	ctx    context.Context
	cfg    *config.Config
	logger *slog.Logger
}

// ItemKeyPair holds a feed item and its generated key
type ItemKeyPair struct {
	Key  string
	Item *gofeed.Item
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
	data, err := io.ReadAll(content)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed content: %w", err)
	}

	// scan for root attributes to preserve namespaces
	rootAttrs := make(map[string]string)
	decoder := xml.NewDecoder(bytes.NewReader(data))
	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}
		if se, ok := t.(xml.StartElement); ok {
			for _, attr := range se.Attr {
				if attr.Name.Local == "version" && attr.Name.Space == "" {
					continue
				}
				key := attr.Name.Local
				if attr.Name.Space != "" {
					key = fmt.Sprintf("%s:%s", attr.Name.Space, attr.Name.Local)
				}
				rootAttrs[key] = attr.Value
			}
			break
		}
	}

	fp := gofeed.NewParser()
	feed, err := fp.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	if feed.Custom == nil {
		feed.Custom = make(map[string]string)
	}
	for k, v := range rootAttrs {
		feed.Custom[k] = v
	}

	return feed, nil
}

// RenderFeed renders a gofeed.Feed to an RSS 2.0 XML string
func (f *feeds) RenderFeed(ctx context.Context, feed *gofeed.Feed) (string, error) {
	if feed == nil {
		return "", fmt.Errorf("feed cannot be nil")
	}

	rss := &rss2{
		Version: "2.0",
		Custom:  feed.Custom,
		Channel: rss2Channel{
			Title:       feed.Title,
			Link:        feed.Link,
			Description: feed.Description,
			Extensions:  feed.Extensions,
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

// GetValidItems iterates over the items, calculates the hash.
// If there is an error calculating the hash, the item is removed from the feed.
// It returns a new feed with valid items and a list of key/item pairs.
func (f *feeds) GetValidItems(ctx context.Context, feed *gofeed.Feed) (*gofeed.Feed, []ItemKeyPair) {
	if feed == nil {
		return nil, nil
	}

	newFeed := *feed
	newFeed.Items = make([]*gofeed.Item, 0, len(feed.Items))
	var results []ItemKeyPair

	for _, item := range feed.Items {
		key, err := ItemHashKey(item)
		if err != nil {
			f.logger.DebugContext(ctx, "skipping item", "error", err, "item_title", item.Title)
			continue
		}
		newFeed.Items = append(newFeed.Items, item)
		results = append(results, ItemKeyPair{Key: key, Item: item})
	}

	return &newFeed, results
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
		Extensions:  item.Extensions,
		Custom:      item.Custom,
	}
}

// Internal structs for XML marshalling
type rss2 struct {
	Version string            `xml:"version,attr"`
	Custom  map[string]string `xml:"-"`
	Channel rss2Channel       `xml:"channel"`
}

func (r *rss2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "rss"
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "version"}, Value: r.Version})

	defaults := map[string]string{
		"xmlns:content":  "http://purl.org/rss/1.0/modules/content/",
		"xmlns:deframer": "https://github.com/egandro/news-deframer/",
		"xmlns:dc":       "http://purl.org/dc/elements/1.1/",
		"xmlns:atom":     "http://www.w3.org/2005/Atom",
	}

	var attrs []xml.Attr
	for k, v := range defaults {
		if _, ok := r.Custom[k]; !ok {
			attrs = append(attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
		}
	}

	for k, v := range r.Custom {
		if k == "version" {
			continue
		}
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Name.Local < attrs[j].Name.Local
	})
	start.Attr = append(start.Attr, attrs...)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeElement(r.Channel, xml.StartElement{Name: xml.Name{Local: "channel"}}); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type rss2Channel struct {
	Title       string            `xml:"title"`
	Link        string            `xml:"link"`
	Description string            `xml:"description"`
	Items       []rss2Item        `xml:"item"`
	Extensions  ext.Extensions    `xml:"-"`
	Custom      map[string]string `xml:"-"`
}

func (r rss2Channel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if r.Title != "" {
		if err := e.EncodeElement(r.Title, xml.StartElement{Name: xml.Name{Local: "title"}}); err != nil {
			return err
		}
	}
	if r.Link != "" {
		if err := e.EncodeElement(r.Link, xml.StartElement{Name: xml.Name{Local: "link"}}); err != nil {
			return err
		}
	}
	if r.Description != "" {
		if err := e.EncodeElement(r.Description, xml.StartElement{Name: xml.Name{Local: "description"}}); err != nil {
			return err
		}
	}

	if err := encodeExtensions(e, r.Extensions); err != nil {
		return err
	}

	for _, item := range r.Items {
		if err := e.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: "item"}}); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type rss2Item struct {
	XMLName     xml.Name          `xml:"item"`
	Title       string            `xml:"title,omitempty"`
	Link        string            `xml:"link,omitempty"`
	Description string            `xml:"description,omitempty"`
	Content     string            `xml:"content:encoded,omitempty"`
	PubDate     string            `xml:"pubDate,omitempty"`
	GUID        string            `xml:"guid,omitempty"`
	CustomField string            `xml:"deframer:custom,omitempty"`
	Extensions  ext.Extensions    `xml:"-"`
	Custom      map[string]string `xml:"-"`
}

func (r rss2Item) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var attrs []xml.Attr
	for k, v := range r.Custom {
		attrs = append(attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Name.Local < attrs[j].Name.Local
	})
	start.Attr = append(start.Attr, attrs...)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	encode := func(name, val string) error {
		if val == "" {
			return nil
		}
		return e.EncodeElement(val, xml.StartElement{Name: xml.Name{Local: name}})
	}

	if err := encode("title", r.Title); err != nil {
		return err
	}
	if err := encode("link", r.Link); err != nil {
		return err
	}
	if err := encode("description", r.Description); err != nil {
		return err
	}
	if err := encode("content:encoded", r.Content); err != nil {
		return err
	}
	if err := encode("pubDate", r.PubDate); err != nil {
		return err
	}
	if err := encode("guid", r.GUID); err != nil {
		return err
	}
	if err := encode("deframer:custom", r.CustomField); err != nil {
		return err
	}

	if err := encodeExtensions(e, r.Extensions); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func encodeExtensions(e *xml.Encoder, extensions ext.Extensions) error {
	var prefixes []string
	for prefix := range extensions {
		prefixes = append(prefixes, prefix)
	}
	sort.Strings(prefixes)

	for _, prefix := range prefixes {
		elements := extensions[prefix]
		var names []string
		for name := range elements {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			exts := elements[name]
			for _, ext := range exts {
				local := name
				if prefix != "" {
					local = prefix + ":" + name
				}
				start := xml.StartElement{Name: xml.Name{Local: local}}
				var attrs []xml.Attr
				for k, v := range ext.Attrs {
					attrs = append(attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
				}
				sort.Slice(attrs, func(i, j int) bool {
					return attrs[i].Name.Local < attrs[j].Name.Local
				})
				start.Attr = append(start.Attr, attrs...)

				if err := e.EncodeToken(start); err != nil {
					return err
				}

				if len(ext.Children) > 0 {
					if err := encodeExtensionChildren(e, ext.Children); err != nil {
						return err
					}
				} else {
					if err := e.EncodeToken(xml.CharData(ext.Value)); err != nil {
						return err
					}
				}

				if err := e.EncodeToken(start.End()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func encodeExtensionChildren(e *xml.Encoder, children map[string][]ext.Extension) error {
	var names []string
	for name := range children {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		exts := children[name]
		for _, ext := range exts {
			start := xml.StartElement{Name: xml.Name{Local: name}}
			var attrs []xml.Attr
			for k, v := range ext.Attrs {
				attrs = append(attrs, xml.Attr{Name: xml.Name{Local: k}, Value: v})
			}
			sort.Slice(attrs, func(i, j int) bool {
				return attrs[i].Name.Local < attrs[j].Name.Local
			})
			start.Attr = append(start.Attr, attrs...)

			if err := e.EncodeToken(start); err != nil {
				return err
			}

			if len(ext.Children) > 0 {
				if err := encodeExtensionChildren(e, ext.Children); err != nil {
					return err
				}
			} else {
				if err := e.EncodeToken(xml.CharData(ext.Value)); err != nil {
					return err
				}
			}

			if err := e.EncodeToken(start.End()); err != nil {
				return err
			}
		}
	}
	return nil
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
