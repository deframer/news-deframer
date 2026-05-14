package think

import (
	"fmt"
	"sort"
	"strings"
)

// en is the default language. Translations must stay in the same position.
var localizedCategories = map[string][]string{
	"en": {"politics", "world", "business", "sport", "culture", "technology", "health", "finance", "science", "environment", "travel", "lifestyle", "games", "history", "opinion", "other"},
	"de": {"politik", "welt", "wirtschaft", "sport", "kultur", "technologie", "gesundheit", "finanzen", "wissenschaft", "umwelt", "reisen", "lifestyle", "spiele", "geschichte", "meinung", "sonstiges"},
	"da": {"politik", "verden", "erhverv", "sport", "kultur", "teknologi", "sundhed", "finans", "videnskab", "miljo", "rejse", "livsstil", "spil", "historie", "mening", "andet"},
}

var canonicalCategories = localizedCategories["en"]

var builtCategoryLookupByLanguage = buildCategoryLookupByLanguage()

func buildCategoryLookupByLanguage() map[string]map[string]string {
	lookup := make(map[string]map[string]string, len(localizedCategories))
	for language, labels := range localizedCategories {
		languageLookup := make(map[string]string, len(labels))
		for idx, localized := range labels {
			if idx >= len(canonicalCategories) {
				break
			}
			languageLookup[localized] = canonicalCategories[idx]
		}
		lookup[language] = languageLookup
	}
	return lookup
}

func validateLocalizedCategory(language, category string) error {
	_, err := normalizeCategory(language, category)
	if err != nil {
		return err
	}

	return nil
}

func firstLocalizedCategory(language string) (string, error) {
	allowed, err := localizedCategoriesFor(language)
	if err != nil {
		return "", err
	}
	if len(allowed) == 0 {
		return "", fmt.Errorf("no localized categories configured for language: %s", language)
	}

	return allowed[0], nil
}

func localizedCategoriesFor(language string) ([]string, error) {
	allowed, ok := localizedCategories[language]
	if !ok || len(allowed) == 0 {
		return nil, fmt.Errorf("no localized categories configured for language: %s", language)
	}

	categories := make([]string, 0, len(allowed))
	categories = append(categories, allowed...)
	sort.Strings(categories)

	return categories, nil
}

func normalizeCategory(language, category string) (string, error) {
	categories, ok := builtCategoryLookupByLanguage[language]
	if !ok || len(categories) == 0 {
		return "", fmt.Errorf("no localized categories configured for language: %s", language)
	}

	for localized, canonical := range categories {
		if strings.EqualFold(localized, category) {
			if canonical == "" {
				return "", fmt.Errorf("no canonical mapping for category %q and language %s", category, language)
			}

			return canonical, nil
		}
	}

	return "", fmt.Errorf("invalid category %q for language %s", category, language)
}
