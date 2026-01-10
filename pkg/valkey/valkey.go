package valkey

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/egandro/news-deframer/pkg/analyzer"
	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	driver "github.com/valkey-io/valkey-go"
)

const (
	key_prefix_feed_uuid    = "feed_uuid:"
	key_prefix_feed_reverse = "feed_reverse:"
)

type Cache int

const (
	Ok Cache = iota
	ValueUnknown
	Updating
)

type FeedUUIDCache struct {
	Cache      Cache     `json:"cache"`
	UUID       uuid.UUID `json:"uuid,omitempty"`
	BaseDomain []string  `json:"base_domain,omitempty"`
}

type ItemHashCache struct {
	Cache          Cache                   `json:"cache"`
	Hash           string                  `json:"hash,omitempty"`
	Content        string                  `json:"content,omitempty"`
	AnalyzerResult analyzer.AnalyzerResult `json:"analyzer_result,omitempty"`
}

type Valkey interface {
	DrainFeed(id uuid.UUID) error
	GetFeedUUID(u *url.URL) (*FeedUUIDCache, error)
	UpdateFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) error
	TryLockFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) (bool, error)
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

func feedUUIDKey(u *url.URL) string {
	return key_prefix_feed_uuid + url.QueryEscape(u.String())
}

func feedReverseKey(id uuid.UUID) string {
	return key_prefix_feed_reverse + id.String()
}

func (v *valkey) GetFeedUUID(u *url.URL) (*FeedUUIDCache, error) {
	val, err := v.client.Do(v.ctx, v.client.B().Get().Key(feedUUIDKey(u)).Build()).AsBytes()
	if driver.IsValkeyNil(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var state FeedUUIDCache
	if err := json.Unmarshal(val, &state); err != nil {
		return nil, fmt.Errorf("invalid json in cache: %w", err)
	}

	return &state, nil
}

func (v *valkey) UpdateFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	key := feedUUIDKey(u)

	// Pipeline: Set the cache key AND update reverse index if it's a valid feed
	cmds := make([]driver.Completed, 0, 2)
	cmds = append(cmds, v.client.B().Set().Key(key).Value(driver.BinaryString(data)).Ex(ttl).Build())

	// Only add to reverse index if we definitively know the UUID
	if state.Cache == Ok && state.UUID != uuid.Nil {
		// We store the encoded URL as the member in the Set
		cmds = append(cmds, v.client.B().Sadd().Key(feedReverseKey(state.UUID)).Member(url.QueryEscape(u.String())).Build())
	}

	// Execute
	for _, cmd := range cmds {
		if err := v.client.Do(v.ctx, cmd).Error(); err != nil {
			return err
		}
	}
	return nil
}

func (v *valkey) TryLockFeedUUID(u *url.URL, state FeedUUIDCache, ttl time.Duration) (bool, error) {
	// TryLockFeedUUID uses Nx() (Set if Not Exists) to ensure atomic locking.
	data, err := json.Marshal(state)
	if err != nil {
		return false, fmt.Errorf("failed to marshal state: %w", err)
	}

	err = v.client.Do(v.ctx, v.client.B().Set().Key(feedUUIDKey(u)).Value(driver.BinaryString(data)).Nx().Ex(ttl).Build()).Error()
	if driver.IsValkeyNil(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	// It returns true if the lock was acquired (key didn't exist), false otherwise.
	return true, nil
}

func (v *valkey) DrainFeed(id uuid.UUID) error {
	revKey := feedReverseKey(id)

	// 1. Get all URLs associated with this Feed ID
	members, err := v.client.Do(v.ctx, v.client.B().Smembers().Key(revKey).Build()).AsStrSlice()
	if err != nil {
		return err
	}

	if len(members) == 0 {
		return nil
	}

	// 2. Build pipeline to delete all individual URL keys
	// Note: members are already query escaped strings from UpdateFeedUUID logic
	cmds := make([]driver.Completed, 0, len(members)+1)

	for _, member := range members {
		// Reconstruct the full key: "feed_uuid:" + escaped_url
		fullKey := key_prefix_feed_uuid + member
		cmds = append(cmds, v.client.B().Del().Key(fullKey).Build())
	}

	// 3. Delete the reverse index itself
	cmds = append(cmds, v.client.B().Del().Key(revKey).Build())

	// 4. Execute batch
	for _, cmd := range cmds {
		if err := v.client.Do(v.ctx, cmd).Error(); err != nil {
			slog.Error("failed to execute drain command", "err", err)
			// continue processing to try best-effort cleanup
		}
	}

	return nil
}

func (v *valkey) Close() error {
	v.client.Close()
	return nil
}
