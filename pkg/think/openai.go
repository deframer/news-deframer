package think

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// OpenAI requires a standard JSON Schema for structured outputs.
var openAISchemaDefinition = jsonschema.Definition{
	Type: jsonschema.Object,
	Properties: map[string]jsonschema.Definition{
		"title_corrected":               {Type: jsonschema.String},
		"title_correction_reason":       {Type: jsonschema.String},
		"description_corrected":         {Type: jsonschema.String},
		"description_correction_reason": {Type: jsonschema.String},
		"framing":                       {Type: jsonschema.Number},
		"framing_reason":                {Type: jsonschema.String},
		"clickbait":                     {Type: jsonschema.Number},
		"clickbait_reason":              {Type: jsonschema.String},
		"persuasive":                    {Type: jsonschema.Number},
		"persuasive_reason":             {Type: jsonschema.String},
		"hyper_stimulus":                {Type: jsonschema.Number},
		"hyper_stimulus_reason":         {Type: jsonschema.String},
		"speculative":                   {Type: jsonschema.Number},
		"speculative_reason":            {Type: jsonschema.String},
		"overall":                       {Type: jsonschema.Number},
		"overall_reason":                {Type: jsonschema.String},
	},
	Required: []string{
		"title_corrected", "title_correction_reason",
		"description_corrected", "description_correction_reason",
		"framing", "framing_reason",
		"clickbait", "clickbait_reason",
		"persuasive", "persuasive_reason",
		"hyper_stimulus", "hyper_stimulus_reason",
		"speculative", "speculative_reason",
		"overall", "overall_reason",
	},
	// strictly required for OpenAI structured outputs
	AdditionalProperties: false,
}

type openaiProvider struct {
	ctx     context.Context
	model   string
	apiKey  string
	baseURL string
	client  *openai.Client
	mu      sync.RWMutex
	// Cache stores the raw system prompt text string
	cache map[string]string
}

func newOpenAI(ctx context.Context, model, apiKey, baseURL string) (*openaiProvider, error) {
	config := openai.DefaultConfig(apiKey)

	// Set the BaseURL (Crucial for LM Studio)
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	client := openai.NewClientWithConfig(config)

	return &openaiProvider{
		ctx:     ctx,
		model:   model,
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  client,
		cache:   make(map[string]string),
	}, nil
}

func (o *openaiProvider) Run(prompt string, language string, request Request) (*database.ThinkResult, error) {
	key := prompt + ":" + language

	o.mu.RLock()
	sysPromptText, ok := o.cache[key]
	o.mu.RUnlock()

	if !ok {
		var err error
		sysPromptText, err = getPrompt(prompt, language)
		if err != nil {
			return nil, err
		}

		o.mu.Lock()
		// Double-check locking
		if cached, exists := o.cache[key]; exists {
			sysPromptText = cached
		} else {
			o.cache[key] = sysPromptText
		}
		o.mu.Unlock()
	}

	var temperature float32 = 0.0

	// Prepare the schema for the API request
	schemaBytes, err := json.Marshal(openAISchemaDefinition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	start := time.Now()
	resp, err := o.client.CreateChatCompletion(
		o.ctx,
		openai.ChatCompletionRequest{
			Model:       o.model,
			Temperature: temperature,
			// Define Structured Output (JSON Schema)
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "analysis_result",
					Strict: true,
					Schema: json.RawMessage(schemaBytes),
				},
			},
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sysPromptText,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Title: %s\nDescription: %s", request.Title, request.Description),
				},
			},
		},
	)
	slog.Debug("openai request duration", "duration", time.Since(start))
	if err != nil {
		return nil, err
	}

	// Logging Usage
	if resp.Usage.TotalTokens > 0 {
		// OpenAI puts "Reasoning/Thought" tokens in CompletionTokensDetails
		// Note: Not all local models return this detail, but GPT-o1/o3 do.
		var thoughts int
		if resp.Usage.CompletionTokensDetails != nil {
			thoughts = resp.Usage.CompletionTokensDetails.ReasoningTokens
		}

		slog.Debug("openai token usage",
			"prompt_tokens", resp.Usage.PromptTokens,
			"completion_tokens", resp.Usage.CompletionTokens,
			"thoughts_tokens", thoughts,
			"total_tokens", resp.Usage.TotalTokens,
		)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from openai provider")
	}

	var result database.ThinkResult
	// OpenAI returns the result in Message.Content
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return &result, nil
}

// Models

/*
M1 16gb / RTX 3060 (16GB)

Since you have about 10-12GB of usable memory for the model (leaving room for the OS), you should look for GGUF versions (Quantization) in LM Studio, specifically Q4_K_M or Q5_K_M.

1. Qwen 2.5 14B (Instruct)
Why: Currently the "king" of mid-sized models. It follows complex formatting instructions (like strict JSON) better than almost any other model in this weight class.
Fit: You need the Q4_K_M quantization (approx. 9GB). It will fill your RTX 3060 nicely but run very efficiently.

2. Llama 3.1 8B (Instruct)
Why: The industry standard for speed and reliability. It is very strict with instructions and rarely breaks JSON formatting.
Fit: You can run higher precision versions like Q6_K or Q8_0 (8-9GB) easily on both machines for maximum smarts without speed loss.

3. Mistral NeMo 12B (Instruct)
Why: A collaboration between NVIDIA and Mistral designed specifically to fit on cards like the RTX 3060. It has a larger context window and excellent logical reasoning for analyzing text bias.
Fit: Use Q4_K_M or Q5_K_M (approx. 8-9GB).

4. Gemma 2 9B (Instruct)
Why: Google’s open model. It punches way above its weight class in reasoning and nuance, which is crucial for determining "framing" and "bias" scores.

Fit: Fits comfortably at Q5_K_M or Q6_K (approx. 7-8GB).
Recommendation: Start with Qwen 2.5 14B for the best intelligence, or Llama 3.1 8B for the fastest reliable speed.
*/
