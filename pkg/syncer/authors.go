package syncer

import (
	"html"
	"regexp"
	"strings"
	"sync"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
)

var authorPipeSuffixRegex = regexp.MustCompile(`\s*\|.*$`)
var authorWhitespaceRegex = regexp.MustCompile(`\s+`)

var authorJoinersByLanguage = map[string][]string{
	"af":  {"en"},
	"am":  {"እና"},
	"ar":  {"و"},
	"as":  nil,
	"az":  {"və"},
	"be":  nil,
	"bg":  {"и"},
	"bn":  {"এবং"},
	"bs":  nil,
	"ca":  {"i"},
	"cs":  {"a"},
	"cy":  nil,
	"da":  {"og"},
	"de":  {"und"},
	"dsb": {"a"},
	"el":  {"και"},
	"en":  {"and"},
	"es":  {"y", "e"},
	"et":  {"ja"},
	"eu":  {"eta"},
	"fa":  {"و"},
	"fi":  {"ja"},
	"fo":  {"og"},
	"fr":  {"et"},
	"ga":  {"agus"},
	"gd":  nil,
	"grc": {"και"},
	"gu":  {"અને"},
	"he":  {"ו"},
	"hi":  {"और"},
	"hr":  {"i"},
	"hsb": {"a"},
	"hu":  {"és"},
	"hy":  {"և"},
	"id":  {"dan"},
	"is":  {"og"},
	"it":  {"e"},
	"ja":  {"と"},
	"kk":  nil,
	"kn":  {"ಮತ್ತು"},
	"ko":  {"및", "와", "과", "그리고"},
	"ku":  nil,
	"ky":  {"жана"},
	"la":  {"et"},
	"lb":  {"an"},
	"lg":  {"ne"},
	"lij": {"e"},
	"lo":  nil,
	"lt":  {"ir"},
	"lv":  {"un"},
	"mi":  nil,
	"mk":  {"и"},
	"ml":  {"ഒപ്പം"},
	"mn":  nil,
	"mr":  {"आणि"},
	"ms":  {"dan"},
	"my":  nil,
	"nb":  {"og"},
	"ne":  {"र"},
	"nl":  {"en"},
	"nn":  {"og"},
	"pa":  nil,
	"pl":  {"i"},
	"ps":  nil,
	"pt":  {"e"},
	"ro":  {"si", "și"},
	"ru":  {"и"},
	"sa":  {"च"},
	"sd":  nil,
	"si":  {"සහ"},
	"sk":  {"a"},
	"sl":  {"in"},
	"sq":  {"dhe"},
	"sr":  {"и", "i"},
	"sw":  nil,
	"sv":  {"och"},
	"ta":  {"மற்றும்"},
	"te":  {"మరియు"},
	"tg":  nil,
	"th":  {"และ"},
	"ti":  {"ን"},
	"tk":  nil,
	"tl":  {"at"},
	"tn":  {"le"},
	"tr":  {"ve"},
	"tt":  {"һәм"},
	"ug":  nil,
	"uk":  {"і", "та"},
	"ur":  {"اور"},
	"uz":  nil,
	"vi":  {"và"},
	"xx":  {"and"},
	"yi":  nil,
	"yo":  {"ati"},
	"zh":  {"和", "與", "与", "及"},
}

var authorJoinerRegexCache sync.Map

func (s *Syncer) extractAndNormalizeAuthors(item *gofeed.Item, language string) database.StringArray {
	if item == nil {
		return nil
	}

	firstAuthor := strings.TrimSpace(s.guessAuthorsByTag(item))
	if firstAuthor == "" {
		return nil
	}

	normalized := html.UnescapeString(firstAuthor)
	normalized = strings.ReplaceAll(normalized, "&", ",")
	normalized = replaceAuthorJoinersByLanguage(normalized, language)

	parts := strings.FieldsFunc(normalized, func(r rune) bool {
		return r == ',' || r == '/'
	})

	authors := make(database.StringArray, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		author := strings.TrimSpace(part)
		author = authorPipeSuffixRegex.ReplaceAllString(author, "")
		author = strings.TrimSpace(authorWhitespaceRegex.ReplaceAllString(author, " "))
		if author == "" {
			continue
		}
		if _, ok := seen[author]; ok {
			continue
		}
		seen[author] = struct{}{}
		authors = append(authors, author)
	}

	if len(authors) == 0 {
		return nil
	}

	return authors
}

func replaceAuthorJoinersByLanguage(value string, language string) string {
	joiners := authorJoinersForLanguage(language)
	for _, joiner := range joiners {
		value = authorJoinerRegex(joiner).ReplaceAllString(value, ",")
	}
	return value
}

func authorJoinersForLanguage(language string) []string {
	language = normalizeAuthorLanguage(language)
	if joiners, ok := authorJoinersByLanguage[language]; ok && len(joiners) > 0 {
		return joiners
	}
	return authorJoinersByLanguage["xx"]
}

func normalizeAuthorLanguage(language string) string {
	language = strings.TrimSpace(strings.ToLower(language))
	if language == "" {
		return "xx"
	}
	return strings.Split(language, "-")[0]
}

func authorJoinerRegex(joiner string) *regexp.Regexp {
	if cached, ok := authorJoinerRegexCache.Load(joiner); ok {
		return cached.(*regexp.Regexp)
	}

	pattern := regexp.MustCompile(`(?i)\s+` + regexp.QuoteMeta(joiner) + `\s+`)
	actual, _ := authorJoinerRegexCache.LoadOrStore(joiner, pattern)
	return actual.(*regexp.Regexp)
}

func (s *Syncer) guessAuthorsByTag(item *gofeed.Item) string {
	if item == nil {
		return ""
	}

	for _, author := range item.Authors {
		if author == nil {
			continue
		}
		if name := strings.TrimSpace(author.Name); name != "" {
			return name
		}
	}

	candidates := []struct {
		namespace string
		name      string
		child     string
	}{
		{namespace: "dc", name: "creator"},
		{namespace: "atom", name: "author", child: "name"},
		{namespace: "atom", name: "name"},
		{namespace: "content", name: "author"},
		{namespace: "", name: "byline"},
		{namespace: "", name: "creator"},
		{namespace: "slash", name: "author"},
		{namespace: "sy", name: "author"},
		{namespace: "foaf", name: "name"},
		{namespace: "creativeCommons", name: "attributionName"},
	}

	for _, candidate := range candidates {
		if value := getAuthorTagText(item.Extensions, candidate.namespace, candidate.name, candidate.child); value != "" {
			return value
		}
	}

	return ""
}

func getAuthorTagText(extensions map[string]map[string][]ext.Extension, namespace string, name string, child string) string {
	if extensions == nil {
		return ""
	}

	namespaceEntries, ok := extensions[namespace]
	if !ok {
		return ""
	}

	entries, ok := namespaceEntries[name]
	if !ok {
		return ""
	}

	for _, entry := range entries {
		if child != "" {
			children, ok := entry.Children[child]
			if !ok {
				continue
			}
			for _, childEntry := range children {
				if value := strings.TrimSpace(childEntry.Value); value != "" {
					return value
				}
			}
			continue
		}

		if value := strings.TrimSpace(entry.Value); value != "" {
			return value
		}
	}

	return ""
}
