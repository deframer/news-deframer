package syncer

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/feeds"
	"github.com/deframer/news-deframer/pkg/think"
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
func (m *mockRepo) CreateFeedSchedule(feedID uuid.UUID) error {
	return nil
}
func (m *mockRepo) FindItemsByRootDomain(rootDomain string, limit int) ([]database.Item, error) {
	return nil, nil
}

func (m *mockRepo) PurgeFeedById(id uuid.UUID) error {
	return nil
}

// Implement EnqueueSync
func (m *mockRepo) EnqueueSync(id uuid.UUID, pollingInterval time.Duration) error {
	m.enqueueSyncCalled = true
	m.lastId = id
	return nil
}

func (m *mockRepo) EnqueueMine(id uuid.UUID, miningInterval time.Duration) error {
	return nil
}

func (m *mockRepo) RemoveMine(id uuid.UUID) error {
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

// func TestSyncFeedInternal(t *testing.T) {
// 	// developer trigger
// 	// if os.Getenv("GITHUB_ACTIONS") == "true" {
// 	// 	t.Skip("Skipping test on GitHub Actions")
// 	// }

// 	t.Skip("Skipping test")

// 	cfg, err := config.Load()
// 	assert.NoError(t, err)

// 	repo, err := database.NewRepository(cfg)
// 	assert.NoError(t, err)

// 	tests := []struct {
// 		name    string
// 		url     string
// 		wantErr bool
// 		enabled bool
// 	}{
// 		{
// 			name:    "Localhost Feed",
// 			url:     "http://localhost:8003/feed",
// 			wantErr: false,
// 			enabled: true,
// 		},
// 		{
// 			name:    "WordPress Feed",
// 			url:     "http://wordpress/feed",
// 			wantErr: true,
// 			enabled: false,
// 		},
// 	}

// 	s, err := New(context.Background(), cfg, repo)
// 	assert.NoError(t, err)

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if !tt.enabled {
// 				t.Skip("Test case disabled")
// 			}

// 			u, err := url.Parse(tt.url)
// 			assert.NoError(t, err)

// 			feed, err := repo.FindFeedByUrl(u)
// 			assert.NoError(t, err)
// 			if assert.NotNil(t, feed) {
// 				err = s.updatingFeed(feed)
// 				if tt.wantErr {
// 					assert.Error(t, err)
// 				} else {
// 					assert.NoError(t, err)
// 				}
// 			}
// 		})
// 	}
// }

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
		{
			name: "Fallback to Thumbnail",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"thumbnail": {{
							Name: "thumbnail",
							Attrs: map[string]string{
								"url":    "https://ichef.bbci.co.uk/ace/standard/240/cpsprodpb/039e/live/9ebb4470-f496-11f0-a422-4ba8a094a8fa.jpg",
								"width":  "240",
								"height": "135",
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "https://ichef.bbci.co.uk/ace/standard/240/cpsprodpb/039e/live/9ebb4470-f496-11f0-a422-4ba8a094a8fa.jpg",
				Medium: "image", // Default injected
				Width:  240,
				Height: 135,
				Thumbnail: &database.MediaThumbnail{
					URL:    "https://ichef.bbci.co.uk/ace/standard/240/cpsprodpb/039e/live/9ebb4470-f496-11f0-a422-4ba8a094a8fa.jpg",
					Width:  240,
					Height: 135,
				},
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
		expectedMediaURL    string
		expectedWidth       string
		expectedHeight      string
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
			expectedMediaURL:    "https://foobar",
			expectedWidth:       "1920",
			expectedHeight:      "1080",
		},
		{
			name: "With Image Enclosure",
			item: &gofeed.Item{
				Title:       "Enclosure Test",
				Description: "Desc",
				Enclosures: []*gofeed.Enclosure{
					{
						URL:  "https://example.com/image.jpg?width=1280",
						Type: "image/jpeg",
					},
				},
			},
			res: &database.ThinkResult{
				TitleCorrected:       "Corrected Enclosure",
				DescriptionCorrected: "Corrected Desc",
			},
			expectedTitle:       "Corrected Enclosure",
			expectedDescription: "Corrected Desc<br/><br/>",
			expectedContent:     "",
			expectedMediaURL:    "https://example.com/image.jpg?width=1280",
			expectedWidth:       "1280",
			expectedHeight:      "720",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.updateContent(tt.item, tt.res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTitle, tt.item.Title)
			assert.Equal(t, tt.expectedDescription, tt.item.Description)
			assert.Equal(t, tt.expectedContent, tt.item.Content)

			if tt.expectedMediaURL != "" {
				if assert.NotNil(t, tt.item.Extensions) {
					media, ok := tt.item.Extensions["media"]
					if assert.True(t, ok) {
						contentExt, ok := media["content"]
						if assert.True(t, ok) && assert.NotEmpty(t, contentExt) {
							assert.Equal(t, tt.expectedMediaURL, contentExt[0].Attrs["url"])
							assert.Equal(t, "image", contentExt[0].Attrs["medium"])
							assert.Equal(t, tt.expectedWidth, contentExt[0].Attrs["width"])
							assert.Equal(t, tt.expectedHeight, contentExt[0].Attrs["height"])
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

func TestDetermineLanguage(t *testing.T) {
	strPtr := func(s string) *string { return &s }
	s := &Syncer{logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))}

	testCases := []struct {
		name         string
		dbFeed       *database.Feed
		parsedFeed   *gofeed.Feed
		expectedLang string
	}{
		{
			name:         "Default to en when no language is specified",
			dbFeed:       &database.Feed{Base: database.Base{ID: uuid.New()}},
			parsedFeed:   &gofeed.Feed{},
			expectedLang: "en",
		},
		{
			name:         "Use language from database feed",
			dbFeed:       &database.Feed{Language: strPtr("de")},
			parsedFeed:   &gofeed.Feed{},
			expectedLang: "de",
		},
		{
			name:         "Use language from parsed feed",
			dbFeed:       &database.Feed{},
			parsedFeed:   &gofeed.Feed{Language: "fr"},
			expectedLang: "fr",
		},
		{
			name:         "Parsed feed language (fr) overrides database language (de)",
			dbFeed:       &database.Feed{Language: strPtr("de")},
			parsedFeed:   &gofeed.Feed{Language: "fr"},
			expectedLang: "fr",
		},
		{
			name:         "Handle complex language codes from parsed feed",
			dbFeed:       &database.Feed{Language: strPtr("de")},
			parsedFeed:   &gofeed.Feed{Language: "en-US"},
			expectedLang: "en",
		},
		{
			name:         "Empty language string in parsed feed uses database language",
			dbFeed:       &database.Feed{Language: strPtr("es")},
			parsedFeed:   &gofeed.Feed{Language: " "},
			expectedLang: "es",
		},
		{
			name:         "Empty language string in database and parsed feed defaults to en",
			dbFeed:       &database.Feed{Base: database.Base{ID: uuid.New()}, Language: strPtr("")},
			parsedFeed:   &gofeed.Feed{Language: ""},
			expectedLang: "en",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			lang := s.determineLanguage(tc.dbFeed, tc.parsedFeed)
			assert.Equal(t, tc.expectedLang, lang)
		})
	}
}

// Mocks for testing processItem
type mockThink struct {
	runFunc func(scope string, language string, req think.Request) (*database.ThinkResult, error)
}

func (m *mockThink) Run(scope string, language string, req think.Request) (*database.ThinkResult, error) {
	if m.runFunc != nil {
		return m.runFunc(scope, language, req)
	}
	return &database.ThinkResult{}, nil
}

type mockFeeds struct {
	extractCategoriesFunc func(item *gofeed.Item) []string
	renderItemFunc        func(ctx context.Context, item *gofeed.Item) (string, error)
}

func (m *mockFeeds) ParseFeed(ctx context.Context, content io.Reader) (*gofeed.Feed, error) {
	return nil, nil
}
func (m *mockFeeds) RenderFeed(ctx context.Context, feed *gofeed.Feed) (string, error) {
	return "", nil
}
func (m *mockFeeds) FilterItems(ctx context.Context, feed *gofeed.Feed, domains []string) []feeds.ItemHashPair {
	return nil
}
func (m *mockFeeds) RenderItem(ctx context.Context, item *gofeed.Item) (string, error) {
	if m.renderItemFunc != nil {
		return m.renderItemFunc(ctx, item)
	}
	return "<item></item>", nil
}
func (m *mockFeeds) ExtractCategories(item *gofeed.Item) []string {
	if m.extractCategoriesFunc != nil {
		return m.extractCategoriesFunc(item)
	}
	return []string{"mocked-cat"}
}

func TestProcessItem_Categories(t *testing.T) {
	repo := &mockRepo{}
	cfg, err := config.Load()
	assert.NoError(t, err)

	s, err := New(context.Background(), cfg, repo)
	assert.NoError(t, err)

	// Replace dependencies with mocks
	s.think = &mockThink{}
	mockFeeds := &mockFeeds{}
	s.feeds = mockFeeds

	feed := &database.Feed{Base: database.Base{ID: uuid.New()}}
	hash := "test-hash"
	item := &gofeed.Item{
		Title: "Test item with categories",
		Link:  "http://example.com/item",
	}
	language := "en"
	currentErrorCount := 0

	var capturedItem *database.Item
	repo.upsertItemFunc = func(dbItem *database.Item) error {
		capturedItem = dbItem
		return nil
	}

	mockFeeds.extractCategoriesFunc = func(item *gofeed.Item) []string {
		return []string{"cat1", "cat2", "cat3"}
	}

	s.processItem(feed, hash, item, language, currentErrorCount)

	assert.NotNil(t, capturedItem)
	assert.Equal(t, database.StringArray{"cat1", "cat2", "cat3"}, capturedItem.Categories)
}
