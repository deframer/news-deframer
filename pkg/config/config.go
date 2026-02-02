package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type LLMType int

const (
	// PollingInterval defines how often a single feed is re-synced.
	// If a feed was synced at T, it will be eligible again at T + PollingInterval.
	PollingInterval = 10 * time.Minute

	// IdleSleepTime defines how long the worker sleeps when no feeds are due for syncing.
	IdleSleepTime = 10 * time.Second

	// ETagTTL defines the cache duration for API endpoints.
	// This is set below half the PollingInterval (the Nyquist limit) to ensure clients never miss an update.
	ETagTTL = 150 * time.Second

	// lock duration
	DefaultLockDuration = 5 * time.Minute
)

const (
	Dummy LLMType = iota
	Gemini
	OpenAI
)

func (t *LLMType) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "dummy":
		*t = Dummy
	case "gemini":
		*t = Gemini
	case "openai":
		*t = OpenAI
	default:
		return fmt.Errorf("unknown LLM type: %s", string(text))
	}
	return nil
}

type Config struct {
	// HTT Port
	Port string `env:"PORT" envDefault:"8080"`

	// Gorm DNS
	DSN         string `env:"DSN" envDefault:"host=postgres user=deframer password=deframer dbname=deframer port=5432 sslmode=disable"`
	LogDatabase bool   `env:"LOG_DATABASE" envDefault:"false"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`

	BasicAuthUser     string `env:"BASIC_AUTH_USER" envDefault:""`
	BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD" envDefault:""`

	DisableETag bool `env:"DISABLE_ETAG" envDefault:"false"`

	LLM_Type    LLMType `env:"LLM_TYPE" envDefault:"dummy"`
	LLM_Model   string  `env:"LLM_MODEL" envDefault:""`
	LLM_APIKey  string  `env:"LLM_API_KEY" envDefault:""`
	LLM_BaseURL string  `env:"LLM_BASE_URL" envDefault:""`
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
