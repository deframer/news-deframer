package service

import (
	"context"
	"fmt"
	"net/url"

	rss "github.com/deframer/news-deframer/gen/rss"
	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	"goa.design/clue/log"
)

type rsssrvc struct {
	facade facade.Facade
}

// NewRss returns the RSS service implementation.
func NewRss(ctx context.Context) rss.Service {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	repo, err := database.NewRepository(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return &rsssrvc{facade: facade.New(ctx, cfg, repo)}
}

// Return a proxied RSS feed.
func (s *rsssrvc) Feed(ctx context.Context, p *rss.RSSPayload) (res string, err error) {
	log.Printf(ctx, "handleRSS url=%s", p.URL)
	if p.URL == "" {
		return "", rss.BadRequest("missing url")
	}

	u, err := url.ParseRequestURI(p.URL)
	if err != nil {
		log.Errorf(ctx, err, "invalid url")
		return "", rss.BadRequest("invalid url")
	}

	filter := facade.RSSProxyFilter{
		URL:      u.String(),
		Lang:     p.Lang,
		Max:      p.MaxScore,
		Embedded: p.Embedded,
	}

	feed, err := s.facade.GetRssProxyFeed(ctx, &filter)
	if err != nil {
		log.Errorf(ctx, err, "failed to get rss proxy feed")
		return "", fmt.Errorf("failed to get rss proxy feed")
	}

	return feed, nil
}
