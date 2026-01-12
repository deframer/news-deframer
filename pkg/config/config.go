package config

import (
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/egandro/news-deframer/pkg/think"
	"github.com/joho/godotenv"
)

type Config struct {
	// HTT Port
	Port string `env:"PORT" envDefault:"8080"`

	// Gorm DNS
	DSN         string `env:"DSN" envDefault:"host=postgres user=deframer password=deframer dbname=deframer port=5432 sslmode=disable"`
	LogDatabase bool   `env:"LOG_DATABASE" envDefault:"false"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`

	LocalFeedFilesDir string `env:"LOCAL_FEED_FILES_DIR" envDefault:""`

	LLMType think.LLMType
}

func Load() (*Config, error) {
	// Load .env file (if exists)
	// We ignore the error because in production we rely solely on real Env Vars
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		if root := os.Getenv("PROJECT_ROOT"); root != "" {
			// running in visual studio code test
			_ = godotenv.Load(filepath.Join(root, ".env"))
		} else {
			// running non docker / cli
			_ = godotenv.Load()
		}
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// If LocalFileRoot is set, verify it exists. If not, disable it.
	if cfg.LocalFeedFilesDir != "" {
		if _, err := os.Stat(cfg.LocalFeedFilesDir); os.IsNotExist(err) {
			cfg.LocalFeedFilesDir = ""
		} else if abs, err := filepath.Abs(cfg.LocalFeedFilesDir); err == nil {
			// Store absolute path to ensure safe scoping later
			cfg.LocalFeedFilesDir = abs
		}
	}

	return cfg, nil
}
