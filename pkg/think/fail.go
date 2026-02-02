package think

import (
	"fmt"

	"github.com/deframer/news-deframer/pkg/database"
)

type fail struct{}

func newFail() *fail {
	return &fail{}
}

func (f *fail) Run(prompt string, language string, request Request) (*database.ThinkResult, error) {
	return nil, fmt.Errorf("intentionally failed")
}
