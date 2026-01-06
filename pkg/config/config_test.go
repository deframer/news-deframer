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
		_ = os.Unsetenv("VALKEY_HOST")
		_ = os.Unsetenv("VALKEY_PASSWORD")
		_ = os.Unsetenv("VALKEY_DB")
		_ = os.Unsetenv("DEBUG_LEVEL")
	}

	t.Run("Defaults", func(t *testing.T) {
		unsetEnv()
		// We assume no .env file is present in pkg/config/ during test execution
		// or that godotenv doesn't find one in the test CWD.
		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "valkey:6379", cfg.ValkeyHost)
		assert.Equal(t, "debug", cfg.DebugLevel)
	})

	t.Run("Environment Variables Override", func(t *testing.T) {
		unsetEnv()
		_ = os.Setenv("PORT", "9090")
		_ = os.Setenv("DEBUG_LEVEL", "info")
		_ = os.Setenv("VALKEY_HOST", "127.0.0.1:6379")

		defer unsetEnv() // Cleanup

		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, "info", cfg.DebugLevel)
		assert.Equal(t, "127.0.0.1:6379", cfg.ValkeyHost)
	})
}
