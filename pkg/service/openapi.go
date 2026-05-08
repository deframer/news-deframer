package service

import (
	"context"
	"fmt"

	openapi "github.com/deframer/news-deframer/gen/openapi"
	"github.com/deframer/news-deframer/pkg/openapiutil"
	"goa.design/clue/log"
)

// openapi service example implementation.
// The example methods log the requests and return zero values.
type openapisrvc struct{}

// NewOpenapi returns the openapi service implementation.
func NewOpenapi() openapi.Service {
	return &openapisrvc{}
}

// OpenAPI implements openApi.
func (s *openapisrvc) OpenAPI(ctx context.Context, p *openapi.OpenAPIPayload) (res string, err error) {
	log.Printf(ctx, "openapi.openApi")

	url := "http://localhost:8000"
	if p.ForwardedProto != nil && p.ForwardedHost != nil {
		// this might be only available if we run in Traefik/Kubernetes
		url = fmt.Sprintf("%s://%s", *p.ForwardedProto, *p.ForwardedHost)
	}

	ctx = log.With(ctx,
		log.KV{K: "url", V: url},
		log.KV{K: "tag", V: *p.Tag},
	)
	log.Printf(ctx, "creating openapi")

	filterPaths := []string{"/ping", "/openapi", "/openapi/{tag}"}
	json, err := openapiutil.OpenAPI(url, *p.Tag, &filterPaths, nil, "./gen/http/openapi3.json")
	if err != nil {
		return "", openapi.MakeInvalidOpenapiFile(err)
	}

	return json, nil
}

// OpenAPIOnlyFor implements openApiOnlyFor.
func (s *openapisrvc) OpenAPIOnlyFor(ctx context.Context, p *openapi.OpenAPIOnlyForPayload) (res string, err error) {
	log.Printf(ctx, "openapi.openApiOnlyFor")

	url := "http://localhost:8000"
	if p.ForwardedProto != nil && p.ForwardedHost != nil {
		// this might be only available if we run in Traefik/Kubernetes
		url = fmt.Sprintf("%s://%s", *p.ForwardedProto, *p.ForwardedHost)
	}

	ctx = log.With(ctx,
		log.KV{K: "url", V: url},
		log.KV{K: "tag", V: *p.Tag},
	)
	log.Printf(ctx, "creating openapi")

	filterPaths := []string{"/ping", "/openapi", "/openapi/{tag}"}
	json, err := openapiutil.OpenAPI(url, *p.Tag, &filterPaths, &p.Paths, "./gen/http/openapi3.json")
	if err != nil {
		return "", openapi.MakeInvalidOpenapiFile(err)
	}

	return json, nil
}
