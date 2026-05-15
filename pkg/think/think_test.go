package think

import (
	"context"
	"testing"

	categorypkg "github.com/deframer/news-deframer/pkg/category"
	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestGetPrompt(t *testing.T) {
	p, err := getPrompt("deframer", "en")
	assert.NoError(t, err)
	assert.Contains(t, p, "System Prompt")

	_, err = getPrompt("foo", "bar")
	assert.Error(t, err)
}

func TestLocalizedCategoryValidation(t *testing.T) {
	err := categorypkg.ValidateLocalizedCategory("en", "business")
	assert.NoError(t, err)

	err = categorypkg.ValidateLocalizedCategory("de", "wirtschaft")
	assert.NoError(t, err)

	err = categorypkg.ValidateLocalizedCategory("da", "erhverv")
	assert.NoError(t, err)

	err = categorypkg.ValidateLocalizedCategory("en", "wirtschaft")
	assert.Error(t, err)

	err = categorypkg.ValidateLocalizedCategory("xx", "business")
	assert.Error(t, err)
}

func TestNormalizeCategory(t *testing.T) {
	value, err := categorypkg.NormalizeCategory("en", "opinion")
	assert.NoError(t, err)
	assert.Equal(t, "opinion", value)

	value, err = categorypkg.NormalizeCategory("de", "meinung")
	assert.NoError(t, err)
	assert.Equal(t, "opinion", value)

	value, err = categorypkg.NormalizeCategory("da", "erhverv")
	assert.NoError(t, err)
	assert.Equal(t, "business", value)

	_, err = categorypkg.NormalizeCategory("de", "business")
	assert.Error(t, err)

	_, err = categorypkg.NormalizeCategory("fr", "opinion")
	assert.Error(t, err)
}

func TestDummy_Run(t *testing.T) {
	d := newDummy()

	_, err := d.Run("deframer", "en", Request{}, false)
	assert.NoError(t, err)

	_, err = d.Run("deframer", "fr", Request{}, false)
	assert.Error(t, err)
}

func TestFail_Run(t *testing.T) {
	f := newFail()

	_, err := f.Run("deframer", "en", Request{}, false)
	assert.Error(t, err)
	assert.EqualError(t, err, "intentionally failed")
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

	resp, err := g.Run("deframer", "en", req, false)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestNew_OpenAI(t *testing.T) {
	t.Skip("skipping test")

	cfg, err := config.Load()
	assert.NoError(t, err)

	if cfg.LLM_Type != config.OpenAI {
		t.Skip("skipping openai test")
	}

	th, err := New(context.Background(), cfg)
	assert.NoError(t, err)

	o, ok := th.(*openaiProvider)
	assert.True(t, ok)
	assert.Equal(t, cfg.LLM_Model, o.model)
	assert.Equal(t, cfg.LLM_APIKey, o.apiKey)
	//assert.Equal(t, cfg.LLM_BaseURL, o.baseURL)

	req := Request{
		Title:       "THE END IS HERE?! SECRET LEAK Exposes The ELITES' Plan To WIPE OUT Your Savings By TUESDAY!!",
		Description: "THEY are coming for you! A terrifying whistleblower tape PROVES the \"Mars Attack\" is starting! YOU will own NOTHING! Don't be a sheep—CLICK HERE to secure the ONLY asset that survives the crash before it’s ILLEGAL! ACT NOW OR PERISH!",
	}

	resp, err := o.Run("deframer", "en", req, false)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestVerifyThinkResult(t *testing.T) {
	err := verifyThinkResult("en", nil, false)
	assert.NoError(t, err)

	err = verifyThinkResult("en", &database.ThinkResult{Category: "other"}, false)
	assert.NoError(t, err)

	res := &database.ThinkResult{Category: "meinung"}
	err = verifyThinkResult("de", res, false)
	assert.NoError(t, err)
	assert.Equal(t, "opinion", res.Category)

	err = verifyThinkResult("en", &database.ThinkResult{
		Framing:       0.5,
		Clickbait:     0.5,
		Persuasive:    0.5,
		HyperStimulus: 0.5,
		Speculative:   0.5,
		Overall:       0.5,
		Category:      "business",
	}, false)
	assert.NoError(t, err)

	err = verifyThinkResult("en", &database.ThinkResult{Framing: -0.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Framing")

	err = verifyThinkResult("en", &database.ThinkResult{Framing: 1.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Framing")

	err = verifyThinkResult("en", &database.ThinkResult{Clickbait: -0.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Clickbait")

	err = verifyThinkResult("en", &database.ThinkResult{Persuasive: 1.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Persuasive")

	err = verifyThinkResult("en", &database.ThinkResult{HyperStimulus: -0.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "HyperStimulus")

	err = verifyThinkResult("en", &database.ThinkResult{Speculative: 1.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Speculative")

	err = verifyThinkResult("en", &database.ThinkResult{Overall: -0.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Overall")

	err = verifyThinkResult("en", &database.ThinkResult{Overall: 1.1, Category: "business"}, false)
	assert.ErrorContains(t, err, "Overall")

	err = verifyThinkResult("en", &database.ThinkResult{Category: "wirtschaft"}, false)
	assert.ErrorContains(t, err, "invalid category")

	err = verifyThinkResult("fr", &database.ThinkResult{Category: "opinion"}, false)
	assert.ErrorContains(t, err, "no localized categories configured")

	res = &database.ThinkResult{Category: "wirtschaft"}
	err = verifyThinkResult("en", res, true)
	assert.NoError(t, err)
	assert.Equal(t, "other", res.Category)
}
