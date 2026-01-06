package valkey

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	driver "github.com/valkey-io/valkey-go"
)

type Valkey interface {
	GetFeedUrl(u *url.URL) (*string, error)
	AddFeedUrl(u *url.URL, value string, ttl time.Duration) error
	DeleteFeedUrl(u *url.URL) error
	Close() error
}

type valkey struct {
	client driver.Client
	ctx    context.Context
}

func New(ctx context.Context, cfg *config.Config) (Valkey, error) {
	db, err := strconv.Atoi(cfg.ValkeyDB)
	if err != nil {
		return nil, fmt.Errorf("invalid valkey db: %w", err)
	}

	opt := driver.ClientOption{
		InitAddress: []string{cfg.ValkeyHost},
		Password:    cfg.ValkeyPassword,
		SelectDB:    db,
	}

	client, err := driver.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	// Verify connection
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Do(pingCtx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	return &valkey{client: client, ctx: ctx}, nil
}

const prefix = "feed:"

func urlToValkeyKey(u *url.URL) string {
	return prefix + url.QueryEscape(u.String())
}

func (v *valkey) GetFeedUrl(u *url.URL) (*string, error) {
	val, err := v.client.Do(v.ctx, v.client.B().Get().Key(urlToValkeyKey(u)).Build()).ToString()
	if driver.IsValkeyNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func (v *valkey) AddFeedUrl(u *url.URL, value string, ttl time.Duration) error {
	return v.client.Do(v.ctx, v.client.B().Set().Key(urlToValkeyKey(u)).Value(value).Ex(ttl).Build()).Error()
}

func (v *valkey) DeleteFeedUrl(u *url.URL) error {
	return v.client.Do(v.ctx, v.client.B().Del().Key(urlToValkeyKey(u)).Build()).Error()
}

func (v *valkey) Close() error {
	v.client.Close()
	return nil
}
