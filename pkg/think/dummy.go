package think

import (
	"math/rand"
	"time"

	"github.com/deframer/news-deframer/pkg/database"
)

type dummy struct {
	rng *rand.Rand
}

func newDummy() *dummy {
	return &dummy{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (d *dummy) Run(prompt string, language string, request Request) (*database.ThinkResult, error) {
	if _, err := getPrompt(prompt, language); err != nil {
		return nil, err
	}

	framing := d.rng.Float64()
	clickbait := d.rng.Float64()
	persuasive := d.rng.Float64()
	hyperStimulus := d.rng.Float64()
	speculative := d.rng.Float64()
	overall := (framing + clickbait + persuasive + hyperStimulus + speculative) / 5.0

	return &database.ThinkResult{
		TitleCorrected:              request.Title,
		TitleCorrectionReason:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		DescriptionCorrected:        request.Description,
		DescriptionCorrectionReason: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		Framing:                     framing,
		FramingReason:               "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		Clickbait:                   clickbait,
		ClickbaitReason:             "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
		Persuasive:                  persuasive,
		PersuasiveReason:            "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
		HyperStimulus:               hyperStimulus,
		HyperStimulusReason:         "Curabitur pretium tincidunt lacus. Nulla gravida orci a odio.",
		Speculative:                 speculative,
		SpeculativeReason:           "Nullam varius, turpis et commodo pharetra, est eros bibendum elit, nec luctus magna felis sollicitudin mauris.",
		Overall:                     overall,
		OverallReason:               "Integer in mauris eu nibh euismod gravida. Duis ac tellus et risus vulputate vehicula.",
	}, nil

	/*
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
			Overall:                     0.5,
			OverallReason:               "The text is sensationalized clickbait exaggerating routine financial news to induce panic.",
		}, nil
	*/
}
