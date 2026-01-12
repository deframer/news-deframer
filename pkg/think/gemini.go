package think

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/genai"
)

type gemini struct {
	ctx    context.Context
	model  string
	apiKey string
}

func (g *gemini) Run(prompt string, language string, request Request) (*Result, error) {
	sysPrompt, err := getPrompt(prompt, language)
	if err != nil {
		return nil, err
	}

	client, err := genai.NewClient(g.ctx, &genai.ClientConfig{
		APIKey:  g.apiKey,
		Backend: genai.BackendGeminiAPI, // Use VertexAI if you are on Google Cloud
	})
	if err != nil {
		return nil, err
	}

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title_corrected":               {Type: genai.TypeString},
			"title_correction_reason":       {Type: genai.TypeString},
			"description_corrected":         {Type: genai.TypeString},
			"description_correction_reason": {Type: genai.TypeString},
			"framing":                       {Type: genai.TypeNumber},
			"framing_reason":                {Type: genai.TypeString},
			"clickbait":                     {Type: genai.TypeNumber},
			"clickbait_reason":              {Type: genai.TypeString},
			"persuasive_intent":             {Type: genai.TypeNumber},
			"persuasive_reason":             {Type: genai.TypeString},
			"hyper_stimulus":                {Type: genai.TypeNumber},
			"hyper_stimulus_reason":         {Type: genai.TypeString},
			"speculative_content":           {Type: genai.TypeNumber},
			"speculative_reason":            {Type: genai.TypeString},
			"overall_reason":                {Type: genai.TypeString},
		},
		Required: []string{
			"title_corrected", "title_correction_reason",
			"description_corrected", "description_correction_reason",
			"framing", "framing_reason",
			"clickbait", "clickbait_reason",
			"persuasive_intent", "persuasive_reason",
			"hyper_stimulus", "hyper_stimulus_reason",
			"speculative_content", "speculative_reason",
			"overall_reason",
		},
	}

	var temperature float32 = 0.0

	resp, err := client.Models.GenerateContent(g.ctx, g.model,
		genai.Text(fmt.Sprintf("Title: %s\nDescription: %s", request.Title, request.Description)),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   schema,
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{
					{Text: sysPrompt},
				},
			},
			Temperature: &temperature,
		},
	)
	if err != nil {
		return nil, err
	}

	var result Result
	if err := json.Unmarshal([]byte(resp.Candidates[0].Content.Parts[0].Text), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
