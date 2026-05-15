package think

import (
	"fmt"

	categorypkg "github.com/deframer/news-deframer/pkg/category"
	"github.com/deframer/news-deframer/pkg/database"
)

func validateAndNormalizeThinkResult(language string, res *database.ThinkResult, ignoreCategoryErrors bool) error {
	const errFmt = "ThinkResult is out of bounds 0.0 - 1.0: %s is %.1f"

	if res == nil {
		return nil
	}

	if res.Framing < 0.0 || res.Framing > 1.0 {
		return fmt.Errorf(errFmt, "Framing", res.Framing)
	}
	if res.Clickbait < 0.0 || res.Clickbait > 1.0 {
		return fmt.Errorf(errFmt, "Clickbait", res.Clickbait)
	}
	if res.Persuasive < 0.0 || res.Persuasive > 1.0 {
		return fmt.Errorf(errFmt, "Persuasive", res.Persuasive)
	}
	if res.HyperStimulus < 0.0 || res.HyperStimulus > 1.0 {
		return fmt.Errorf(errFmt, "HyperStimulus", res.HyperStimulus)
	}
	if res.Speculative < 0.0 || res.Speculative > 1.0 {
		return fmt.Errorf(errFmt, "Speculative", res.Speculative)
	}
	if res.Overall < 0.0 || res.Overall > 1.0 {
		return fmt.Errorf(errFmt, "Overall", res.Overall)
	}

	category, err := normalizeThinkResultCategory(language, res.Category, ignoreCategoryErrors)
	if err != nil {
		return err
	}
	res.Category = category

	return nil
}

func normalizeThinkResultCategory(language, category string, ignoreCategoryErrors bool) (string, error) {
	normalizedCategory, err := categorypkg.NormalizeCategory(language, category)
	if err != nil {
		if ignoreCategoryErrors {
			return categorypkg.GetUnknowCategory(), nil
		}
		return "", err
	}

	return normalizedCategory, nil
}
