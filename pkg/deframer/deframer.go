package deframer

import (
	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/openai"
	"github.com/egandro/news-deframer/pkg/source"
)

type Deframer struct {
	db  *database.Database
	ai  *openai.OpenAI
	src *source.Source
}

// NewDeframer initializes a new deframer
func NewDeframer() (*Deframer, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(cfg.DatabaseFile)
	if err != nil {
		return nil, err
	}

	ai := openai.NewAI(cfg.AI_URL, cfg.AI_Model, "")

	src, err := source.ParseFile(cfg.Source)
	if err != nil {
		return nil, err
	}

	res := &Deframer{
		db:  db,
		ai:  ai,
		src: src,
	}

	return res, nil
}
