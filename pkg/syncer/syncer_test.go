package syncer

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	enqueueSyncCalled bool
	lastId            uuid.UUID
	removeSyncCalled  bool
	upsertItemFunc    func(item *database.Item) error
}

// Implement database.Repository interface stubs
func (m *mockRepo) FindFeedByUrl(u *url.URL) (*database.Feed, error) { return nil, nil }
func (m *mockRepo) FindFeedByUrlAndAvailability(u *url.URL, onlyEnabled bool) (*database.Feed, error) {
	return nil, nil
}
func (m *mockRepo) FindFeedById(feedID uuid.UUID) (*database.Feed, error) { return nil, nil }
func (m *mockRepo) UpsertFeed(feed *database.Feed) error                  { return nil }
func (m *mockRepo) FindItemsByUrl(u *url.URL) ([]database.Item, error)    { return nil, nil }
func (m *mockRepo) GetAllFeeds(deleted bool) ([]database.Feed, error)     { return nil, nil }
func (m *mockRepo) DeleteFeedById(id uuid.UUID) error                     { return nil }
func (m *mockRepo) RemoveSync(id uuid.UUID) error {
	m.removeSyncCalled = true
	m.lastId = id
	return nil
}
func (m *mockRepo) BeginFeedUpdate(lockDuration time.Duration) (*database.Feed, error) {
	return nil, nil
}
func (m *mockRepo) EndFeedUpdate(id uuid.UUID, err error, successDelay time.Duration) error {
	return nil
}
func (m *mockRepo) GetPendingItems(feedID uuid.UUID, hashes []string, maxRetries int) (map[string]int, error) {
	res := make(map[string]int, len(hashes))
	for _, h := range hashes {
		res[h] = 0
	}
	return res, nil
}
func (m *mockRepo) UpsertItem(item *database.Item) error {
	if m.upsertItemFunc != nil {
		return m.upsertItemFunc(item)
	}
	return nil
}
func (m *mockRepo) GetItemsByHashes(feedID uuid.UUID, hashes []string) ([]database.Item, error) {
	return nil, nil
}
func (m *mockRepo) UpsertCachedFeed(cachedFeed *database.CachedFeed) error { return nil }
func (m *mockRepo) FindCachedFeedById(feedID uuid.UUID) (*database.CachedFeed, error) {
	return nil, nil
}
func (m *mockRepo) FindFeedScheduleById(feedID uuid.UUID) (*database.FeedSchedule, error) {
	return nil, nil
}
func (m *mockRepo) FindItemsByRootDomain(rootDomain string, limit int) ([]database.Item, error) {
	return nil, nil
}

// Implement EnqueueSync
func (m *mockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration, lockDuration time.Duration) error {
	m.enqueueSyncCalled = true
	m.lastId = id
	return nil
}

func TestSyncFeed(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	id := uuid.New()
	err = s.SyncFeed(id)

	assert.NoError(t, err)
	assert.True(t, repo.enqueueSyncCalled)
	assert.Equal(t, id, repo.lastId)
}

func TestPoll(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to test exit

	s, err := New(ctx, cfg, repo)
	assert.NoError(t, err)
	s.Poll()
}

func TestStopPolling(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	id := uuid.New()
	err = s.StopPolling(id)

	assert.NoError(t, err)
	assert.True(t, repo.removeSyncCalled)
	assert.Equal(t, id, repo.lastId)
}

