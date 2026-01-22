package syncer

import (
	"testing"

	"github.com/deframer/news-deframer/pkg/database"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestApplyFancyRatingText(t *testing.T) {
	tests := []struct {
		name                 string
		item                 *gofeed.Item
		res                  *database.ThinkResult
		thinkRating          float64
		language             string
		expectedTitle        string
		expectedDescContains []string
	}{
		{
			name: "Basic English",
			item: &gofeed.Item{
				Title:       "Corrected Title",
				Description: "Original Description",
			},
			res: &database.ThinkResult{
				Clickbait:           0.0,
				ClickbaitReason:     "None",
				Framing:             0.0,
				FramingReason:       "None",
				Persuasive:          0.0,
				PersuasiveReason:    "None",
				HyperStimulus:       0.0,
				HyperStimulusReason: "None",
				Speculative:         0.0,
				SpeculativeReason:   "None",
			},
			thinkRating:   0.0,
			language:      "en",
			expectedTitle: "★★★★★ Corrected Title",
			expectedDescContains: []string{
				"Original Description",
			},
		},
		{
			name: "Nil Item",
			item: nil,
			res:  &database.ThinkResult{},
		},
		{
			name: "Nil Result",
			item: &gofeed.Item{},
			res:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyFancyRatingText(tt.item, tt.res, tt.thinkRating, tt.language)
			if tt.item != nil && tt.res != nil {
				assert.Equal(t, tt.expectedTitle, tt.item.Title)
				for _, s := range tt.expectedDescContains {
					assert.Contains(t, tt.item.Description, s)
				}
			}
		})
	}
}
