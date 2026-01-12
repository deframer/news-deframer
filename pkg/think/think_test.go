package think

import (
	"testing"

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
	d := &dummy{}

	_, err := d.Run("deframer", "en", nil)
	assert.NoError(t, err)

	_, err = d.Run("unknown", "en", nil)
	assert.Error(t, err)
}
