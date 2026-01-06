package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/valkey-io/valkey-go"
)

type Cache interface {
	HasFeed(ctx context.Context, key uuid.UUID) (bool, error)
	GetFeed(ctx context.Context, key uuid.UUID) (*uuid.UUID, error)
	AddFeed(ctx context.Context, key uuid.UUID) error
	DeleteFeed(ctx context.Context, key uuid.UUID) error
	Close() error
}

type cache struct {
	client valkey.Client
}

func New(cfg *config.Config) (Cache, error) {
	db, err := strconv.Atoi(cfg.ValkeyDB)
	if err != nil {
		return nil, fmt.Errorf("invalid valkey db: %w", err)
	}

	opt := valkey.ClientOption{
		InitAddress: []string{cfg.ValkeyHost},
		Password:    cfg.ValkeyPassword,
		SelectDB:    db,
	}

	client, err := valkey.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	return &cache{client: client}, nil
}

const prefix = "feed:"

func (c *cache) HasFeed(ctx context.Context, key uuid.UUID) (bool, error) {
	n, err := c.client.Do(ctx, c.client.B().Exists().Key(prefix+key.String()).Build()).ToInt64()
	return n > 0, err
}

func (c *cache) GetFeed(ctx context.Context, key uuid.UUID) (*uuid.UUID, error) {
	_, err := c.client.Do(ctx, c.client.B().Get().Key(prefix+key.String()).Build()).ToString()
	if valkey.IsValkeyNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (c *cache) AddFeed(ctx context.Context, key uuid.UUID) error {
	return c.client.Do(ctx, c.client.B().Set().Key(prefix+key.String()).Value(time.Now().Format(time.RFC3339)).Build()).Error()
}

func (c *cache) DeleteFeed(ctx context.Context, key uuid.UUID) error {
	return c.client.Do(ctx, c.client.B().Del().Key(prefix+key.String()).Build()).Error()
}

func (c *cache) Close() error {
	c.client.Close()
	return nil
}
