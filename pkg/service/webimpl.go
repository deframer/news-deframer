package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	web "github.com/deframer/news-deframer/gen/web"
	"github.com/deframer/news-deframer/pkg/config"
	"github.com/deframer/news-deframer/pkg/database"
	"github.com/deframer/news-deframer/pkg/facade"
	"goa.design/clue/log"
	"goa.design/goa/v3/security"
)

// WebImpl implements the web service example behavior.
type WebImpl struct {
	facade facade.Facade
	cfg    *config.Config
}

// NewWebImplementation returns the shared web service implementation.
func NewWebImplementation(ctx context.Context) web.Service {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	repo, err := database.NewRepository(ctx, cfg)
	if err != nil {
		panic(err)
	}

	return &WebImpl{
		facade: facade.New(ctx, cfg, repo),
		cfg:    cfg,
	}
}

func (w *WebImpl) Item(ctx context.Context, p *web.ItemPayload) (res *web.AnalyzedItem, err error) {
	log.Printf(ctx, "handleItem url=%s", p.URL)
	reqURL := strings.TrimSuffix(p.URL, "/")
	if reqURL == "" {
		return nil, fmt.Errorf("missing url")
	}

	u, err := url.ParseRequestURI(reqURL)
	if err != nil {
		log.Errorf(ctx, err, "invalid url")
		return nil, fmt.Errorf("invalid url")
	}

	item, err := w.facade.GetFirstItemForUrl(ctx, u)
	if err != nil {
		log.Errorf(ctx, err, "failed to get item")
		return nil, web.NotFound("not found")
	}
	if item == nil || item.ThinkResult == nil {
		return nil, web.NotFound("not found")
	}

	return convertAnalyzedItem(item), nil
}

func (w *WebImpl) Site(ctx context.Context, p *web.SitePayload) (res []*web.AnalyzedItem, err error) {
	log.Printf(ctx, "handleSite root=%s", p.Root)
	rootDomain := strings.TrimSuffix(p.Root, "/")
	if rootDomain == "" {
		return nil, fmt.Errorf("missing root")
	}

	items, err := w.facade.GetItemsForRootDomain(ctx, rootDomain, p.MaxScore)
	if err != nil || len(items) == 0 {
		if err != nil {
			log.Errorf(ctx, err, "GetItemsForRootDomain failed")
		}
		return nil, web.NotFound("not found")
	}

	res = make([]*web.AnalyzedItem, 0, len(items))
	for i := range items {
		res = append(res, convertAnalyzedItem(&items[i]))
	}
	return res, nil
}

func (w *WebImpl) Articles(ctx context.Context, p *web.ArticlesPayload) (res []*web.AnalyzedArticle, err error) {
	log.Printf(ctx, "handleArticles root=%s term=%s", p.Root, p.Term)
	rootDomain := strings.TrimSuffix(p.Root, "/")
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if rootDomain == "" || p.Term == "" {
		return nil, fmt.Errorf("missing root or term")
	}

	items, err := w.facade.GetArticlesByTrend(ctx, p.Term, rootDomain, date, p.Days, p.Offset, p.Limit)
	if err != nil || len(items) == 0 {
		if err != nil {
			log.Errorf(ctx, err, "GetArticlesByTrend failed")
		}
		return nil, web.NotFound("not found")
	}

	res = make([]*web.AnalyzedArticle, 0, len(items))
	for i := range items {
		res = append(res, convertAnalyzedArticle(&items[i]))
	}
	return res, nil
}

func (w *WebImpl) Sentiments(ctx context.Context, p *web.SentimentsPayload) (res *web.SentimentItem, err error) {
	log.Printf(ctx, "handleSentiments root=%s term=%s", p.Root, p.Term)
	rootDomain := strings.TrimSuffix(p.Root, "/")
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if rootDomain == "" || p.Term == "" {
		return nil, fmt.Errorf("missing root or term")
	}

	item, err := w.facade.GetSentimentsByTrend(ctx, p.Term, rootDomain, date, p.Days)
	if err != nil || item == nil {
		if err != nil {
			log.Errorf(ctx, err, "GetSentimentsByTrend failed")
		}
		return nil, web.NotFound("not found")
	}

	return convertSentimentItem(item), nil
}

