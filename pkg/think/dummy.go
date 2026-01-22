package think

import (
	"github.com/deframer/news-deframer/pkg/database"
)

type dummy struct{}

func newDummy() *dummy {
	return &dummy{}
}

func (d *dummy) Run(prompt string, language string, request Request) (*database.ThinkResult, error) {
	if _, err := getPrompt(prompt, language); err != nil {
		return nil, err
	}

	return &database.ThinkResult{
		TitleCorrected:              "from ai: " + request.Title,
		TitleCorrectionReason:       "Removed 'shocking' and 'disaster'; standard business reporting tone used.",
		DescriptionCorrected:        "from ai: " + request.Description,
		DescriptionCorrectionReason: "Removed dramatic language and focused on the reported statistics.",
		Framing:                     0.4,
		FramingReason:               "Negative framing of standard market fluctuation.",
		Clickbait:                   0.8,
		ClickbaitReason:             "Used 'You won't believe' curiosity gap.",
		Persuasive:                  0.0,
		PersuasiveReason:            "No call to action detected.",
		HyperStimulus:               0.6,
		HyperStimulusReason:         "Use of all-caps on key emotional words.",
		Speculative:                 0.2,
		SpeculativeReason:           "Implies bankruptcy without official filing source.",
		OverallReason:               "The text is sensationalized clickbait exaggerating routine financial news to induce panic.",
	}, nil
}
