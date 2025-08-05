package deframer

import (
	_ "embed"
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/downloader"
	"github.com/egandro/news-deframer/pkg/downloader/mock_downloader"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/openai/mock_openai"
	"github.com/egandro/news-deframer/pkg/source"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

//go:embed testing/feed.xml.testing
var rssContent string

//go:embed testing/source.json
var sourceContent string

func setupTestDeframer(t *testing.T, ai openai.OpenAI, src *source.Source, downloader downloader.Downloader) (Deframer, error) {
	// Use in-memory SQLite for testing
	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	prompts := make(map[string]source.Prompt)
	for _, prompt := range src.Prompts {
		prompts[prompt.Language] = prompt
	}

	res := &deframer{
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
	source, err := source.ParseString(sourceContent)
	downloaderMock := mock_downloader.NewMockDownloader(ctrl)
	downloaderMock.EXPECT().DownloadRSSFeed(gomock.Any()).Return(rssContent, nil).Times(1)

	d, err := setupTestDeframer(t, openAIMock, source, downloaderMock)

	assert.NoError(t, err)
	assert.NotNil(t, d, "Deframer should be initialized")

	err = d.UpdateFeeds()
	assert.NoError(t, err)

}
