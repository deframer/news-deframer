package main

import (
	"fmt"
	"os"
	"strings"

	categorypkg "github.com/deframer/news-deframer/pkg/category"
)

func validateCategoryList(categories []string) []string {
	validated, errs := validateCategoryListStrict(categories)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
	return validated
}

func validateCategoryListStrict(categories []string) ([]string, []string) {
	validated := make([]string, 0, len(categories))
	var errs []string
	for _, raw := range categories {
		category := strings.TrimSpace(raw)
		if category == "" {
			continue
		}

		if err := categorypkg.ValidateEnglishCategory(category); err != nil {
			errs = append(errs, fmt.Sprintf("invalid category %q; allowed: %s", category, strings.Join(categorypkg.EnglishCategories, ",")))
			continue
		}
		validated = append(validated, category)
	}
	return validated, errs
}
