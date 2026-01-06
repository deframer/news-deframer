package feeds

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/mmcdole/gofeed"
)

type Feeds interface {
	ParseFeed(ctx context.Context, content io.Reader) (*gofeed.Feed, error)
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
