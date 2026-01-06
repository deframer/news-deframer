package database

import (
	"testing"
	"time"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGetTime(t *testing.T) {
	cfg, err := config.Load()
	assert.NoError(t, err)

	repo, err := NewRepository(cfg)
	assert.NoError(t, err)

	got, err := repo.GetTime()
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), got, 1*time.Minute)
}
