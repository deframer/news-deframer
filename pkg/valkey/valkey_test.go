package valkey

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFeedUrlKeys(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	ctx := context.Background()
	c, err := New(ctx, cfg)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()

	rawURL := "http://example.com/feed"
	u, err := url.Parse(rawURL)
	assert.NoError(t, err)
	value := uuid.New()

	// Ensure cleanup before and after
	if existing, _ := c.GetFeedByUrl(u); existing != nil && existing.UUID != uuid.Nil {
		_ = c.DrainFeed(existing.UUID)
	}
	defer func() {
		_ = c.DrainFeed(value)
	}()

	foundVal, err := c.GetFeedByUrl(u)
	assert.NoError(t, err)
	assert.Nil(t, foundVal)

	err = c.UpdateFeedByUrl(FeedUrlToUUID{
		Cache: Ok,
		UUID:  value,
	}, FeedInfo{BaseDomain: []string{"example.com", "sub.example.com"}, URL: u.String()}, time.Minute)
	assert.NoError(t, err)

	foundVal, err = c.GetFeedByUrl(u)
	assert.NoError(t, err)
	assert.NotNil(t, foundVal)
	assert.Equal(t, Ok, foundVal.Cache)
	assert.Equal(t, value, foundVal.UUID)

	err = c.DrainFeed(value)
	assert.NoError(t, err)

	foundVal, err = c.GetFeedByUrl(u)
	assert.NoError(t, err)
	assert.Nil(t, foundVal)
}
