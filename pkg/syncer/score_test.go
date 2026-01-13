package syncer

import (
	"testing"

	"github.com/egandro/news-deframer/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHybrid(t *testing.T) {
	tests := []struct {
		name     string
		res      *database.ThinkResult
		expected float64
	}{
		{
			name:     "Nil Result",
			res:      nil,
			expected: 0.0,
		},
		{
			name: "All Ones",
			res: &database.ThinkResult{
				Clickbait: 1.0, Framing: 1.0, Persuasive: 1.0, HyperStimulus: 1.0, Speculative: 1.0,
			},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CalculateHybrid(tt.res, 0.7))
		})
	}
}
