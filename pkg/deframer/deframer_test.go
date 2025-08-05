package deframer

import (
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/stretchr/testify/assert"
)

func setupTestDeframer(t *testing.T) (*Deframer, error) {
	// Use in-memory SQLite for testing
	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	res := &Deframer{
		db: db,
	}

	return res, nil
}

func TestNewDeframer(t *testing.T) {
	s, err := setupTestDeframer(t)
	//s, err := NewDeframer()
	assert.NoError(t, err)
	assert.NotNil(t, s, "Deframer should be initialized")
}
