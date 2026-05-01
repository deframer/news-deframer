package servicenew

import (
	"context"

	mobile "github.com/deframer/news-deframer/gen/mobile"
	"goa.design/clue/log"
)

// mobile service example implementation.
// The example methods log the requests and return zero values.
type mobilesrvc struct{}

// NewMobile returns the mobile service implementation.
func NewMobile() mobile.Service {
	return &mobilesrvc{}
}

// Fetch a single analyzed item by URL.
func (s *mobilesrvc) Item(ctx context.Context, p *mobile.ItemPayload) (res *mobile.AnalyzedItem, err error) {
	res = &mobile.AnalyzedItem{}
	log.Printf(ctx, "mobile.item")
	return
}

// List analyzed items for a root domain.
func (s *mobilesrvc) Site(ctx context.Context, p *mobile.SitePayload) (res []*mobile.AnalyzedItem, err error) {
	log.Printf(ctx, "mobile.site")
	return
}

// List articles for a trend and domain.
func (s *mobilesrvc) Articles(ctx context.Context, p *mobile.ArticlesPayload) (res []*mobile.AnalyzedArticle, err error) {
	log.Printf(ctx, "mobile.articles")
	return
}

// Get sentiment scores for a trend.
func (s *mobilesrvc) Sentiments(ctx context.Context, p *mobile.SentimentsPayload) (res *mobile.SentimentItem, err error) {
	res = &mobile.SentimentItem{}
	log.Printf(ctx, "mobile.sentiments")
	return
}

// List root domains.
func (s *mobilesrvc) Domains(ctx context.Context) (res []*mobile.DomainEntry, err error) {
	log.Printf(ctx, "mobile.domains")
	return
}

// List top trends for a domain.
func (s *mobilesrvc) TopTrendsByDomain(ctx context.Context, p *mobile.TopTrendsByDomainPayload) (res []*mobile.TrendMetric, err error) {
	log.Printf(ctx, "mobile.topTrendsByDomain")
	return
}

// List trend context for a domain.
func (s *mobilesrvc) ContextByDomain(ctx context.Context, p *mobile.ContextByDomainPayload) (res []*mobile.TrendContext, err error) {
	log.Printf(ctx, "mobile.contextByDomain")
	return
}

// List trend lifecycle data for a domain.
func (s *mobilesrvc) LifecycleByDomain(ctx context.Context, p *mobile.LifecycleByDomainPayload) (res []*mobile.Lifecycle, err error) {
	log.Printf(ctx, "mobile.lifecycleByDomain")
	return
}

// Compare two domains for a trend.
func (s *mobilesrvc) DomainComparisonEndpoint(ctx context.Context, p *mobile.DomainComparisonPayload) (res []*mobile.DomainComparison, err error) {
	log.Printf(ctx, "mobile.domainComparison")
	return
}
