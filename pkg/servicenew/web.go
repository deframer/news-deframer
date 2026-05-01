package servicenew

import (
	"context"

	web "github.com/deframer/news-deframer/gen/web"
	"goa.design/clue/log"
	"goa.design/goa/v3/security"
)

// web service example implementation.
// The example methods log the requests and return zero values.
type websrvc struct{}

// NewWeb returns the web service implementation.
func NewWeb() web.Service {
	return &websrvc{}
}

// Fetch a single analyzed item by URL.
func (s *websrvc) Item(ctx context.Context, p *web.ItemPayload) (res *web.AnalyzedItem, err error) {
	res = &web.AnalyzedItem{}
	log.Printf(ctx, "web.item")
	return
}

// List analyzed items for a root domain.
func (s *websrvc) Site(ctx context.Context, p *web.SitePayload) (res []*web.AnalyzedItem, err error) {
	log.Printf(ctx, "web.site")
	return
}

// List articles for a trend and domain.
func (s *websrvc) Articles(ctx context.Context, p *web.ArticlesPayload) (res []*web.AnalyzedArticle, err error) {
	log.Printf(ctx, "web.articles")
	return
}

// Get sentiment scores for a trend.
func (s *websrvc) Sentiments(ctx context.Context, p *web.SentimentsPayload) (res *web.SentimentItem, err error) {
	res = &web.SentimentItem{}
	log.Printf(ctx, "web.sentiments")
	return
}

// List root domains.
func (s *websrvc) Domains(ctx context.Context, p *web.DomainsPayload) (res []*web.DomainEntry, err error) {
	log.Printf(ctx, "web.domains")
	return
}

// List top trends for a domain.
func (s *websrvc) TopTrendsByDomain(ctx context.Context, p *web.TopTrendsByDomainPayload) (res []*web.TrendMetric, err error) {
	log.Printf(ctx, "web.topTrendsByDomain")
	return
}

// List trend context for a domain.
func (s *websrvc) ContextByDomain(ctx context.Context, p *web.ContextByDomainPayload) (res []*web.TrendContext, err error) {
	log.Printf(ctx, "web.contextByDomain")
	return
}

// List trend lifecycle data for a domain.
func (s *websrvc) LifecycleByDomain(ctx context.Context, p *web.LifecycleByDomainPayload) (res []*web.Lifecycle, err error) {
	log.Printf(ctx, "web.lifecycleByDomain")
	return
}

// Compare two domains for a trend.
func (s *websrvc) DomainComparisonEndpoint(ctx context.Context, p *web.DomainComparisonPayload) (res []*web.DomainComparison, err error) {
	log.Printf(ctx, "web.domainComparison")
	return
}

func (s *websrvc) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	return ctx, nil
}
