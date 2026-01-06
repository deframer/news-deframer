package facade

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasFeed(t *testing.T) {
	ctx := context.Background()

	// Dependencies are nil as the current implementation is a stub
	f := New(ctx, nil, nil, nil)

	exists, err := f.HasFeed(ctx, "http://example.com/feed.xml")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestHasArticle(t *testing.T) {
	ctx := context.Background()

	// Dependencies are nil as the current implementation is a stub
	f := New(ctx, nil, nil, nil)

	exists, err := f.HasArticle(ctx, "http://example.com/article")
	assert.NoError(t, err)
	assert.False(t, exists)
}
