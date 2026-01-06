package config

import (
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// HTT Port
	Port string `env:"PORT" envDefault:"8080"`

	// Gorm DNS
	DSN         string `env:"DSN" envDefault:"host=postgres user=deframer password=deframer dbname=deframer port=5432 sslmode=disable"`
	LogDatabase bool   `env:"LOG_DATABASE" envDefault:"false"`

	// Valkey
	ValkeyHost     string `env:"VALKEY_HOST" envDefault:"valkey:6379"`
	ValkeyPassword string `env:"VALKEY_PASSWORD" envDefault:"deframer"`
	ValkeyDB       string `env:"VALKEY_DB" envDefault:"0"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
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

	return cfg, nil
}
