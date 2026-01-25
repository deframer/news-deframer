package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripHTML(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple paragraph",
			input:    "<p>Hello, world!</p>",
			expected: "Hello, world!",
		},
		{
			name:     "String with no HTML",
			input:    "This is a plain string.",
			expected: "This is a plain string.",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "String with only HTML tags",
			input:    "<div><p></p></div>",
			expected: "",
		},
		{
			name:     "Multiple tags",
			input:    "<div><h1>Title</h1><p>Some text.</p></div>",
			expected: "Title Some text.",
		},
		{
			name:     "Line breaks",
			input:    "First line.<br>Second line.",
			expected: "First line. Second line.",
		},
		{
			name:     "HTML entities that form tags",
			input:    "&lt;p&gt;Hello&lt;/p&gt;",
			expected: "<p>Hello</p>",
		},
		{
			name:     "Regular HTML entities",
			input:    "Me &amp; you",
			expected: "Me & you",
		},
		{
			name:     "Mixed content",
			input:    "Some <b>bold</b> text and <i>italic</i>.",
			expected: "Some bold text and italic.",
		},
		{
			name:     "Extra whitespace",
			input:    "  <p>  leading and trailing spaces  </p>  ",
			expected: "leading and trailing spaces",
		},
		{
			name:     "Image with alt text",
			input:    `Text <img src="image.png" alt="This should be removed"> more text`,
			expected: "Text more text",
		},
		{
			name:     "Only image with alt",
			input:    `<img src="image.png" alt="Hidden">`,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := StripHTML(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
