package facade

import (
	"context"

	"github.com/egandro/news-deframer/pkg/config"
	"github.com/egandro/news-deframer/pkg/database"
	"github.com/egandro/news-deframer/pkg/valkey"
)

type Facade struct {
	ctx    context.Context
	cfg    *config.Config
	valkey valkey.Valkey
	repo   database.Repository
}

func New(ctx context.Context, cfg *config.Config, v valkey.Valkey, repo database.Repository) *Facade {
	return &Facade{
		ctx:    ctx,
		cfg:    cfg,
		valkey: v,
		repo:   repo,
	}
}

func (f *Facade) HasFeed(ctx context.Context, url string) (bool, error) {
	return false, nil
}

func (f *Facade) HasArticle(ctx context.Context, url string) (bool, error) {
	return false, nil
}
