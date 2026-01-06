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

const (
	statusValid   = "valid"
	statusInvalid = "invalid"
	statusPending = "pending"
)

func (f *Facade) HasFeed(ctx context.Context, u *url.URL) (bool, error) {
	val, err := f.valkey.GetFeedUrl(u)
	if err != nil {
		return false, err
	}

	// If we have a definitive result, return it. (= cache hit)
	if val != nil && *val != statusPending {
		return *val == statusValid, nil
	}

	// If it was missing (nil), try to acquire the lock to be the one fetching it.
	if val == nil {
		acquired, err := f.valkey.TryLockFeedUrl(u, statusPending, maxPendingTimeout)
		if err != nil {
			return false, err
		}
		if !acquired {
			// We failed to acquire lock, meaning someone else just set it to pending.
			// Fall through to the wait loop.
			goto WAIT_LOOP
		}

		// We acquired the lock. Proceed to DB lookup.
		return f.fetchAndCache(u)
	}

	// It is pending (either from GetFeedUrl or failed UpdateFeedUrl). Wait for it.
WAIT_LOOP:
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	timeout := time.After(maxPendingTimeout)

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timeout:
			return false, fmt.Errorf("timeout waiting for pending feed")
		case <-ticker.C:
			val, err := f.valkey.GetFeedUrl(u)
			if err != nil {
				return false, err
			}
			if val != nil && *val != statusPending {
				return *val == statusValid, nil
			}
			// If val is nil here, it means the pending key expired or was deleted without result.
			// we don't care and wait for the next query
		}
	}
}

func (f *Facade) fetchAndCache(u *url.URL) (bool, error) {
	feed, err := f.repo.FindFeedByUrl(u)
	if err != nil {
		return false, err
	}

	status := statusInvalid
	if feed != nil {
		status = statusValid
	}

	if err := f.valkey.UpdateFeedUrl(u, status, maxTimeout); err != nil {
		return false, err
	}

	return status == statusValid, nil
}

func (f *Facade) HasArticle(ctx context.Context, u *url.URL) (bool, error) {
	return false, nil
}
