package main

import (
	"fmt"
	"os"
	"strings"

	categorypkg "github.com/deframer/news-deframer/pkg/category"
)

func validateCategoryList(categories []string) []string {
	validated := make([]string, 0, len(categories))
	for _, raw := range categories {
		category := strings.TrimSpace(raw)
		if category == "" {
			continue
		}

		err := categorypkg.ValidateEnglishCategory(category)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid category %q; allowed: %s\n", category, strings.Join(categorypkg.EnglishCategories, ","))
			os.Exit(1)
		}
		validated = append(validated, category)
	}
	return validated
}
