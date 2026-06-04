package service

import (
	"context"

	infra "github.com/deframer/news-deframer/gen/infra"
	"goa.design/clue/log"
)

type infrasrvc struct {
}

// NewInfra returns the infra service implementation.
func NewInfra() infra.Service {
	return &infrasrvc{}
}

func (s *infrasrvc) Ping(ctx context.Context) (res string, err error) {
	log.Printf(ctx, "infra.ping")
	return "pong", nil
}
