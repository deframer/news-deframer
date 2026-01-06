package facade

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasFeed(t *testing.T) {
	ctx := context.Background()

	// Dependencies are nil as the current implementation is a stub
	f := New(ctx, nil, nil, nil)

	u, err := url.Parse("http://example.com/feed.xml")
	assert.NoError(t, err)
	exists, err := f.HasFeed(ctx, u)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestHasArticle(t *testing.T) {
	ctx := context.Background()

	// Dependencies are nil as the current implementation is a stub
	f := New(ctx, nil, nil, nil)

	u, err := url.Parse("http://example.com/article")
	assert.NoError(t, err)
	exists, err := f.HasArticle(ctx, u)
	assert.NoError(t, err)
	assert.False(t, exists)
}
