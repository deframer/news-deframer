package think

import (
	"crypto/rand"
	"encoding/binary"

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

	framing := secureFloat64()
	clickbait := secureFloat64()
	persuasive := secureFloat64()
	hyperStimulus := secureFloat64()
	speculative := secureFloat64()
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

func secureFloat64() float64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return float64(binary.LittleEndian.Uint64(b[:])>>11) / float64(1<<53)
}
