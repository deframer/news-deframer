package source

import (
	"encoding/json"
	"os"
)

type Feed struct {
	RSS_URL  string `json:"rss_url"`
	Language string `json:"language"`
}

type Prompt struct {
	User     string `json:"user"`
	System   string `json:"system"`
	Language string `json:"language"`
}

type Source struct {
	Feeds   []Feed   `json:"feeds"`
	Prompts []Prompt `json:"prompts"`
}

// ParseString parses the feed from a JSON string and returns feeds
func ParseString(feedJSON string) (Source, error) {
	var source Source
	err := json.Unmarshal([]byte(feedJSON), &source)
	return source, err
}

// ParseFile parses the feed from a JSON file and returns feeds
func ParseFile(filePath string) (Source, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Source{}, err
	}
	return ParseString(string(data))
}