func (w *WebImpl) Domains(ctx context.Context, p *web.DomainsPayload) (res []*web.DomainEntry, err error) {
	log.Printf(ctx, "handleDomains")
	domains, err := w.facade.GetRootDomains(ctx)
	if err != nil {
		log.Errorf(ctx, err, "failed to get root domains")
		return nil, web.NotFound("not found")
	}

	res = make([]*web.DomainEntry, 0, len(domains))
	for i := range domains {
		res = append(res, &web.DomainEntry{
			Domain:    domains[i].Domain,
			Language:  domains[i].Language,
			PortalURL: domains[i].PortalUrl,
		})
	}
	return res, nil
}

func (w *WebImpl) TopTrendsByDomain(ctx context.Context, p *web.TopTrendsByDomainPayload) (res []*web.TrendMetric, err error) {
	log.Printf(ctx, "handleTopTrendsByDomain domain=%s", p.Domain)
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if p.Domain == "" || p.Lang == "" {
		return nil, fmt.Errorf("missing domain or lang")
	}

	trends, err := w.facade.GetTopTrendByDomain(ctx, p.Domain, p.Lang, date, p.Days)
	if err != nil {
		log.Errorf(ctx, err, "failed to get top trends")
		return nil, web.NotFound("not found")
	}

	res = make([]*web.TrendMetric, 0, len(trends))
	for i := range trends {
		res = append(res, &web.TrendMetric{
			TrendTopic:   trends[i].TrendTopic,
			Frequency:    trends[i].Frequency,
			Utility:      trends[i].Utility,
			OutlierRatio: trends[i].OutlierRatio,
			TimeSlice:    trends[i].TimeSlice.Format(time.RFC3339),
		})
	}
	return res, nil
}

func (w *WebImpl) ContextByDomain(ctx context.Context, p *web.ContextByDomainPayload) (res []*web.TrendContext, err error) {
	log.Printf(ctx, "handleContextByDomain domain=%s term=%s", p.Domain, p.Term)
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if p.Term == "" || p.Domain == "" || p.Lang == "" {
		return nil, fmt.Errorf("missing term, domain or lang")
	}

	contexts, err := w.facade.GetContextByDomain(ctx, p.Term, p.Domain, p.Lang, date, p.Days)
	if err != nil {
		log.Errorf(ctx, err, "failed to get context by domain")
		return nil, web.NotFound("not found")
	}

	res = make([]*web.TrendContext, 0, len(contexts))
	for i := range contexts {
		res = append(res, &web.TrendContext{
			Context:   contexts[i].Context,
			Frequency: contexts[i].Frequency,
		})
	}
	return res, nil
}

func (w *WebImpl) LifecycleByDomain(ctx context.Context, p *web.LifecycleByDomainPayload) (res []*web.Lifecycle, err error) {
	log.Printf(ctx, "handleLifecycleByDomain domain=%s term=%s", p.Domain, p.Term)
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if p.Term == "" || p.Domain == "" || p.Lang == "" {
		return nil, fmt.Errorf("missing term, domain or lang")
	}

	lifecycles, err := w.facade.GetLifecycleByDomain(ctx, p.Term, p.Domain, p.Lang, date, p.Days)
	if err != nil {
		log.Errorf(ctx, err, "failed to get lifecycle by domain")
		return nil, web.NotFound("not found")
	}

	res = make([]*web.Lifecycle, 0, len(lifecycles))
	for i := range lifecycles {
		res = append(res, &web.Lifecycle{
			TimeSlice: lifecycles[i].TimeSlice.Format(time.RFC3339),
			Frequency: lifecycles[i].Frequency,
			Velocity:  lifecycles[i].Velocity,
		})
	}
	return res, nil
}

