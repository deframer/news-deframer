package syncer

import (
	"testing"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
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
				Alt:    "Description | Credit Name",
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
				Alt:    "Just Description",
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

func TestExtractMediaContent(t *testing.T) {
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
			name: "Specific Media Content Tag",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"content": {{
							Name: "content",
							Attrs: map[string]string{
								"url":        "https://bar/foo.jpg",
								"type":       "image/jpeg",
								"expression": "full",
								"width":      "931",
								"height":     "523",
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "https://bar/foo.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  931,
				Height: 523,
			},
		},
		{
			name: "Media Group Content",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"group": {{
							Name: "group",
							Children: map[string][]ext.Extension{
								"content": {{
									Name: "content",
									Attrs: map[string]string{
										"url":        "https://foo.jpg",
										"type":       "image/jpeg",
										"medium":     "image",
										"expression": "full",
										"width":      "300",
										"height":     "180",
									},
								}},
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "https://foo.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  300,
				Height: 180,
			},
		},
		{
			name: "Media Group Prefers Larger Rendition",
			item: &gofeed.Item{
				Extensions: map[string]map[string][]ext.Extension{
					"media": {
						"group": {{
							Name: "group",
							Children: map[string][]ext.Extension{
								"content": {
									{
										Name: "content",
										Attrs: map[string]string{
											"url":        "https://example.com/image-large.jpg",
											"type":       "image/jpeg",
											"medium":     "image",
											"expression": "full",
											"width":      "870",
											"height":     "522",
										},
									},
									{
										Name: "content",
										Attrs: map[string]string{
											"url":        "https://example.com/image-small.jpg",
											"type":       "image/jpeg",
											"medium":     "image",
											"expression": "full",
											"width":      "300",
											"height":     "180",
										},
									},
								},
							},
						}},
					},
				},
			},
			expected: &database.MediaContent{
				URL:    "https://example.com/image-large.jpg",
				Type:   "image/jpeg",
				Medium: "image",
				Width:  870,
				Height: 522,
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
			got, err := extractMediaContent(tt.item)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseDimensions(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		width  int
		height int
	}{
		{
			name:   "Width Query Param",
			url:    "https://example.com/image.jpg?width=800",
			width:  800,
			height: 450,
		},
		{
			name:   "Invalid URL Defaults",
			url:    "%",
			width:  1920,
			height: 1080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := parseDimensions(tt.url)
			assert.Equal(t, tt.width, w)
			assert.Equal(t, tt.height, h)
		})
	}
}
