package think

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/egandro/news-deframer/pkg/database"
	"google.golang.org/genai"
)

var geminiSchema = &genai.Schema{
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

type gemini struct {
	ctx    context.Context
	model  string
	apiKey string
	client *genai.Client
	mu     sync.RWMutex
	cache  map[string]*genai.Content
}

func newGemini(ctx context.Context, model, apiKey string) (*gemini, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI, // Use VertexAI if you are on Google Cloud
	})
	if err != nil {
		return nil, err
	}
	return &gemini{
		ctx:    ctx,
		model:  model,
		apiKey: apiKey,
		client: client,
		cache:  make(map[string]*genai.Content),
	}, nil
}

func (g *gemini) Run(prompt string, language string, request Request) (*database.ThinkResult, error) {
	key := prompt + ":" + language

	g.mu.RLock()
	sysInstruction, ok := g.cache[key]
	g.mu.RUnlock()

	if !ok {
		sysPrompt, err := getPrompt(prompt, language)
		if err != nil {
			return nil, err
		}
		sysInstruction = &genai.Content{
			Parts: []*genai.Part{
				{Text: sysPrompt},
			},
		}
		g.mu.Lock()
		// Double-check locking: ensure another goroutine didn't populate the cache
		// while we were loading the file and creating the struct.
		if cached, exists := g.cache[key]; exists {
			sysInstruction = cached
		} else {
			g.cache[key] = sysInstruction
		}
		g.mu.Unlock()
	}

	var temperature float32 = 0.0

	userPrompt := genai.Text(fmt.Sprintf("Title: %s\nDescription: %s", request.Title, request.Description))

	// While Gemini does offer a Context Caching feature (which stores prompts on the server to avoid re-transmission),
	// it currently has a minimum requirement of 32,768 tokens (roughly 25,000 words). The system prompt is significantly
	// smaller than this threshold, so using Context Caching is not applicable or cost-effective for this specific use case.
	// Sending the text every time is the correct approach here.

	start := time.Now()
	resp, err := g.client.Models.GenerateContent(g.ctx, g.model,
		userPrompt,
		&genai.GenerateContentConfig{
			ResponseMIMEType:  "application/json",
			ResponseSchema:    geminiSchema,
			SystemInstruction: sysInstruction,
			Temperature:       &temperature,
		},
	)
	slog.Debug("gemini request duration", "duration", time.Since(start))
	if err != nil {
		return nil, err
	}

	if resp.UsageMetadata != nil {
		var thoughts int32
		if resp.UsageMetadata.ThoughtsTokenCount != 0 {
			thoughts = resp.UsageMetadata.ThoughtsTokenCount
		}

		slog.Debug("gemini token usage",
			"prompt_tokens", resp.UsageMetadata.PromptTokenCount,
			"candidate_tokens", resp.UsageMetadata.CandidatesTokenCount,
			"thoughts_tokens", thoughts,
			"total_tokens", resp.UsageMetadata.TotalTokenCount,
		)
	}

	var result database.ThinkResult
	if err := json.Unmarshal([]byte(resp.Candidates[0].Content.Parts[0].Text), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
