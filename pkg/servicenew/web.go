package servicenew

import (
	"context"

	web "github.com/deframer/news-deframer/gen/web"
	sharedservice "github.com/deframer/news-deframer/pkg/service"
)

// NewWeb returns the web service implementation.
func NewWeb(ctx context.Context) web.Service {
	return sharedservice.NewService(ctx)
}
