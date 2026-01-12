package think

import (
	"context"
	"testing"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetPrompt(t *testing.T) {
	p, err := getPrompt("deframer", "en")
	assert.NoError(t, err)
	assert.Contains(t, p, "System Prompt")

	_, err = getPrompt("foo", "bar")
	assert.Error(t, err)
}

func TestDummy_Run(t *testing.T) {
	d := newDummy()

	_, err := d.Run("deframer", "en", Request{})
	assert.NoError(t, err)

	_, err = d.Run("unknown", "en", Request{})
	assert.Error(t, err)
}

func TestNew_Gemini(t *testing.T) {
	t.Skip("skipping test")

	cfg, err := config.Load()
	assert.NoError(t, err)

	if cfg.LLM_Type != config.Gemini {
		t.Skip("skipping gemini test")
	}

	th, err := New(context.Background(), cfg)
	assert.NoError(t, err)

	g, ok := th.(*gemini)
	assert.True(t, ok)
	assert.Equal(t, cfg.LLM_Model, g.model)
	assert.Equal(t, cfg.LLM_APIKey, g.apiKey)

	req := Request{
		Title:       "THE END IS HERE?! SECRET LEAK Exposes The ELITES' Plan To WIPE OUT Your Savings By TUESDAY!!",
		Description: "THEY are coming for you! A terrifying whistleblower tape PROVES the \"Mars Attack\" is starting! YOU will own NOTHING! Don't be a sheep—CLICK HERE to secure the ONLY asset that survives the crash before it’s ILLEGAL! ACT NOW OR PERISH!",
	}

	resp, err := g.Run("deframer", "en", req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
