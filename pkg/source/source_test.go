package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	feedJSON := `
	{
		"feeds": [
			{
				"rss_url": "https://example.com/rss",
				"language": "en"
			},
			{
				"rss_url": "https://example.com/rss2",
				"language": "de"
			}
		],
		"prompts": [
			{
				"user": "Summarize this feed in English.",
				"system": "You are a reporter.",
				"language": "en"
			},
			{
				"user": "Fasse den Text in Deutsch zusammen.",
				"system": "Du bist ein Reporter.",
				"language": "de"
			}
		]
	}
	`
	source, err := ParseString(feedJSON)
	assert.NoError(t, err)
	assert.Len(t, source.Feeds, 2)

	feed1 := source.Feeds[0]
	assert.EqualValues(t, "https://example.com/rss", feed1.RSS_URL)
	assert.EqualValues(t, "en", feed1.Language)

	feed2 := source.Feeds[1]
	assert.EqualValues(t, "https://example.com/rss2", feed2.RSS_URL)
	assert.EqualValues(t, "de", feed2.Language)

	assert.Len(t, source.Prompts, 2)

	prompt1 := source.Prompts[0]
	assert.EqualValues(t, "en", prompt1.Language)
	assert.EqualValues(t, "Summarize this feed in English.", prompt1.User)
	assert.EqualValues(t, "You are a reporter.", prompt1.System)

	prompt2 := source.Prompts[1]
	assert.EqualValues(t, "de", prompt2.Language)
	assert.EqualValues(t, "Fasse den Text in Deutsch zusammen.", prompt2.User)
	assert.EqualValues(t, "Du bist ein Reporter.", prompt2.System)
}
