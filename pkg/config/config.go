package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// HTT Port
	Port string `env:"PORT" envDefault:"8080"`

	// Gorm DNS
	DSN string `env:"DSN" envDefault:"host=postgres user=deframer password=deframer dbname=deframer port=5432 sslmode=disable"`

	// Valkey
	ValkeyHost     string `env:"VALKEY_HOST" envDefault:"valkey:6379"`
	ValkeyPassword string `env:"VALKEY_PASSWORD" envDefault:"deframer"`
	ValkeyDB       string `env:"VALKEY_DB" envDefault:"0"`

	DebugLevel string `env:"DEBUG_LEVEL" envDefault:"debug"`
}

func Load() *Config {
	// 1. Load .env file (if exists)
	// We ignore the error because in production we might rely solely on real Env Vars
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("Config parse error: %v\n", err)
	}

	return cfg
}
