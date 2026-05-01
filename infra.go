package servicenew

import (
	"context"

	infra "github.com/deframer/news-deframer/gen/infra"
	"goa.design/clue/log"
)

// infra service example implementation.
// The example methods log the requests and return zero values.
type infrasrvc struct{}

// NewInfra returns the infra service implementation.
func NewInfra() infra.Service {
	return &infrasrvc{}
}

// Health check endpoint.
func (s *infrasrvc) Ping(ctx context.Context) (res string, err error) {
	log.Printf(ctx, "infra.ping")
	return
}

// Return the current host name.
func (s *infrasrvc) Hostname(ctx context.Context) (res *infra.HostnameResponse, err error) {
	res = &infra.HostnameResponse{}
	log.Printf(ctx, "infra.hostname")
	return
}
