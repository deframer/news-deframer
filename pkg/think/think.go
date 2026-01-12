package think

import (
	"embed"
	"fmt"

	"github.com/egandro/news-deframer/pkg/config"
)

//go:embed prompts/*.md
var promptFS embed.FS

type Think interface {
	Run(prompt string, language string, request Request) (*Result, error)
}

func New(cfg *config.Config) (Think, error) {
	t := cfg.LLM_Type
	switch t {
	case config.Dummy:
		return &dummy{}, nil
	default:
		return nil, fmt.Errorf("unknown think type: %v", t)
	}
}

func getPrompt(name, lang string) (string, error) {
	filename := fmt.Sprintf("prompts/%s-prompt-%s.md", name, lang)
	content, err := promptFS.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

type Request struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Result struct {
	TitleCorrected              string  `json:"title_corrected,omitempty"`
	TitleCorrectionReason       string  `json:"title_correction_reason,omitempty"`
	DescriptionCorrected        string  `json:"description_corrected,omitempty"`
	DescriptionCorrectionReason string  `json:"description_correction_reason,omitempty"`
	ClickbaitScore              float64 `json:"clickbait,omitempty"`
	ClickbaitReason             string  `json:"clickbait_reason,omitempty"`
	FramingScore                float64 `json:"framing,omitempty"`
	FramingReason               string  `json:"framing_reason,omitempty"`
	PersuasiveIntentScore       float64 `json:"persuasive_intent,omitempty"`
	PersuasiveReason            string  `json:"persuasive_reason,omitempty"`
	HyperStimulusScore          float64 `json:"hyper_stimulus,omitempty"`
	HyperStimulusReason         string  `json:"hyper_stimulus_reason,omitempty"`
	SpeculativeContentScore     float64 `json:"speculative_content,omitempty"`
	SpeculativeReason           string  `json:"speculative_reason,omitempty"`
	OverallReason               string  `json:"overall_reason,omitempty"`
}
