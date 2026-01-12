package syncer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformContent(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected MediaData
	}{
		{
			name: "News Example",
			// Simulates structure: Image inside link, text mixed with em tags, link at end (ignored)
			content: `<p> <a href="link"><img src="https://example.com/image.jpg?width=800" alt="Description | Credit Name" /></a> <br/> <br/>Main text start. <em>Emphasized</em> text end.[<a href="link">more</a>]</p>`,
			expected: MediaData{
				URL:    "https://example.com/image.jpg?width=800",
				Width:  800,
				Height: 450, // 800 * 9 / 16
				Medium: "image",
			},
		},
		{
			name:    "No Image",
			content: `<p>Just text content.</p>`,
			expected: MediaData{
				URL:    "",
				Width:  0,
				Height: 0,
				Medium: "image",
			},
		},
		{
			name:    "Image No Dimensions No Credit",
			content: `<img src="https://example.com/image.jpg" alt="Just Description" />`,
			expected: MediaData{
				URL:    "https://example.com/image.jpg",
				Width:  1920, // Default
				Height: 1080, // Default
				Medium: "image",
			},
		},
		{
			name:    "CDATA Wrapped",
			content: `<![CDATA[<p>Text inside CDATA</p>]]>`,
			expected: MediaData{
				URL:    "",
				Width:  0,
				Height: 0,
				Medium: "image",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformContent(tt.content)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
