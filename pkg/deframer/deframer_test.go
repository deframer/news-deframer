package deframer

import (
	"context"
	_ "embed"
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/downloader/mock_downloader"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/openai/mock_openai"
	"github.com/egandro/news-deframer/pkg/source"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed testing/feed.xml.testing
var rssContent string

//go:embed testing/source.json
var sourceContent string

func setupTestDeframer(t *testing.T, ai openai.OpenAI, src *source.Source, downloader downloader.Downloader) (Deframer, error) {
	ctx := context.Background()

	// Use in-memory SQLite for testing
	db, err := database.NewDatabase(":memory:")
	//db, err := database.NewDatabase("./test.db")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	prompts := make(map[string]source.Prompt)
	if src != nil {
		for _, prompt := range src.Prompts {
			prompts[prompt.Language] = prompt
		}
	}

	res := &deframer{
		ctx:        ctx,
		db:         db,
		ai:         ai,
		src:        src,
		downloader: downloader,
		prompts:    prompts,
	}

	return res, nil
}

func TestNewDeframer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	openAIMock := mock_openai.NewMockOpenAI(ctrl)
	source, err := source.ParseString(sourceContent)
	downloader := mock_downloader.NewMockDownloader(ctrl)
	d, err := setupTestDeframer(t, openAIMock, source, downloader)

	//s, err := NewDeframer()
	assert.NoError(t, err)
	assert.NotNil(t, d, "Deframer should be initialized")
}

func TestNewUpdateFeeds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	openAIMock := mock_openai.NewMockOpenAI(ctrl)
	openAIMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).Times(3)
	openAIMock.EXPECT().FuzzyParseJSON(gomock.Any()).
		Return(nil, nil).Times(3)

	source, err := source.ParseString(sourceContent)
	downloaderMock := mock_downloader.NewMockDownloader(ctrl)
	downloaderMock.EXPECT().DownloadRSSFeed(gomock.Any()).Return(rssContent, nil).Times(1)

	d, err := setupTestDeframer(t, openAIMock, source, downloaderMock)

	assert.NoError(t, err)
	assert.NotNil(t, d, "Deframer should be initialized")

	expected := 1
	count, err := d.UpdateFeeds()
	assert.NoError(t, err)
	assert.Equal(t, count, expected)

	// 2nd call - it should take the data from the cache
	expected = 0
	count, err = d.UpdateFeeds()
	assert.NoError(t, err)
	assert.Equal(t, count, expected)
}

func TestDeframe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	openAIMock := mock_openai.NewMockOpenAI(ctrl)
	openAIMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).Times(3)
	openAIMock.EXPECT().FuzzyParseJSON(gomock.Any()).
		Return(nil, nil).Times(3)

	source, err := source.ParseString(sourceContent)
	d, err := setupTestDeframer(t, openAIMock, source, nil)

	assert.NoError(t, err)
	assert.NotNil(t, d, "Deframer should be initialized")

	parser := gofeed.NewParser()
	parsedData, err := parser.ParseString(string(rssContent))
	assert.NoError(t, err)

	str, err := d.DeframeFeed(parsedData, source.Feeds[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, str, "")
}

func TestDeframeItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jsonString := `json
	{
		"title_corrected": "dummy title",
		"framing": 0.2,
		"reason": "My Reason"
	}
	`

	openAI := openai.NewAI("", "", "dummy")
	fuzzy, errFuzzy := openAI.FuzzyParseJSON(jsonString)

	openAIMock := mock_openai.NewMockOpenAI(ctrl)
	openAIMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(jsonString, nil).Times(1)
	openAIMock.EXPECT().FuzzyParseJSON(gomock.Any()).
		Return(fuzzy, errFuzzy).Times(1)
	source, err := source.ParseString(sourceContent)
	d, err := setupTestDeframer(t, openAIMock, source, nil)

	assert.NoError(t, err)
	assert.NotNil(t, d, "Deframer should be initialized")

	parser := gofeed.NewParser()
	parsedData, err := parser.ParseString(string(rssContent))
	assert.NoError(t, err)

	str, err := d.DeframeItem(parsedData.Items[0], source.Feeds[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, str, "")
}
