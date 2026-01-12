package syncer

import (
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// MediaData represents our target structure
type MediaData struct {
	URL    string
	Width  int
	Height int
	Medium string
}

func transformContent(content string) (MediaData, error) {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "<![CDATA[") && strings.HasSuffix(content, "]]>") {
		content = strings.TrimPrefix(content, "<![CDATA[")
		content = strings.TrimSuffix(content, "]]>")
	}
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return MediaData{}, err
	}

	data := MediaData{Medium: "image"}

	// Recursive function to traverse the HTML tree
	var f func(*html.Node)
	f = func(n *html.Node) {
		// 1. Extract Info from Image Tag
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					data.URL = a.Val
					// Extract width from URL query params
					data.Width, data.Height = parseDimensions(a.Val)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return data, nil
}

// parseDimensions pulls width from query string and calculates 16:9 height
func parseDimensions(imgURL string) (int, int) {
	u, err := url.Parse(imgURL)
	if err != nil {
		return 1920, 1080
	}
	widthStr := u.Query().Get("width")
	w, err := strconv.Atoi(widthStr)
	if err != nil {
		return 1920, 1080
	}
	// Calculate height based on standard 16:9 ratio
	h := (w * 9) / 16
	return w, h
}
