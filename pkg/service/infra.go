package service

import (
	"context"
	"os"

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

func (s *infrasrvc) Hostname(ctx context.Context) (res *infra.HostnameResponse, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf(ctx, err, "failed to get hostname")
		return nil, err
	}
	log.Printf(ctx, "infra.hostname")
	return &infra.HostnameResponse{Hostname: hostname}, nil
}
