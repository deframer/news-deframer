package valkey

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	driver "github.com/valkey-io/valkey-go"
)

type Cache int

const (
	Ok Cache = iota
	ValueUnknown
	Updating
)

type FeedUUIDCache struct {
	Cache Cache
	UUID  uuid.UUID
}

type Valkey interface {
	GetFeedUUID(u *url.URL) (*FeedUUIDCache, error)
	UpdateFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) error
	TryLockFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) (bool, error)
	DeleteFeedUUID(u *url.URL) error
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

const key_prefix_feed_uuid = "feed_uuid:"

func feedUUIDKey(u *url.URL) string {
	return key_prefix_feed_uuid + url.QueryEscape(u.String())
}

func (v *valkey) GetFeedUUID(u *url.URL) (*FeedUUIDCache, error) {
	val, err := v.client.Do(v.ctx, v.client.B().Get().Key(feedUUIDKey(u)).Build()).ToString()
	if driver.IsValkeyNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	state := &FeedUUIDCache{}
	if i, err := strconv.Atoi(val); err == nil {
		state.Cache = Cache(i)
		state.UUID = uuid.Nil
		return state, nil
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid in cache %q: %w", val, err)
	}
	state.Cache = Ok
	state.UUID = id
	return state, nil
}

func (v *valkey) UpdateFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) error {
	val := state.UUID.String()
	if state.Cache != Ok {
		val = strconv.Itoa(int(state.Cache))
	}
	return v.client.Do(v.ctx, v.client.B().Set().Key(feedUUIDKey(u)).Value(val).Ex(ttl).Build()).Error()
}

func (v *valkey) TryLockFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) (bool, error) {
	// TryLockFeedUUID uses Nx() (Set if Not Exists) to ensure atomic locking.
	// It returns true if the lock was acquired (key didn't exist), false otherwise.
	val := state.UUID.String()
	if state.Cache != Ok {
		val = strconv.Itoa(int(state.Cache))
	}
	err := v.client.Do(v.ctx, v.client.B().Set().Key(feedUUIDKey(u)).Value(val).Nx().Ex(ttl).Build()).Error()
	if driver.IsValkeyNil(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (v *valkey) DeleteFeedUUID(u *url.URL) error {
	return v.client.Do(v.ctx, v.client.B().Del().Key(feedUUIDKey(u)).Build()).Error()
}

func (v *valkey) Close() error {
	v.client.Close()
	return nil
}
