package think

import (
	"fmt"
	"strings"
)

var localizedCategories = map[string][]string{
	"en": {"politics", "world", "business", "sport", "culture", "technology", "health", "finance", "science", "environment", "travel", "lifestyle", "games", "history", "opinion", "other"},
	"de": {"politik", "welt", "wirtschaft", "sport", "kultur", "technologie", "gesundheit", "finanzen", "wissenschaft", "umwelt", "reisen", "lifestyle", "spiele", "geschichte", "meinung", "sonstiges"},
	"da": {"politik", "verden", "erhverv", "sport", "kultur", "teknologi", "sundhed", "finans", "videnskab", "miljo", "rejse", "livsstil", "spil", "historie", "mening", "andet"},
}

func validateLocalizedCategory(language, category string) error {
	allowed, err := localizedCategoriesFor(language)
	if err != nil {
		return err
	}

	for _, candidate := range allowed {
		if strings.EqualFold(candidate, category) {
			return nil
		}
	}

	return fmt.Errorf("invalid category %q for language %s", category, language)
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

	return append([]string(nil), allowed...), nil
}
