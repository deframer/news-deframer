package netutil

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRootDomain(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		want     string
	}{
		{
			name:     "Simple .com",
			inputURL: "http://example.com/foo",
			want:     "example.com",
		},
		{
			name:     "Subdomain",
			inputURL: "http://blog.example.com/foo",
			want:     "example.com",
		},
		{
			name:     "Complex TLD",
			inputURL: "http://news.bbc.co.uk/story",
			want:     "bbc.co.uk",
		},
		{
			name:     "Localhost",
			inputURL: "http://localhost:8080/rss",
			want:     "localhost",
		},
		{
			name:     "IP Address",
			inputURL: "http://127.0.0.1/rss",
			want:     "127.0.0.1",
		},
		{
			name:     "Deep Subdomain",
			inputURL: "https://a.b.c.example.org",
			want:     "example.org",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.inputURL)
			assert.NoError(t, err)
			got := GetRootDomain(u)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Basic",
			input:    "http://example.com/page",
			expected: "http://example.com/page",
		},
		{
			name:     "Trailing Slash",
			input:    "http://example.com/page/",
			expected: "http://example.com/page",
		},
		{
			name:     "Fragment",
			input:    "http://example.com/page#section",
			expected: "http://example.com/page",
		},
		{
			name:     "Tracking Params",
			input:    "http://example.com/page?utm_source=twitter&utm_medium=social&id=123",
			expected: "http://example.com/page?id=123",
		},
		{
			name:     "All Tracking Params",
			input:    "http://example.com/page?utm_source=a&utm_medium=b&utm_campaign=c&fbclid=d",
			expected: "http://example.com/page",
		},
		{
			name:     "Mixed Params",
			input:    "http://example.com/page?keep=me&utm_source=remove",
			expected: "http://example.com/page?keep=me",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeURL(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}
