package think

import (
	"crypto/rand"
	"encoding/binary"

	categorypkg "github.com/deframer/news-deframer/pkg/category"
	"github.com/deframer/news-deframer/pkg/database"
)

type dummy struct{}

func newDummy() *dummy {
	return &dummy{}
}

func (d *dummy) Run(prompt string, language string, request Request, ignoreCategoryErrors bool) (*database.ThinkResult, error) {
	if _, err := categorypkg.LocalizedCategoriesFor(language); err != nil {
		return nil, err
	}

	if _, err := getPrompt(prompt, language); err != nil {
		return nil, err
	}

	framing := secureFloat64()
	clickbait := secureFloat64()
	persuasive := secureFloat64()
	hyperStimulus := secureFloat64()
	speculative := secureFloat64()
	overall := (framing + clickbait + persuasive + hyperStimulus + speculative) / 5.0

	result := &database.ThinkResult{
		TitleCorrected:              request.Title,
		TitleCorrectionReason:       "Dummy output, not AI-generated.",
		DescriptionCorrected:        request.Description,
		DescriptionCorrectionReason: "Dummy output, not AI-generated.",
		Framing:                     framing,
		FramingReason:               "Dummy output, not AI-generated.",
		Clickbait:                   clickbait,
		ClickbaitReason:             "Dummy output, not AI-generated.",
		Persuasive:                  persuasive,
		PersuasiveReason:            "Dummy output, not AI-generated.",
		HyperStimulus:               hyperStimulus,
		HyperStimulusReason:         "Dummy output, not AI-generated.",
		Speculative:                 speculative,
		SpeculativeReason:           "Dummy output, not AI-generated.",
		Overall:                     overall,
		OverallReason:               "Dummy output, not AI-generated.",
	}

	category, err := categorypkg.FirstLocalizedCategory(language)
	if err != nil {
		return nil, err
	}
	result.Category = category

	if err := verifyThinkResult(language, result, ignoreCategoryErrors); err != nil {
		return nil, err
	}

	return result, nil

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
