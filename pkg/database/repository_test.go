package database

import (
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetTime(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)

	got, err := repo.GetTime()
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), got, 1*time.Minute)
}

func TestFindFeedByUrl(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)

	t.Run("Found", func(t *testing.T) {
		// This URL is seeded in pkg/database/migrate.go
		rawURL := "http://dummy"
		u, err := url.Parse(rawURL)
		assert.NoError(t, err)

		feed, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		assert.NotNil(t, feed)
		assert.Equal(t, rawURL, feed.URL)
	})

	t.Run("NotFound", func(t *testing.T) {
		u, err := url.Parse("http://does-not-exist.com/feed")
		assert.NoError(t, err)

		feed, err := repo.FindFeedByUrl(u)
		assert.NoError(t, err)
		assert.Nil(t, feed)
	})
}
