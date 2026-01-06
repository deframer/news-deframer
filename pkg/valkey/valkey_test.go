package valkey

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
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
	value := "some-value"

	// Ensure cleanup before and after
	_ = c.DeleteFeedUrl(u)
	defer func() {
		_ = c.DeleteFeedUrl(u)
	}()

	foundVal, err := c.GetFeedUrl(u)
	assert.NoError(t, err)
	assert.Nil(t, foundVal)

	err = c.AddFeedUrl(u, value, time.Minute)
	assert.NoError(t, err)

	foundVal, err = c.GetFeedUrl(u)
	assert.NoError(t, err)
	assert.NotNil(t, foundVal)
	assert.Equal(t, value, *foundVal)

	err = c.DeleteFeedUrl(u)
	assert.NoError(t, err)

	foundVal, err = c.GetFeedUrl(u)
	assert.NoError(t, err)
	assert.Nil(t, foundVal)
}
