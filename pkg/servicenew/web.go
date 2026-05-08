package servicenew

import (
	web "github.com/deframer/news-deframer/gen/web"
	sharedservice "github.com/deframer/news-deframer/pkg/service"
)

// NewWeb returns the web service implementation.
func NewWeb() web.Service {
	return sharedservice.NewService()
}
