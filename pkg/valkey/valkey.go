package valkey

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	driver "github.com/valkey-io/valkey-go"
)

type Valkey interface {
	HasFeed(ctx context.Context, key uuid.UUID) (bool, error)
	GetFeed(ctx context.Context, key uuid.UUID) (*uuid.UUID, error)
	AddFeed(ctx context.Context, key uuid.UUID) error
	DeleteFeed(ctx context.Context, key uuid.UUID) error
	Close() error
}

type valkey struct {
	client driver.Client
}

func New(cfg *config.Config) (Valkey, error) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Do(ctx, client.B().Ping().Build()).Error(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to connect to valkey: %w", err)
	}

	return &valkey{client: client}, nil
}

const prefix = "feed:"

func (v *valkey) HasFeed(ctx context.Context, key uuid.UUID) (bool, error) {
	n, err := v.client.Do(ctx, v.client.B().Exists().Key(prefix+key.String()).Build()).ToInt64()
	return n > 0, err
}

func (v *valkey) GetFeed(ctx context.Context, key uuid.UUID) (*uuid.UUID, error) {
	_, err := v.client.Do(ctx, v.client.B().Get().Key(prefix+key.String()).Build()).ToString()
	if driver.IsValkeyNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (v *valkey) AddFeed(ctx context.Context, key uuid.UUID) error {
	return v.client.Do(ctx, v.client.B().Set().Key(prefix+key.String()).Value(time.Now().Format(time.RFC3339)).Build()).Error()
}

func (v *valkey) DeleteFeed(ctx context.Context, key uuid.UUID) error {
	return v.client.Do(ctx, v.client.B().Del().Key(prefix+key.String()).Build()).Error()
}

func (v *valkey) Close() error {
	v.client.Close()
	return nil
}
