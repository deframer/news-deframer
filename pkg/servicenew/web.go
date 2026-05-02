package servicenew

import (
	"context"

	web "github.com/deframer/news-deframer/gen/web"
	sharedservice "github.com/deframer/news-deframer/pkg/service"
)

type websrvc struct {
	svc web.Service
}

// NewWeb returns the web service implementation.
func NewWeb() web.Service {
	return &websrvc{svc: sharedservice.NewService()}
}

func (s *websrvc) Item(ctx context.Context, p *web.ItemPayload) (res *web.AnalyzedItem, err error) {
	return s.svc.Item(ctx, p)
}

func (s *websrvc) Site(ctx context.Context, p *web.SitePayload) (res []*web.AnalyzedItem, err error) {
	return s.svc.Site(ctx, p)
}

func (s *websrvc) Articles(ctx context.Context, p *web.ArticlesPayload) (res []*web.AnalyzedArticle, err error) {
	return s.svc.Articles(ctx, p)
}

func (s *websrvc) Sentiments(ctx context.Context, p *web.SentimentsPayload) (res *web.SentimentItem, err error) {
	return s.svc.Sentiments(ctx, p)
}

func (s *websrvc) Domains(ctx context.Context, p *web.DomainsPayload) (res []*web.DomainEntry, err error) {
	return s.svc.Domains(ctx, p)
}

func (s *websrvc) TopTrendsByDomain(ctx context.Context, p *web.TopTrendsByDomainPayload) (res []*web.TrendMetric, err error) {
	return s.svc.TopTrendsByDomain(ctx, p)
}

func (s *websrvc) ContextByDomain(ctx context.Context, p *web.ContextByDomainPayload) (res []*web.TrendContext, err error) {
	return s.svc.ContextByDomain(ctx, p)
}

func (s *websrvc) LifecycleByDomain(ctx context.Context, p *web.LifecycleByDomainPayload) (res []*web.Lifecycle, err error) {
	return s.svc.LifecycleByDomain(ctx, p)
}

func (s *websrvc) DomainComparisonEndpoint(ctx context.Context, p *web.DomainComparisonPayload) (res []*web.DomainComparison, err error) {
	return s.svc.DomainComparisonEndpoint(ctx, p)
}
