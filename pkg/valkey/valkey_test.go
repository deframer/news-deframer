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

	ctx := context.Background()
	c, err := New(ctx, cfg)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()

	key := uuid.New()

	// Ensure cleanup before and after
	_ = c.DeleteFeed(key)
	defer func() {
		_ = c.DeleteFeed(key)
	}()

	err = c.AddFeed(key)
	assert.NoError(t, err)

	exists, err := c.HasFeed(key)
	assert.NoError(t, err)
	assert.True(t, exists)

	foundKey, err := c.GetFeed(key)
	assert.NoError(t, err)
	assert.NotNil(t, foundKey)
	assert.Equal(t, key, *foundKey)

	err = c.DeleteFeed(key)
	assert.NoError(t, err)

	exists, err = c.HasFeed(key)
	assert.NoError(t, err)
	assert.False(t, exists)

	foundKey, err = c.GetFeed(key)
	assert.NoError(t, err)
	assert.Nil(t, foundKey)
}
