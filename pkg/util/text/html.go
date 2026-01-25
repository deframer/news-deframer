package text

import (
	"strings"

	"golang.org/x/net/html"
)

// StripHTML removes HTML tags from a string, leaving the text content.
// It also normalizes whitespace, replacing block-level tags with a single space
// and collapsing multiple whitespace characters into one.
func StripHTML(s string) string {
	if s == "" {
		return ""
	}

	// Optimization: if it doesn't look like HTML, return as is
	if !strings.ContainsAny(s, "<&") {
		return s
	}

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		// Return original string on parse error
		return s
	}

	var sb strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "div", "br", "li", "tr", "h1", "h2", "h3", "h4", "h5", "h6":
				sb.WriteByte(' ')
			}
		}
	}
	f(doc)

	return strings.Join(strings.Fields(sb.String()), " ")
}
