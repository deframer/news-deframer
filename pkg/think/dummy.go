package think

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
