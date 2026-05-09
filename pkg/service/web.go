package service

import (
	"context"

	web "github.com/deframer/news-deframer/gen/web"
)

// NewWeb returns the web service implementation.
func NewWeb(ctx context.Context) web.Service {
	return NewWebImplementation(ctx)
}
