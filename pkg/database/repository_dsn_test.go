package database

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDSNWithApplicationName(t *testing.T) {
	tests := []struct {
		name string
		app  string
	}{
		{name: "spaces", app: "News Deframer Browser Extension Service"},
		{name: "apostrophe", app: "Worker's Mode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := postgresDSNWithApplicationName("host=localhost user=deframer dbname=deframer", tt.app)

			cfg, err := pgconn.ParseConfig(dsn)
			assert.NoError(t, err)
			assert.Equal(t, tt.app, cfg.RuntimeParams["application_name"])
		})
	}
}
