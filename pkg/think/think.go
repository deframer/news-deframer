package think

import (
	"context"
	"embed"
	"fmt"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
)

//go:embed prompts/*.md
var promptFS embed.FS

type Think interface {
	Run(prompt string, language string, request Request) (*database.ThinkResult, error)
}

func New(ctx context.Context, cfg *config.Config) (Think, error) {
	t := cfg.LLM_Type
	switch t {
	case config.Dummy:
		return newDummy(), nil
	case config.Gemini:
		return newGemini(ctx, cfg.LLM_Model, cfg.LLM_APIKey)
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
