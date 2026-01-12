package think

import (
	"embed"
	"fmt"

	"github.com/egandro/news-deframer/pkg/config"
)

//go:embed prompts/*.md
var promptFS embed.FS

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
