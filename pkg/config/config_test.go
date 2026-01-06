package config

import (
	"os"
	"testing"
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
		cfg := Load()

		if cfg.Port != "8080" {
			t.Errorf("expected default Port '8080', got '%s'", cfg.Port)
		}
		if cfg.ValkeyHost != "valkey:6379" {
			t.Errorf("expected default ValkeyHost 'valkey:6379', got '%s'", cfg.ValkeyHost)
		}
		if cfg.DebugLevel != "debug" {
			t.Errorf("expected default DebugLevel 'debug', got '%s'", cfg.DebugLevel)
		}
	})

	t.Run("Environment Variables Override", func(t *testing.T) {
		unsetEnv()
		_ = os.Setenv("PORT", "9090")
		_ = os.Setenv("DEBUG_LEVEL", "info")
		_ = os.Setenv("VALKEY_HOST", "127.0.0.1:6379")

		defer unsetEnv() // Cleanup

		cfg := Load()

		if cfg.Port != "9090" {
			t.Errorf("expected Port '9090', got '%s'", cfg.Port)
		}
		if cfg.DebugLevel != "info" {
			t.Errorf("expected DebugLevel 'info', got '%s'", cfg.DebugLevel)
		}
		if cfg.ValkeyHost != "127.0.0.1:6379" {
			t.Errorf("expected ValkeyHost '127.0.0.1:6379', got '%s'", cfg.ValkeyHost)
		}
	})
}
