package servicenew

import (
	"context"
	"log/slog"
	"os"

	infra "github.com/deframer/news-deframer/gen/infra"
)

type infrasrvc struct {
	logger *slog.Logger
}

// NewInfra returns the infra service implementation.
func NewInfra() infra.Service {
	return &infrasrvc{logger: slog.With("component", "infra")}
}

func (s *infrasrvc) Ping(ctx context.Context) (res string, err error) {
	s.logger.Debug("ping")
	return "pong", nil
}

func (s *infrasrvc) Hostname(ctx context.Context) (res *infra.HostnameResponse, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		s.logger.Error("failed to get hostname", "error", err)
		return nil, err
	}
	return &infra.HostnameResponse{Hostname: hostname}, nil
}