func (w *WebImpl) DomainComparisonEndpoint(ctx context.Context, p *web.DomainComparisonPayload) (res []*web.DomainComparison, err error) {
	log.Printf(ctx, "handleDomainComparison domain_a=%s domain_b=%s", p.DomainA, p.DomainB)
	date, err := parseOptionalDateParam(p.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if p.DomainA == "" || p.DomainB == "" || p.Lang == "" {
		return nil, fmt.Errorf("missing domain_a, domain_b or lang")
	}

	comparisons, err := w.facade.GetDomainComparison(ctx, p.DomainA, p.DomainB, p.Lang, date, p.Days)
	if err != nil {
		log.Errorf(ctx, err, "failed to get domain comparison")
		return nil, web.NotFound("not found")
	}

	res = make([]*web.DomainComparison, 0, len(comparisons))
	for i := range comparisons {
		res = append(res, &web.DomainComparison{
			Classification: comparisons[i].Classification,
			RankGroup:      comparisons[i].RankGroup,
			TrendTopic:     comparisons[i].TrendTopic,
			ScoreA:         comparisons[i].ScoreA,
			ScoreB:         comparisons[i].ScoreB,
		})
	}
	return res, nil
}

func (w *WebImpl) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	return ctx, nil
}

func convertAnalyzedItem(item *database.AnalyzedItem) *web.AnalyzedItem {
	if item == nil {
		return nil
	}
	return &web.AnalyzedItem{
		Hash:               item.Hash,
		URL:                item.URL,
		Sentiments:         convertSentimentScores(item.Sentiments),
		SentimentsDeframed: convertSentimentScores(item.SentimentsDeframed),
		Media:              convertMediaContent(item.MediaContent),
		Rating:             item.ThinkRating,
		Authors:            append([]string(nil), item.Authors...),
		PubDate:            item.PubDate.Format(time.RFC3339),
	}
}

func convertAnalyzedArticle(item *database.AnalyzedArticle) *web.AnalyzedArticle {
	if item == nil {
		return nil
	}
	return &web.AnalyzedArticle{
		URL:     item.URL,
		Title:   item.Title,
		Rating:  item.Rating,
		Authors: append([]string(nil), item.Authors...),
		PubDate: item.PubDate.Format(time.RFC3339),
	}
}

func convertSentimentItem(item *database.SentimentItem) *web.SentimentItem {
	if item == nil {
		return nil
	}
	return &web.SentimentItem{
		Sentiments:         convertSentimentScores(item.Sentiments),
		SentimentsDeframed: convertSentimentScores(item.SentimentsDeframed),
	}
}

func convertSentimentScores(scores *database.SentimentScores) *web.SentimentScores {
	if scores == nil {
		return nil
	}
	return &web.SentimentScores{
		Valence:   float64Ptr(scores.Valence),
		Arousal:   float64Ptr(scores.Arousal),
		Dominance: float64Ptr(scores.Dominance),
		Joy:       float64Ptr(scores.Joy),
		Anger:     float64Ptr(scores.Anger),
		Sadness:   float64Ptr(scores.Sadness),
		Fear:      float64Ptr(scores.Fear),
		Disgust:   float64Ptr(scores.Disgust),
	}
}

func convertMediaContent(content *database.MediaContent) *web.MediaContent {
	if content == nil {
		return nil
	}
	return &web.MediaContent{
		URL:         content.URL,
		Type:        stringPtr(content.Type),
		Medium:      stringPtr(content.Medium),
		Height:      intPtr(content.Height),
		Width:       intPtr(content.Width),
		Title:       stringPtr(content.Title),
		Description: stringPtr(content.Description),
		Thumbnail:   convertMediaThumbnail(content.Thumbnail),
		Credit:      stringPtr(content.Credit),
	}
}

func convertMediaThumbnail(thumb *database.MediaThumbnail) *web.MediaThumbnail {
	if thumb == nil {
		return nil
	}
	return &web.MediaThumbnail{
		URL:    thumb.URL,
		Height: intPtr(thumb.Height),
		Width:  intPtr(thumb.Width),
	}
}

func parseOptionalDateParam(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", *raw)
	if err != nil {
		return nil, err
	}
	normalized := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
	return &normalized, nil
}

func float64Ptr(v float64) *float64 {
	if v == 0 {
		return nil
	}
	return &v
}

func intPtr(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
