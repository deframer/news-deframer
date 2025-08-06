package deframer

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/source"
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

const maxAge = time.Minute * 90

type deframer struct {
	ctx        context.Context
	db         *database.Database
	ai         openai.OpenAI
	src        *source.Source
	downloader downloader.Downloader
	prompts    map[string]source.Prompt
}

type Deframer interface {
	UpdateFeeds() (int, error)
	DeframeFeed(parsedData *gofeed.Feed, feed source.Feed) (string, error)
	DeframeItem(item *gofeed.Item, feed source.Feed) (*gofeed.Item, error)
}

// NewDeframer initializes a new deframer
func NewDeframer(ctx context.Context) (Deframer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(cfg.DatabaseFile)
	if err != nil {
		return nil, err
	}

	ai := openai.NewAI(cfg.AI_URL, cfg.AI_Model, "")

	src, err := source.ParseFile(cfg.Source)
	if err != nil {
		return nil, err
	}

	downloader := downloader.NewDownloader()

	prompts := make(map[string]source.Prompt)
	for _, prompt := range src.Prompts {
		prompts[prompt.Language] = prompt
	}

	res := &deframer{
		ctx:        ctx,
		db:         db,
		ai:         ai,
		src:        src,
		downloader: downloader,
		prompts:    prompts,
	}

	return res, nil
}

func (d *deframer) UpdateFeeds() (int, error) {
	numberOfDownloads := 0
	parser := gofeed.NewParser()

	for _, feed := range d.src.Feeds {
		cache, err := d.db.FindCacheByFeedUrl(feed.RSS_URL, maxAge)
		if err != nil {
			return numberOfDownloads, err
		}

		if cache != nil {
			continue
		}

		data, err := d.downloader.DownloadRSSFeed(feed.RSS_URL)
		if err != nil {
			return numberOfDownloads, err
		}

		parsedData, err := parser.ParseString(string(data))
		if err != nil {
			return numberOfDownloads, err
		}

		title := parsedData.Title
		if title == "" {
			// some fallback
			title = feed.RSS_URL
		}

		title = fmt.Sprintf("%v (%v)", title, feed.Language)

		unframed, err := d.DeframeFeed(parsedData, feed)
		if err != nil {
			return numberOfDownloads, err
		}

		cache = &database.Cache{
			FeedUrl: feed.RSS_URL,
			Title:   title,
			Cache:   unframed,
		}

		err = d.db.CreateCache(cache)
		if err != nil {
			return numberOfDownloads, err
		}

		numberOfDownloads++
	}
	return numberOfDownloads, nil
}