func TestSyncFeedInternal(t *testing.T) {
	// developer trigger
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping test on GitHub Actions")
	}

	t.Skip("Skipping test")

	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := database.NewRepository(cfg)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		url     string
		wantErr bool
		enabled bool
	}{
		{
			name:    "Localhost Feed",
			url:     "http://localhost:8003/feed",
			wantErr: false,
			enabled: true,
		},
		{
			name:    "WordPress Feed",
			url:     "http://wordpress/feed",
			wantErr: true,
			enabled: false,
		},
	}

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.enabled {
				t.Skip("Test case disabled")
			}

			u, err := url.Parse(tt.url)
			assert.NoError(t, err)

			feed, err := repo.FindFeedByUrl(u)
			assert.NoError(t, err)
			if assert.NotNil(t, feed) {
				err = s.updatingFeed(feed)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestExtractMediaContent(t *testing.T) {
	s := &Syncer{logger: slog.Default()}

	tests := []struct {
		name     string
		item     *gofeed.Item
		expected *database.MediaContent
	}{
		{
			name:     "No Extensions",
			item:     &gofeed.Item{},
			expected: nil,
		},
		{
			name: "Full Media Content",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":    "http://example.com/video.mp4",
								"type":   "video/mp4",
								"medium": "video",
								"width":  "1920",
								"height": "1080",
							},
							Children: map[string][]ext.Extension{
								"title":  {{Name: "title", Value: "Video Title"}},
								"credit": {{Name: "credit", Value: "Getty Images"}},
							},
						}},
						"thumbnail": {{
							Name: "thumbnail",
							Attrs: map[string]string{
								"url":    "http://example.com/thumb.jpg",
								"width":  "300",
								"height": "200",
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "http://example.com/video.mp4",
				Type:   "video/mp4",
				Medium: "video",
				Width:  1920,
				Height: 1080,
				Title:  "Video Title",
				Credit: "Getty Images",
				Thumbnail: &database.MediaThumbnail{
					URL:    "http://example.com/thumb.jpg",
					Width:  300,
					Height: 200,
				},
			},
		},
		{
			name: "Item Level Credit",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":    "http://example.com/image.jpg",
								"type":   "image/jpeg",
								"medium": "image",
							},
						}},
						"credit": {{Name: "credit", Value: "AP Photo"}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "http://example.com/image.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  1920,
				Height: 1080,
				Credit: "AP Photo",
			},
		},
		{
			name: "Image No Thumbnail",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":    "http://example.com/image.jpg",
								"type":   "image/jpeg",
								"medium": "image",
								"width":  "1920",
								"height": "1080",
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "http://example.com/image.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  1920,
				Height: 1080,
			},
		},
		{
			name: "Missing Dimensions (Defaults)",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":    "http://example.com/image.jpg",
								"type":   "image/jpeg",
								"medium": "image",
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "http://example.com/image.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  1920,
				Height: 1080,
			},
		},
		{
			name: "Partial Dimensions (Recalculate from URL)",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":    "http://example.com/image.jpg?width=800",
								"type":   "image/jpeg",
								"medium": "image",
								"width":  "800",
								// height missing
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "http://example.com/image.jpg?width=800",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  800,
				Height: 450, // 800 * 9 / 16
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.extractMediaContent(tt.item)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestUpdateContent(t *testing.T) {
	s := &Syncer{
		logger: slog.Default(),
	}

	tests := []struct {
		name                string
		item                *gofeed.Item
		res                 *database.ThinkResult
		expectedTitle       string
		expectedDescription string
		expectedContent     string
		checkExtensions     bool
	}{
		{
			name: "Basic Update",
			item: &gofeed.Item{
				Title:       "Original Title",
				Description: "Original Description",
			},
			res: &database.ThinkResult{
				TitleCorrected:       "Corrected Title",
				DescriptionCorrected: "Corrected Description",
			},
			expectedTitle:       "Corrected Title",
			expectedDescription: "Corrected Description<br/><br/>",
			expectedContent:     "",
			checkExtensions:     false,
		},
		{
			name: "With Content Encoded",
			item: &gofeed.Item{
				Title:       "Foobar",
				Description: "Short desc",
				Content:     `<![CDATA[ <p> <a href="https://wwwwhatever...."> <img src="https://foobar" alt="ALT TEXT" /></a><br /><br /></p> ]]>`,
			},
			res: &database.ThinkResult{
				TitleCorrected:       "Corrected Foobar",
				DescriptionCorrected: "Corrected Short desc",
				OverallReason:        "Analysis",
				Clickbait:            1.0,
				Framing:              1.0,
				Persuasive:           1.0,
				HyperStimulus:        1.0,
				Speculative:          1.0,
			},
			expectedTitle:       "Corrected Foobar",
			expectedDescription: "Corrected Short desc<br/><br/>Analysis",
			expectedContent:     "",
			checkExtensions:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.updateContent(tt.item, tt.res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTitle, tt.item.Title)
			assert.Equal(t, tt.expectedDescription, tt.item.Description)
			assert.Equal(t, tt.expectedContent, tt.item.Content)

			if tt.checkExtensions {
				if assert.NotNil(t, tt.item.Extensions) {
					media, ok := tt.item.Extensions["media"]
					if assert.True(t, ok) {
						contentExt, ok := media["content"]
						if assert.True(t, ok) && assert.NotEmpty(t, contentExt) {
							assert.Equal(t, "https://foobar", contentExt[0].Attrs["url"])
							assert.Equal(t, "image", contentExt[0].Attrs["medium"])
							assert.Equal(t, "1920", contentExt[0].Attrs["width"])
							assert.Equal(t, "1080", contentExt[0].Attrs["height"])
						}

						_, creditExists := media["credit"]
						assert.False(t, creditExists, "media:credit should be removed")
					}
				}
			}
		})
	}
}

func TestWantedDomains(t *testing.T) {
	repo := &mockRepo{}
	cfg := &config.Config{}
	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		feed     *database.Feed
		expected []string
		wantErr  bool
	}{
		{
			name: "Enforce Disabled",
			feed: &database.Feed{
				EnforceFeedDomain: false,
				URL:               "http://example.com/feed",
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Simple Domain",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://example.com/feed",
			},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Subdomain",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://blog.example.com/feed",
			},
			expected: []string{"example.com"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Public Suffix",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://example.co.uk/feed",
			},
			expected: []string{"example.co.uk"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - Localhost",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://localhost:8003/feed",
			},
			expected: []string{"localhost"},
			wantErr:  false,
		},
		{
			name: "Enforce Enabled - WordPress",
			feed: &database.Feed{
				EnforceFeedDomain: true,
				URL:               "http://wordpress/feed",
			},
			expected: []string{"wordpress"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domains, err := s.wantedDomains(tt.feed)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, domains)
			}
		})
	}
}
