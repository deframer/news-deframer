package syncer

import (
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHybrid(t *testing.T) {
	tests := []struct {
		name     string
		res      *database.ThinkResult
		expected float64
	}{
		{
			name:     "Nil Result",
			res:      nil,
			expected: 0.0,
		},
		{
			name: "All Ones",
			res: &database.ThinkResult{
				Clickbait: 1.0, Framing: 1.0, Persuasive: 1.0, HyperStimulus: 1.0, Speculative: 1.0,
			},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, calculateHybrid(tt.res, 0.7))
		})
	}
}

func TestGetRatingIcon(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{0.0, "❂"},
		{0.1, "❂"},
		{0.2, "◈"},
		{0.3, "◈"},
		{0.5, "◬"},
		{0.6, "◬"},
		{0.7, "⦸"},
		{0.8, "⦸"},
		{0.9, "☒"},
		{1.0, "☒"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, getRatingIcon(tt.score))
		})
	}
}

func TestCreateStarRating(t *testing.T) {
	tests := []struct {
		name     string
		res      *database.ThinkResult
		expected StarRating
	}{
		{
			name: "Perfect Score",
			res: &database.ThinkResult{
				Clickbait:     0.0,
				Framing:       0.0,
				Persuasive:    0.0,
				HyperStimulus: 0.0,
				Speculative:   0.0,
			},
			expected: StarRating{
				Clickbait:          "★★★★★",
				Framing:            "★★★★★",
				PersuasiveIntent:   "★★★★★",
				HyperStimulus:      "★★★★★",
				SpeculativeContent: "★★★★★",
				Overall:            "★★★★★",
			},
		},
		{
			name: "Worst Score",
			res: &database.ThinkResult{
				Clickbait:     1.0,
				Framing:       1.0,
				Persuasive:    1.0,
				HyperStimulus: 1.0,
				Speculative:   1.0,
			},
			expected: StarRating{
				Clickbait:          "☆☆☆☆☆",
				Framing:            "☆☆☆☆☆",
				PersuasiveIntent:   "☆☆☆☆☆",
				HyperStimulus:      "☆☆☆☆☆",
				SpeculativeContent: "☆☆☆☆☆",
				Overall:            "☆☆☆☆☆",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createStarRating(tt.res)
			assert.Equal(t, tt.expected, got)
		})
	}
}
