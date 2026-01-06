package facade

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/valkey"
)

type Facade struct {
	ctx    context.Context
	cfg    *config.Config
	valkey valkey.Valkey
	repo   database.Repository
}

func New(ctx context.Context, cfg *config.Config, v valkey.Valkey, repo database.Repository) *Facade {
	return &Facade{
		ctx:    ctx,
		cfg:    cfg,
		valkey: v,
		repo:   repo,
	}
}

var (
	maxPendingTimeout = 30 * time.Second
	checkInterval     = 5 * time.Second
	maxTimeout        = 5 * time.Minute
)

func (f *Facade) HasFeed(ctx context.Context, u *url.URL) (bool, error) {
	val, err := f.valkey.GetFeedUrl(u)
	if err != nil {
		return false, err
	}

	if val != nil {
		switch *val {
		case "valid":
			return true, nil
		case "invalid":
			return false, nil
		case "pending":
			ticker := time.NewTicker(checkInterval)
			defer ticker.Stop()
			timeout := time.After(maxPendingTimeout)
			for {
				select {
				case <-timeout:
					return false, fmt.Errorf("timeout waiting for pending feed")
				case <-ticker.C:
					val, err := f.valkey.GetFeedUrl(u)
					if err != nil {
						return false, err
					}
					if val != nil && *val != "pending" {
						return *val == "valid", nil
					}
				}
			}
		}
	}

	if err := f.valkey.AddFeedUrl(u, "pending", maxPendingTimeout); err != nil {
		return false, err
	}

	feed, err := f.repo.FindFeedByUrl(u)
	if err != nil {
		return false, err
	}

	status := "invalid"
	if feed != nil {
		status = "valid"
	}

	if err := f.valkey.AddFeedUrl(u, status, maxTimeout); err != nil {
		return false, err
	}

	return status == "valid", nil
}

func (f *Facade) HasArticle(ctx context.Context, u *url.URL) (bool, error) {
	return false, nil
}
