package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Helper to unset env vars to ensure clean state for defaults
	unsetEnv := func() {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("DSN")
		_ = os.Unsetenv("LOG_LEVEL")
	}

	t.Run("Defaults", func(t *testing.T) {
		unsetEnv()
		// We assume no .env file is present in pkg/config/ during test execution
		// or that godotenv doesn't find one in the test CWD.
		cfg, err := Load()
		assert.NoError(t, err)

		assert.NotEmpty(t, cfg.Port)
		assert.NotEmpty(t, cfg.LogLevel)
		assert.NotEmpty(t, cfg.DSN)
	})

	t.Run("Environment Variables Override", func(t *testing.T) {
		unsetEnv()
		_ = os.Setenv("PORT", "9090")
		_ = os.Setenv("LOG_LEVEL", "info")

		defer unsetEnv() // Cleanup

		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, "info", cfg.LogLevel)
	})
}
