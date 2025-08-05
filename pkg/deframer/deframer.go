package deframer

import (
	"fmt"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/source"
)

type deframer struct {
	db         *database.Database
	ai         openai.OpenAI
	src        *source.Source
	downloader downloader.Downloader
	prompts    map[string]source.Prompt
}

type Deframer interface {
	UpdateFeeds() error
}

// NewDeframer initializes a new deframer
func NewDeframer() (Deframer, error) {
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
		db:         db,
		ai:         ai,
		src:        src,
		downloader: downloader,
		prompts:    prompts,
	}

	return res, nil
}

func (d *deframer) UpdateFeeds() error {
	for _, feed := range d.src.Feeds {
		data, err := d.downloader.DownloadRSSFeed(feed.RSS_URL)
		if err != nil {
			return err
		}
		fmt.Printf("data: %v", data)
	}
	return nil
}
