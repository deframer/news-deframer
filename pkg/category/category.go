package category

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

var englishCategoryLookup = func() map[string]struct{} {
	lookup := make(map[string]struct{}, len(EnglishCategories))
	for _, category := range EnglishCategories {
		lookup[category] = struct{}{}
	}
	return lookup
}()

// EnglishCategories is the canonical set used by admin input and stored results.
var EnglishCategories = append([]string(nil), canonicalCategories...)

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

func ValidateLocalizedCategory(language, category string) error {
	_, err := NormalizeCategory(language, category)
	return err
}

func ValidateEnglishCategory(category string) error {
	if _, ok := englishCategoryLookup[category]; !ok {
		return fmt.Errorf("invalid category %q", category)
	}

	return nil
}

func FirstLocalizedCategory(language string) (string, error) {
	allowed, err := LocalizedCategoriesFor(language)
	if err != nil {
		return "", err
	}
	if len(allowed) == 0 {
		return "", fmt.Errorf("no localized categories configured for language: %s", language)
	}

	return allowed[0], nil
}

func LocalizedCategoriesFor(language string) ([]string, error) {
	allowed, ok := localizedCategories[language]
	if !ok || len(allowed) == 0 {
		return nil, fmt.Errorf("no localized categories configured for language: %s", language)
	}

	categories := make([]string, 0, len(allowed))
	categories = append(categories, allowed...)
	sort.Strings(categories)

	return categories, nil
}

func NormalizeCategory(language, category string) (string, error) {
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

func GetUnknowCategory() string {
	return "other"
}