func (d *deframer) DeframeFeed(parsedData *gofeed.Feed, feed source.Feed) (string, error) {
	// Update channel title with prefix
	prefix := "[Prefix] "
	parsedData.Title = prefix + parsedData.Title

	newFeed := &feeds.Feed{
		Title: parsedData.Title,
		Link: &feeds.Link{
			Href: parsedData.Link,
			Type: parsedData.FeedType,
		},
		Description: parsedData.Description,
		Author:      &feeds.Author{},
		// // Id:
		// // Subtitle:
		// // Items:
		Copyright: parsedData.Copyright,
		// Image:
	}

	if parsedData.Author != nil {
		newFeed.Author.Name = parsedData.Author.Name
		newFeed.Author.Email = parsedData.Author.Email
	}

	if parsedData.PublishedParsed != nil {
		newFeed.Created = *parsedData.PublishedParsed
	}

	if parsedData.UpdatedParsed != nil {
		newFeed.Updated = *parsedData.UpdatedParsed
	}

	if parsedData.Image != nil {
		newFeed.Image = &feeds.Image{}
		newFeed.Image.Url = parsedData.Image.URL
		newFeed.Image.Title = parsedData.Image.Title
	}

	// TODO: add a dummy feed

	// item := &feeds.Item{
	// 	Title:       "Ihr Post Titel",
	// 	Link:        &feeds.Link{Href: "http://example.com/post-url"},
	// 	Description: "Eine kurze Beschreibung zu Ihrem Post",
	// 	Author:      &feeds.Author{Name: "Your Name", Email: "yourname@example.com"},
	// 	Created:     time.Now(),
	// 	Id:          "my id",
	// }
	// newFeed.Add(item)

	for _, current := range parsedData.Items {
		deframed, err := d.DeframeItem(current, feed)
		if err != nil {
			// maybe just continue?
			return "", err
		}

		item := &feeds.Item{
			Title:       deframed.Title,
			Link:        &feeds.Link{Href: deframed.Link},
			Description: deframed.Description,
			Content:     deframed.Content,
			Id:          deframed.GUID,
		}

		if deframed.PublishedParsed != nil {
			item.Created = *deframed.PublishedParsed
		}

		if deframed.UpdatedParsed != nil {
			item.Updated = *deframed.UpdatedParsed
		}

		if len(deframed.Authors) > 0 {
			item.Author = &feeds.Author{
				Name:  deframed.Authors[0].Name,
				Email: deframed.Authors[0].Email,
			}
		}

		newFeed.Add(item)
	}

	result, err := newFeed.ToRss()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (d *deframer) DeframeItem(item *gofeed.Item, feed source.Feed) (*gofeed.Item, error) {
	key := fmt.Sprintf("%v-%v", feed.RSS_URL, item.GUID)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(key)))

	dbItem, err := d.db.FindItemByHash(hash)
	if err != nil {
		return nil, err
	}

	if dbItem != nil {
		item.Link = dbItem.Link
		item.GUID = dbItem.Guid
		item.Title = dbItem.Title
		item.Description = dbItem.Description
		item.Content = dbItem.Content
		return item, nil
	}

	dbItem, err = d.deframeItemInternal(item, feed)
	if err != nil {
		return nil, err
	}

	err = d.db.CreateItem(dbItem)
	if err != nil {
		return nil, err
	}

	item.Link = dbItem.Link
	item.GUID = dbItem.Guid
	item.Title = dbItem.Title
	item.Description = dbItem.Description
	item.Content = dbItem.Content

	return item, err
}

func (d *deframer) deframeItemInternal(item *gofeed.Item, feed source.Feed) (*database.Item, error) {
	res := &database.Item{
		Link:        item.Link,
		Guid:        item.GUID,
		Title:       item.Title,
		Description: item.Description,
		Content:     item.Content,
	}

	if _, ok := d.prompts[feed.Language]; !ok {
		// we don't know this language
		return res, nil
	}

	prompt := d.prompts[feed.Language]
	user := prompt.User
	system := prompt.System

	user = strings.ReplaceAll(user, "$TITLE", item.Title)
	user = strings.ReplaceAll(user, "$DESCRIPTION", item.Description)

	system = strings.ReplaceAll(system, "$TITLE", item.Title)
	system = strings.ReplaceAll(system, "$DESCRIPTION", item.Description)

	resultString, err := d.ai.Query(d.ctx, user, system)
	if err != nil {
		return nil, err
	}

	resultAny, err := d.ai.FuzzyParseJSON(resultString)
	if err != nil {
		return nil, err
	}

	title_corrected := ""
	reason := ""
	var framing float64

	if resultMap, ok := resultAny.(map[string]any); ok {
		if v, ok := resultMap["title_corrected"]; ok {
			if d, ok := v.(string); ok {
				title_corrected = d
			}
		}
		if v, ok := resultMap["framing"]; ok {
			if d, ok := v.(json.Number); ok {
				f, err := d.Float64()
				if err == nil {
					framing = f
				}
			}
		}
		if v, ok := resultMap["reason"]; ok {
			if d, ok := v.(string); ok {
				reason = d
			}
		}
	}

	if title_corrected != "" && framing > 0.0 {
		res.TitleAI = &title_corrected
		res.Framing = &framing
		res.ReasonAI = &reason
		title := fmt.Sprintf("Framing: %v - %v", framing, title_corrected)
		if res.Content != "" {
			res.Content = fmt.Sprintf("Original title: %v <br/> Reason: %v <br/> %v", res.Title, reason, res.Content)
		}
		res.Title = title
	}

	return res, nil
}
