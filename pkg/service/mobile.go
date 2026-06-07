package service

import (
	"context"

	web "github.com/deframer/news-deframer/gen/web"

	mobile "github.com/deframer/news-deframer/gen/mobile"
	"goa.design/goa/v3/security"
)

// WARNING: This file is a thin adapter over webimpl.go.
// Keep it in sync with pkg/service/webimpl.go for every endpoint and DTO change.
// If webimpl.go changes, this file must be updated too, or mobile will silently diverge.
type mobilesrvc struct {
	svc web.Service
}

// NewMobile returns the mobile service implementation.
func NewMobile(ctx context.Context) mobile.Service {
	return &mobilesrvc{svc: NewWebImplementation(ctx)}
}

func (s *mobilesrvc) Item(ctx context.Context, p *mobile.ItemPayload) (res *mobile.AnalyzedItem, err error) {
	item, err := s.svc.Item(ctx, &web.ItemPayload{URL: p.URL, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	return convertMobileAnalyzedItem(item), nil
}

func (s *mobilesrvc) Site(ctx context.Context, p *mobile.SitePayload) (res []*mobile.AnalyzedSiteItem, err error) {
	items, err := s.svc.Site(ctx, &web.SitePayload{Root: p.Root, MaxScore: p.MaxScore, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.AnalyzedSiteItem, 0, len(items))
	for i := range items {
		res = append(res, convertMobileAnalyzedSiteItem(items[i]))
	}
	return res, nil
}

func (s *mobilesrvc) Articles(ctx context.Context, p *mobile.ArticlesPayload) (res []*mobile.AnalyzedArticle, err error) {
	items, err := s.svc.Articles(ctx, &web.ArticlesPayload{Root: p.Root, Term: p.Term, Date: p.Date, Days: p.Days, Offset: p.Offset, Limit: p.Limit, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.AnalyzedArticle, 0, len(items))
	for i := range items {
		res = append(res, convertMobileAnalyzedArticle(items[i]))
	}
	return res, nil
}

func (s *mobilesrvc) Sentiments(ctx context.Context, p *mobile.SentimentsPayload) (res *mobile.SentimentItem, err error) {
	item, err := s.svc.Sentiments(ctx, &web.SentimentsPayload{Root: p.Root, Term: p.Term, Date: p.Date, Days: p.Days, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	return convertMobileSentimentItem(item), nil
}

func (s *mobilesrvc) Domains(ctx context.Context, p *mobile.DomainsPayload) (res []*mobile.DomainEntry, err error) {
	domains, err := s.svc.Domains(ctx, &web.DomainsPayload{User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.DomainEntry, 0, len(domains))
	for i := range domains {
		res = append(res, convertMobileDomainEntry(domains[i]))
	}
	return res, nil
}

func (s *mobilesrvc) TopTrendsByDomain(ctx context.Context, p *mobile.TopTrendsByDomainPayload) (res []*mobile.TrendMetric, err error) {
	trends, err := s.svc.TopTrendsByDomain(ctx, &web.TopTrendsByDomainPayload{Domain: p.Domain, Lang: p.Lang, Days: p.Days, Date: p.Date, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.TrendMetric, 0, len(trends))
	for i := range trends {
		res = append(res, convertMobileTrendMetric(trends[i]))
	}
	return res, nil
}

func (s *mobilesrvc) ContextByDomain(ctx context.Context, p *mobile.ContextByDomainPayload) (res []*mobile.TrendContext, err error) {
	contexts, err := s.svc.ContextByDomain(ctx, &web.ContextByDomainPayload{Term: p.Term, Domain: p.Domain, Lang: p.Lang, Days: p.Days, Date: p.Date, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.TrendContext, 0, len(contexts))
	for i := range contexts {
		res = append(res, convertMobileTrendContext(contexts[i]))
	}
	return res, nil
}

func (s *mobilesrvc) LifecycleByDomain(ctx context.Context, p *mobile.LifecycleByDomainPayload) (res []*mobile.Lifecycle, err error) {
	lifecycles, err := s.svc.LifecycleByDomain(ctx, &web.LifecycleByDomainPayload{Term: p.Term, Domain: p.Domain, Lang: p.Lang, Days: p.Days, Date: p.Date, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.Lifecycle, 0, len(lifecycles))
	for i := range lifecycles {
		res = append(res, convertMobileLifecycle(lifecycles[i]))
	}
	return res, nil
}

func (s *mobilesrvc) DomainComparisonEndpoint(ctx context.Context, p *mobile.DomainComparisonPayload) (res []*mobile.DomainComparison, err error) {
	comparisons, err := s.svc.DomainComparisonEndpoint(ctx, &web.DomainComparisonPayload{DomainA: p.DomainA, DomainB: p.DomainB, Lang: p.Lang, Days: p.Days, Date: p.Date, User: p.User, Pass: p.Pass})
	if err != nil {
		return nil, translateMobileError(err)
	}
	res = make([]*mobile.DomainComparison, 0, len(comparisons))
	for i := range comparisons {
		res = append(res, convertMobileDomainComparison(comparisons[i]))
	}
	return res, nil
}

func (s *mobilesrvc) BasicAuth(ctx context.Context, user, pass string, scheme *security.BasicScheme) (context.Context, error) {
	if auther, ok := s.svc.(interface {
		BasicAuth(context.Context, string, string, *security.BasicScheme) (context.Context, error)
	}); ok {
		return auther.BasicAuth(ctx, user, pass, scheme)
	}
	return ctx, nil
}

func translateMobileError(err error) error {
	if err == nil {
		return nil
	}
	if nf, ok := err.(web.NotFound); ok {
		return mobile.NotFound(nf)
	}
	return err
}

func convertMobileAnalyzedItem(item *web.AnalyzedItem) *mobile.AnalyzedItem {
	if item == nil {
		return nil
	}
	return &mobile.AnalyzedItem{
		Hash:                        item.Hash,
		Tags:                        append([]string{}, item.Tags...),
		URL:                         item.URL,
		TitleOriginal:               item.TitleOriginal,
		DescriptionOriginal:         item.DescriptionOriginal,
		TitleCorrected:              item.TitleCorrected,
		TitleCorrectionReason:       item.TitleCorrectionReason,
		DescriptionCorrected:        item.DescriptionCorrected,
		DescriptionCorrectionReason: item.DescriptionCorrectionReason,
		Framing:                     item.Framing,
		FramingReason:               item.FramingReason,
		Clickbait:                   item.Clickbait,
		ClickbaitReason:             item.ClickbaitReason,
		Persuasive:                  item.Persuasive,
		PersuasiveReason:            item.PersuasiveReason,
		HyperStimulus:               item.HyperStimulus,
		HyperStimulusReason:         item.HyperStimulusReason,
		Speculative:                 item.Speculative,
		SpeculativeReason:           item.SpeculativeReason,
		Overall:                     item.Overall,
		OverallReason:               item.OverallReason,
		Category:                    item.Category,
		Sentiments:                  convertMobileSentimentScores(item.Sentiments),
		SentimentsDeframed:          convertMobileSentimentScores(item.SentimentsDeframed),
		Media:                       convertMobileMediaContent(item.Media),
		Rating:                      item.Rating,
		Authors:                     append([]string{}, item.Authors...),
		PubDate:                     item.PubDate,
	}
}

func convertMobileAnalyzedSiteItem(item *web.AnalyzedSiteItem) *mobile.AnalyzedSiteItem {
	if item == nil {
		return nil
	}
	return &mobile.AnalyzedSiteItem{
		URL:                  item.URL,
		TitleOriginal:        item.TitleOriginal,
		DescriptionOriginal:  item.DescriptionOriginal,
		TitleCorrected:       item.TitleCorrected,
		DescriptionCorrected: item.DescriptionCorrected,
		OverallReason:        item.OverallReason,
		Media:                convertMobileMediaContent(item.Media),
		Rating:               item.Rating,
		Authors:              append([]string{}, item.Authors...),
		PubDate:              item.PubDate,
	}
}

func convertMobileAnalyzedArticle(item *web.AnalyzedArticle) *mobile.AnalyzedArticle {
	if item == nil {
		return nil
	}
	return &mobile.AnalyzedArticle{
		URL:     item.URL,
		Title:   item.Title,
		Rating:  item.Rating,
		Authors: append([]string{}, item.Authors...),
		PubDate: item.PubDate,
	}
}

func convertMobileSentimentItem(item *web.SentimentItem) *mobile.SentimentItem {
	if item == nil {
		return nil
	}
	return &mobile.SentimentItem{
		Sentiments:         convertMobileSentimentScores(item.Sentiments),
		SentimentsDeframed: convertMobileSentimentScores(item.SentimentsDeframed),
	}
}

func convertMobileSentimentScores(scores *web.SentimentScores) *mobile.SentimentScores {
	if scores == nil {
		return nil
	}
	return &mobile.SentimentScores{
		Valence:   scores.Valence,
		Arousal:   scores.Arousal,
		Dominance: scores.Dominance,
		Joy:       scores.Joy,
		Anger:     scores.Anger,
		Sadness:   scores.Sadness,
		Fear:      scores.Fear,
		Disgust:   scores.Disgust,
	}
}

func convertMobileMediaContent(content *web.MediaContent) *mobile.MediaContent {
	if content == nil {
		return nil
	}
	return &mobile.MediaContent{
		URL:         content.URL,
		Type:        content.Type,
		Medium:      content.Medium,
		Height:      content.Height,
		Width:       content.Width,
		Title:       content.Title,
		Description: content.Description,
		Thumbnail:   convertMobileMediaThumbnail(content.Thumbnail),
		Credit:      content.Credit,
	}
}

func convertMobileMediaThumbnail(thumb *web.MediaThumbnail) *mobile.MediaThumbnail {
	if thumb == nil {
		return nil
	}
	return &mobile.MediaThumbnail{
		URL:    thumb.URL,
		Height: thumb.Height,
		Width:  thumb.Width,
	}
}

func convertMobileDomainEntry(entry *web.DomainEntry) *mobile.DomainEntry {
	if entry == nil {
		return nil
	}
	return &mobile.DomainEntry{
		Domain:    entry.Domain,
		Language:  entry.Language,
		Tags:      append([]string{}, entry.Tags...),
		PortalURL: entry.PortalURL,
	}
}

func convertMobileTrendMetric(metric *web.TrendMetric) *mobile.TrendMetric {
	if metric == nil {
		return nil
	}
	return &mobile.TrendMetric{
		TrendTopic:   metric.TrendTopic,
		Frequency:    metric.Frequency,
		Utility:      metric.Utility,
		OutlierRatio: metric.OutlierRatio,
		TimeSlice:    metric.TimeSlice,
	}
}

func convertMobileTrendContext(context *web.TrendContext) *mobile.TrendContext {
	if context == nil {
		return nil
	}
	return &mobile.TrendContext{
		Context:   context.Context,
		Frequency: context.Frequency,
	}
}

func convertMobileLifecycle(lifecycle *web.Lifecycle) *mobile.Lifecycle {
	if lifecycle == nil {
		return nil
	}
	return &mobile.Lifecycle{
		TimeSlice: lifecycle.TimeSlice,
		Frequency: lifecycle.Frequency,
		Velocity:  lifecycle.Velocity,
	}
}

func convertMobileDomainComparison(comparison *web.DomainComparison) *mobile.DomainComparison {
	if comparison == nil {
		return nil
	}
	return &mobile.DomainComparison{
		Classification: comparison.Classification,
		RankGroup:      comparison.RankGroup,
		TrendTopic:     comparison.TrendTopic,
		ScoreA:         comparison.ScoreA,
		ScoreB:         comparison.ScoreB,
	}
}
