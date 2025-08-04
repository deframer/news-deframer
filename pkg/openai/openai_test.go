package openai

import (
	"context"
	"testing"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func TestQuery(t *testing.T) {
	t.Skip("this requires a LM-Studio connection")
	ctx := context.Background()

	cfg, err := config.GetConfig()
	assert.NoError(t, err)

	const user = `
		Create 10 cat names.
		Strictly return it as json array. Sample: ["name1", "name2"]`
	const system = "You are a cat name expert."

	ai := NewAI(cfg.AI_URL, cfg.AI_Model, "")
	assert.NotNil(t, ai)

	res, err := ai.Query(ctx, user, system)
	assert.NoError(t, err)
	assert.NotEqual(t, res, "")
}

func TestFuzzyParseJSON(t *testing.T) {
	input := `json

	[
		"Whiskers Shadow", "Mittens Blaze", "Shadow Purrfect", "Midnight Whisker", "Aurora Claw", "Sapphire Meowster", "Basil Flufftail", "Glimmer Paw", "Twilight Velvet", "Cocoa Munchkin"
	]
`
	ai := NewAI("", "", "")
	parsed, err := ai.FuzzyParseJSON(input)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)
}
