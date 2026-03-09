export interface DomainEntry {
  domain: string;
  language: string;
  portal_url?: string;
}

export interface TrendMetric {
  trend_topic: string;
  frequency: number;
  utility: number;
  outlier_ratio: number;
  time_slice: string;
}

export interface TrendContext {
  context: string;
  frequency: number;
}

export interface Lifecycle {
  time_slice: string;
  frequency: number;
  velocity: number;
}

export interface DomainComparison {
  classification: 'BLINDSPOT_A' | 'BLINDSPOT_B' | 'INTERSECT';
  rank_group: number;
  trend_topic: string;
  score_a: number;
  score_b: number;
}

export interface ThinkResult {
  title_original?: string;
  description_original?: string;
  title_corrected?: string;
  title_correction_reason?: string;
  description_corrected?: string;
  description_correction_reason?: string;
  framing?: number;
  framing_reason?: string;
  clickbait?: number;
  clickbait_reason?: string;
  persuasive?: number;
  persuasive_reason?: string;
  hyper_stimulus?: number;
  hyper_stimulus_reason?: string;
  speculative?: number;
  speculative_reason?: string;
  overall?: number;
  overall_reason?: string;
}

export interface MediaThumbnail {
  url?: string;
  height?: number;
  width?: number;
}

export interface MediaContent {
  url?: string;
  type?: string;
  medium?: string;
  height?: number;
  width?: number;
  title?: string;
  description?: string;
  thumbnail?: MediaThumbnail;
  credit?: string;
}

export interface AnalyzedItem extends ThinkResult {
  hash: string;
  url: string;
  authors?: string[];
  media?: MediaContent;
  rating: number;
  pubDate?: string;
}

export interface AnalyzedArticle {
  url: string;
  title?: string;
  rating?: number;
  authors?: string[];
  pub_date: string;
}

export interface NewsDeframerApi {
  getDomains(): Promise<DomainEntry[]>;
  getItem(url: string): Promise<AnalyzedItem | null>;
  getSite(root: string, maxScore?: number): Promise<AnalyzedItem[]>;
  getTopTrendByDomain(domain: string, language: string, daysInPast: number): Promise<TrendMetric[]>;
  getContextByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<TrendContext[]>;
  getLifecycleByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<Lifecycle[]>;
  getDomainComparison(domainA: string, domainB: string, language: string, daysInPast: number): Promise<DomainComparison[]>;
  getArticlesByTrend(root: string, term: string, date?: string, days?: number, offset?: number, limit?: number): Promise<AnalyzedArticle[]>;
}
