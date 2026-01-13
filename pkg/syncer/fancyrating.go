package syncer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sync"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
)

//go:embed templates/*.txt
var templateFS embed.FS

var (
	templateCache = make(map[string]*template.Template)
	templateMu    sync.RWMutex
)

// applyFancyRatingText call this after you have the Title = TitleCorrected
func applyFancyRatingText(item *gofeed.Item, res *database.ThinkResult, thinkRating float64, language string) {
	if item == nil || res == nil {
		return
	}

	item.Title = fmt.Sprintf("%s %s", getRatingIcon(thinkRating), item.Title)

	stars := createStarRating(res)

	data := struct {
		Stars  StarRating
		Result *database.ThinkResult
	}{
		Stars:  stars,
		Result: res,
	}

	tmpl, err := getRatingTemplate(language)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err == nil {
		item.Description += buf.String()
	}
}

func getRatingTemplate(lang string) (*template.Template, error) {
	if lang == "" {
		lang = "en"
	}

	templateMu.RLock()
	tmpl, ok := templateCache[lang]
	templateMu.RUnlock()
	if ok {
		return tmpl, nil
	}

	filename := fmt.Sprintf("templates/starrating-%s.txt", lang)
	content, err := templateFS.ReadFile(filename)
	if err != nil {
		if lang != "en" {
			return getRatingTemplate("en")
		}
		return nil, err
	}

	tmpl, err = template.New(lang).Parse(string(content))
	if err != nil {
		return nil, err
	}

	templateMu.Lock()
	templateCache[lang] = tmpl
	templateMu.Unlock()

	return tmpl, nil
}
