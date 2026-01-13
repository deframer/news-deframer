package syncer

import (
	"github.com/egandro/news-deframer/pkg/database"
)

// -----------------------------------------------------------------------------
// Hybrid Score (Max + Mean) - RECOMMENDED
// -----------------------------------------------------------------------------

// CalculateHybrid returns a score based on a weighted mix of the Maximum value and the Mean.
//
// parameters:
//   - res: The data pointer.
//   - maxWeight: How much influence the worst outlier has (0.0 to 1.0).
//     A good default is 0.7 (70% Max, 30% Mean).
//
// usage: Best for risk/spam scores. If 'Speculative' is 0.9, the score stays high
//
//	even if all other values are 0.0.
func CalculateHybrid(res *database.ThinkResult, maxWeight float64) float64 {
	if res == nil {
		return 0.0
	}

	// 1. Find the Maximum value and the Sum
	values := []float64{
		res.Framing,
		res.Clickbait,
		res.Persuasive,
		res.HyperStimulus,
		res.Speculative,
	}

	maxVal := 0.0
	sum := 0.0

	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
		sum += v
	}

	// 2. Calculate Mean
	mean := sum / float64(len(values))

	// 3. Mix Max and Mean
	// If maxWeight is 0.7: Score = (0.7 * Max) + (0.3 * Mean)
	return (maxWeight * maxVal) + ((1.0 - maxWeight) * mean)
}

/*
// -----------------------------------------------------------------------------
// Arithmetic Mean
// -----------------------------------------------------------------------------

// CalculateMean returns the simple arithmetic average.
//
// logic: Sum / 5
// pros:  Easy to understand.
// cons:  Hides extreme outliers (e.g., 4 values are 0.0, one is 1.0 -> result is only 0.2).
func CalculateMean(res *database.ThinkResult) float64 {
	if res == nil {
		return 0.0
	}

	sum := res.Framing +
		res.Clickbait +
		res.Persuasive +
		res.HyperStimulus +
		res.Speculative

	return sum / 5.0
}

// -----------------------------------------------------------------------------
// Root Mean Square (RMS)
// -----------------------------------------------------------------------------

// CalculateRMS returns the quadratic mean.
//
// logic: Sqrt(Sum(x^2) / 5)
// pros:  Penalizes outliers. Higher values carry more weight than lower ones.
// usage: Good if you want a mathematical balance that is stricter than the average.
func CalculateRMS(res *database.ThinkResult) float64 {
	if res == nil {
		return 0.0
	}

	sumSquares := (res.Framing * res.Framing) +
		(res.Clickbait * res.Clickbait) +
		(res.Persuasive * res.Persuasive) +
		(res.HyperStimulus * res.HyperStimulus) +
		(res.Speculative * res.Speculative)

	return math.Sqrt(sumSquares / 5.0)
}



type StarRating struct {
	Overall            string
	Clickbait          string
	Framing            string
	PersuasiveIntent   string
	HyperStimulus      string
	SpeculativeContent string
}
func createStarRating(res *database.ThinkResult) StarRating {
	avg := CalculateHybrid(res)

	return StarRating{
		Clickbait:          scoreToStars(res.Clickbait),
		Framing:            scoreToStars(res.Framing),
		PersuasiveIntent:   scoreToStars(res.Persuasive),
		HyperStimulus:      scoreToStars(res.HyperStimulus),
		SpeculativeContent: scoreToStars(res.Speculative),
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
*/
