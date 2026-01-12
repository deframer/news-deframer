package think

import (
	"embed"
	"fmt"
)

//go:embed prompts/*.md
var promptFS embed.FS

type LLMType int

const (
	Dummy LLMType = iota
)

// type ThinkerResult struct {
// 	Hash                     string  `json:"hash,omitempty"`     // will be set to '' when persisted
// 	Url                      string  `json:"url,omitempty"`      // will be set to '' when persisted
// 	FeedUrl                  string  `json:"feed_url,omitempty"` // will be set to '' when persisted
// 	TitleCorrected           string  `json:"title_corrected,omitempty"`
// 	DescriptionCorrected     string  `json:"description_corrected,omitempty"`
// 	ClickbaitScore           float64 `json:"clickbait,omitempty"`
// 	FramingScore             float64 `json:"framing,omitempty"`
// 	PersuasiveIntentScore    float64 `json:"persuasive_intent,omitempty"`
// 	HyperStimulusScore       float64 `json:"hyper_stimulus,omitempty"`
// 	SpeculativeContentScore  float64 `json:"speculative_content,omitempty"`
// 	ClickbaitReason          string  `json:"reason_clickbait,omitempty"`
// 	FramingReason            string  `json:"reason_framing,omitempty"`
// 	PersuasiveIntentReason   string  `json:"reason_persuasive,omitempty"`
// 	HyperStimulusReason      string  `json:"reason_stimulus,omitempty"`
// 	SpeculativeContentReason string  `json:"reason_speculative,omitempty"`
// }

type Think interface {
	Run(prompt string, language string, data map[string]interface{}) (map[string]interface{}, error)
}

func New(t LLMType) (Think, error) {
	switch t {
	case Dummy:
		return &dummy{}, nil
	default:
		return nil, fmt.Errorf("unknown think type: %v", t)
	}
}

type dummy struct{}

func (d *dummy) Run(prompt string, language string, data map[string]interface{}) (map[string]interface{}, error) {
	if _, err := getPrompt(prompt, language); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"title_corrected":               "Company X announces Q3 earnings",
		"title_correction_reason":       "Removed 'shocking' and 'disaster'; standard business reporting tone used.",
		"description_corrected":         "The company reported a 10% decline in revenue due to supply chain issues.",
		"description_correction_reason": "Removed dramatic language and focused on the reported statistics.",
		"framing":                       0.4,
		"framing_reason":                "Negative framing of standard market fluctuation.",
		"clickbait":                     0.8,
		"clickbait_reason":              "Used 'You won't believe' curiosity gap.",
		"persuasive_intent":             0.0,
		"persuasive_reason":             "No call to action detected.",
		"hyper_stimulus":                0.6,
		"hyper_stimulus_reason":         "Use of all-caps on key emotional words.",
		"speculative_content":           0.2,
		"speculative_reason":            "Implies bankruptcy without official filing source.",
		"overall_reason":                "The text is sensationalized clickbait exaggerating routine financial news to induce panic.",
	}

	if v, ok := data["title"].(string); ok {
		result["title_corrected"] = "dummy ai - " + v
	}
	if v, ok := data["description"].(string); ok {
		result["description_corrected"] = "dummy ai - " + v
	}

	return result, nil
}

func getPrompt(name, lang string) (string, error) {
	filename := fmt.Sprintf("prompts/%s-prompt-%s.md", name, lang)
	content, err := promptFS.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
