package analyzer

type AnalyzerRequest struct {
	Hash        string `json:"hash"`
	Url         string `json:"url"`
	FeedUrl     string `json:"feed_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
}

type AnalyzerResult struct {
	Hash                     string  `json:"hash,omitempty"`     // will be set to '' when persisted
	Url                      string  `json:"url,omitempty"`      // will be set to '' when persisted
	FeedUrl                  string  `json:"feed_url,omitempty"` // will be set to '' when persisted
	TitleCorrected           string  `json:"title_corrected,omitempty"`
	DescriptionCorrected     string  `json:"description_corrected,omitempty"`
	ClickbaitScore           float64 `json:"clickbait,omitempty"`
	FramingScore             float64 `json:"framing,omitempty"`
	PersuasiveIntentScore    float64 `json:"persuasive_intent,omitempty"`
	HyperStimulusScore       float64 `json:"hyper_stimulus,omitempty"`
	SpeculativeContentScore  float64 `json:"speculative_content,omitempty"`
	ClickbaitReason          string  `json:"reason_clickbait,omitempty"`
	FramingReason            string  `json:"reason_framing,omitempty"`
	PersuasiveIntentReason   string  `json:"reason_persuasive,omitempty"`
	HyperStimulusReason      string  `json:"reason_stimulus,omitempty"`
	SpeculativeContentReason string  `json:"reason_speculative,omitempty"`
}
