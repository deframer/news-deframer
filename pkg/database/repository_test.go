package database

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/google/uuid"
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

func TestFindCachedFeedById(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	tests := []struct {
		name          string
		setup         func() uuid.UUID
		expectedFound bool
	}{
		{
			name: "Found (Enabled)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: true,
		},
		{
			name: "NotFound (Disabled)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (Deleted)",
			setup: func() uuid.UUID {
				feed := Feed{
					URL:     "http://test-cached-feed-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)

				assert.NoError(t, db.Delete(&feed).Error)
				return feed.ID
			},
			expectedFound: false,
		},
		{
			name: "NotFound (NonExistent)",
			setup: func() uuid.UUID {
				return uuid.New()
			},
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.setup()
			found, err := repo.FindCachedFeedById(id)
			assert.NoError(t, err)
			if tt.expectedFound {
				assert.NotNil(t, found)
				assert.Equal(t, id, found.ID)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}

func TestFindItemsByCachedFeedFeedId(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)
	db := repo.(*repository).db

	// Helper to create valid SHA256 hash string
	makeHash := func(s string) string {
		h := sha256.Sum256([]byte(s))
		return hex.EncodeToString(h[:])
	}

	tests := []struct {
		name          string
		setup         func() (uuid.UUID, string)
		expectedCount int
	}{
		{
			name: "Found (Enabled)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-enabled-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-enabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 1,
		},
		{
			name: "Empty (Disabled)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-disabled-" + uuid.New().String(),
					Enabled: false,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-disabled-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (Deleted)",
			setup: func() (uuid.UUID, string) {
				feed := Feed{
					URL:     "http://test-items-deleted-" + uuid.New().String(),
					Enabled: true,
				}
				assert.NoError(t, db.Create(&feed).Error)

				hash := makeHash("item-deleted-" + feed.ID.String())
				item := Item{Hash: hash, FeedID: feed.ID, URL: "http://item", AIResult: JSONB{"a": 1}}
				assert.NoError(t, db.Create(&item).Error)

				cachedFeed := CachedFeed{
					ID:        feed.ID,
					XMLHeader: "<rss></rss>",
					ItemRefs:  StringArray{hash},
				}
				assert.NoError(t, db.Create(&cachedFeed).Error)

				assert.NoError(t, db.Delete(&feed).Error)
				return feed.ID, hash
			},
			expectedCount: 0,
		},
		{
			name: "Empty (NonExistent)",
			setup: func() (uuid.UUID, string) {
				return uuid.New(), ""
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, hash := tt.setup()
			items, err := repo.FindItemsByCachedFeedFeedId(id)
			assert.NoError(t, err)
			assert.Len(t, items, tt.expectedCount)
			if tt.expectedCount > 0 {
				assert.Equal(t, hash, items[0].Hash)
			}
		})
	}
}
