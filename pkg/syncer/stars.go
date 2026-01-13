package syncer

import (
	"strings"

	"github.com/egandro/news-deframer/pkg/database"
)

type StarRating struct {
	Overall            string
	Clickbait          string
	Framing            string
	PersuasiveIntent   string
	HyperStimulus      string
	SpeculativeContent string
}

func createStarRating(res *database.ThinkResult) StarRating {
	// Average
	sum := res.ClickbaitScore + res.FramingScore + res.PersuasiveScore + res.HyperStimulusScore + res.SpeculativeScore
	avg := sum / 5.0

	return StarRating{
		Clickbait:          scoreToStars(res.ClickbaitScore),
		Framing:            scoreToStars(res.FramingScore),
		PersuasiveIntent:   scoreToStars(res.PersuasiveScore),
		HyperStimulus:      scoreToStars(res.HyperStimulusScore),
		SpeculativeContent: scoreToStars(res.SpeculativeScore),
		Overall:            scoreToStars(avg),
	}
}

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
