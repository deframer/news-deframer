package syncer

import (
	"strings"
)

func scoreToStars(score float64) string {
	// 0.0 = 5 stars (Good/Neutral)
	// 1.0 = 0 stars (Bad/Biased)
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	stars := int((1.0-score)*5.0 + 0.5)
	return strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars)
}
