package valkey

import (
	"context"
	"testing"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFeedKeys(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	c, err := New(cfg)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()

	ctx := context.Background()
	key := uuid.New()

	// Ensure cleanup before and after
	_ = c.DeleteFeed(ctx, key)
	defer func() {
		_ = c.DeleteFeed(ctx, key)
	}()

	err = c.AddFeed(ctx, key)
	assert.NoError(t, err)

	exists, err := c.HasFeed(ctx, key)
	assert.NoError(t, err)
	assert.True(t, exists)

	foundKey, err := c.GetFeed(ctx, key)
	assert.NoError(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, key, *foundKey)

	err = c.DeleteFeed(ctx, key)
	assert.NoError(t, err)

	exists, err = c.HasFeed(ctx, key)
	assert.NoError(t, err)
	assert.False(t, exists)

	foundKey, err = c.GetFeed(ctx, key)
	assert.NoError(t, err)
	assert.Nil(t, foundKey)
}
