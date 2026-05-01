package syncer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"net/url"

	"golang.org/x/net/html"
)

// MediaData represents our target structure
type MediaData struct {
	URL    string
	Width  int
	Height int
	Medium string
	Alt    string
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
				switch a.Key {
				case "src":
					data.URL = a.Val
					// Extract width from URL query params
					data.Width, data.Height = parseDimensions(a.Val)
				case "alt":
					data.Alt = a.Val
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

func extractFromContentTag(item *gofeed.Item) *MediaData {
	if item == nil || item.Content == "" {
		return nil
	}
	mediaData, err := transformContent(item.Content)
	if err == nil && mediaData.URL != "" {
		return &mediaData
	}
	return nil
}

func extractFromEnclosureTag(item *gofeed.Item) *MediaData {
	if item == nil {
		return nil
	}
	for _, enc := range item.Enclosures {
		if strings.HasPrefix(enc.Type, "image/") && enc.URL != "" {
			w, h := parseDimensions(enc.URL)
			return &MediaData{
				URL:    enc.URL,
				Medium: "image",
				Width:  w,
				Height: h,
			}
		}
	}
	return nil
}

func extractFromDescriptionFallback(res *database.ThinkResult) *MediaData {
	if res == nil {
		return nil
	}

	target := res.DescriptionOriginal
	mediaData, err := transformContent(target)
	if err != nil || mediaData.URL == "" {
		target = res.TitleOriginal
		mediaData, err = transformContent(target)
	}

	if err == nil && mediaData.URL != "" {
		return &mediaData
	}
	return nil
}

func applyMediaData(item *gofeed.Item, data *MediaData) {
	if item == nil || data == nil {
		return
	}
	if item.Extensions == nil {
		item.Extensions = make(map[string]map[string][]ext.Extension)
	}
	if item.Extensions["media"] == nil {
		item.Extensions["media"] = make(map[string][]ext.Extension)
	}

	// <media:credit> always delete this
	delete(item.Extensions["media"], "credit")

	attrs := map[string]string{
		"url":    data.URL,
		"medium": data.Medium,
	}
	if data.Width > 0 {
		attrs["width"] = fmt.Sprintf("%d", data.Width)
	}
	if data.Height > 0 {
		attrs["height"] = fmt.Sprintf("%d", data.Height)
	}

	item.Extensions["media"]["content"] = []ext.Extension{{
		Name:  "content",
		Attrs: attrs,
	}}

	if len(data.Alt) > 0 {
		item.Extensions["media"]["description"] = []ext.Extension{{
			Name:  "description",
			Value: data.Alt,
		}}
	}
}

func extractMediaContent(item *gofeed.Item) (*database.MediaContent, error) {
	if item == nil || item.Extensions == nil {
		return nil, nil
	}

	mediaExt, ok := item.Extensions["media"]
	if !ok {
		return nil, nil
	}

	contentExt := findBestMediaExtension(mediaExt, "content")
	if contentExt == nil {
		contentExt = findBestMediaExtension(mediaExt, "thumbnail")
	}
	if contentExt == nil {
		return nil, nil
	}

	mc := &database.MediaContent{
		URL:    contentExt.Attrs["url"],
		Type:   contentExt.Attrs["type"],
		Medium: contentExt.Attrs["medium"],
	}

	if mc.Medium == "" {
		if strings.HasPrefix(mc.Type, "image/") {
			mc.Medium = "image"
		} else if strings.HasPrefix(mc.Type, "video/") {
			mc.Medium = "video"
		}
	}
	if mc.Medium == "" && contentExt.Name == "thumbnail" {
		mc.Medium = "image"
	}

	w, _ := strconv.Atoi(contentExt.Attrs["width"])
	h, _ := strconv.Atoi(contentExt.Attrs["height"])
	if w == 0 || h == 0 {
		w, h = parseDimensions(mc.URL)
	}
	mc.Width = w
	mc.Height = h

	if titles, ok := contentExt.Children["title"]; ok && len(titles) > 0 {
		mc.Title = titles[0].Value
	}
	if descs, ok := contentExt.Children["description"]; ok && len(descs) > 0 {
		mc.Description = descs[0].Value
	}
	if credits, ok := contentExt.Children["credit"]; ok && len(credits) > 0 {
		mc.Credit = credits[0].Value
	}

	if mc.Title == "" {
		if titles, ok := mediaExt["title"]; ok && len(titles) > 0 {
			mc.Title = titles[0].Value
		}
	}
	if mc.Description == "" {
		if descs, ok := mediaExt["description"]; ok && len(descs) > 0 {
			mc.Description = descs[0].Value
		}
	}
	if mc.Credit == "" {
		if credits, ok := mediaExt["credit"]; ok && len(credits) > 0 {
			mc.Credit = credits[0].Value
		}
	}

	if thumbs, ok := mediaExt["thumbnail"]; ok && len(thumbs) > 0 {
		tExt := thumbs[0]
		if tURL := tExt.Attrs["url"]; tURL != "" {
			mt := &database.MediaThumbnail{URL: tURL}
			tw, _ := strconv.Atoi(tExt.Attrs["width"])
			th, _ := strconv.Atoi(tExt.Attrs["height"])
			if tw == 0 || th == 0 {
				tw, th = parseDimensions(tURL)
			}
			mt.Width = tw
			mt.Height = th
			mc.Thumbnail = mt
		}
	}

	return mc, nil
}

func findBestMediaExtension(exts map[string][]ext.Extension, name string) *ext.Extension {
	var best *ext.Extension
	for _, values := range exts {
		best = betterMediaExtension(best, findBestMediaExtensionInList(values, name))
	}

	return best
}

func findBestMediaExtensionInList(exts []ext.Extension, name string) *ext.Extension {
	var best *ext.Extension
	for i := range exts {
		if exts[i].Name == name && exts[i].Attrs["url"] != "" {
			val := exts[i]
			best = betterMediaExtension(best, &val)
		}
		for _, children := range exts[i].Children {
			best = betterMediaExtension(best, findBestMediaExtensionInList(children, name))
		}
	}

	return best
}

func betterMediaExtension(current, candidate *ext.Extension) *ext.Extension {
	if candidate == nil {
		return current
	}
	if current == nil {
		return candidate
	}

	if mediaExtensionScore(candidate) > mediaExtensionScore(current) {
		return candidate
	}
	return current
}

func mediaExtensionScore(ext *ext.Extension) int {
	w, _ := strconv.Atoi(ext.Attrs["width"])
	h, _ := strconv.Atoi(ext.Attrs["height"])
	if w == 0 || h == 0 {
		w, h = parseDimensions(ext.Attrs["url"])
	}
	return w * h
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
