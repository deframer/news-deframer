import { encode as encodeBase64 } from 'base-64';

import { logger } from './logger';
import { Settings } from './settingsService';

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

export interface SentimentItem {
  valence?: number | null;
  arousal?: number | null;
  dominance?: number | null;
  joy?: number | null;
  anger?: number | null;
  sadness?: number | null;
  fear?: number | null;
  disgust?: number | null;
}

export class NewsDeframerClient {
  constructor(private config: Settings) {}

  // Mobile must use the top-level /mobile/api namespace for every backend endpoint.
  private readonly apiBase = '/mobile/api';

  private getBackendBaseUrl(): string {
    return this.config.backendUrl
      .trim()
      .replace(/\/$/, '');
  }

  private async request<T>(endpoint: string, params: Record<string, string>): Promise<T | null> {
    const headers: Record<string, string> = {};

    if (this.config.username && this.config.password) {
      headers.Authorization = `Basic ${encodeBase64(`${this.config.username}:${this.config.password}`)}`;
    }

    const baseUrl = this.getBackendBaseUrl().replace(/\/+$/, '');
    const endpointPath = endpoint.replace(/^\/+/, '');
    const searchParams = new URLSearchParams(params);
    const queryString = searchParams.toString();
    const requestUrl = queryString ? `${baseUrl}/${endpointPath}?${queryString}` : `${baseUrl}/${endpointPath}`;

    logger.info('API request', { endpoint, params, url: requestUrl });

    let response: Response;
    try {
      response = await fetch(requestUrl, { headers });
    } catch (error) {
      logger.error('API network error', { endpoint, url: requestUrl, error: String(error) });
      const rawMessage = error instanceof Error ? error.message : String(error);
      throw new Error(`GET ${requestUrl} failed: ${rawMessage}`);
    }
    if (response.status === 404) {
      logger.warn('API returned 404', { endpoint, url: requestUrl });
      return null;
    }
    if (!response.ok) {
      logger.error('API non-ok response', { endpoint, url: requestUrl, status: response.status });
      throw new Error(`API Error: ${response.status}`);
    }

    return (await response.json()) as T;
  }

  async getDomains(): Promise<DomainEntry[]> {
    const response = await this.request<DomainEntry[]>(`${this.apiBase}/domains`, {});
    if (response === null) {
      throw new Error('Domains endpoint returned 404');
    }
    return response;
  }

  async getItem(url: string): Promise<AnalyzedItem | null> {
    return this.request<AnalyzedItem>(`${this.apiBase}/item`, { url });
  }

  async getSite(root: string, maxScore?: number): Promise<AnalyzedItem[]> {
    const params: Record<string, string> = { root };
    if (maxScore !== undefined) {
      params.max_score = String(maxScore);
    }
    return (await this.request<AnalyzedItem[]>(`${this.apiBase}/site`, params)) ?? [];
  }

  async getTopTrendByDomain(domain: string, language: string, daysInPast: number): Promise<TrendMetric[]> {
    logger.info('Fetching top trends by domain', { domain, language, daysInPast });
    return (await this.request<TrendMetric[]>(`${this.apiBase}/trends/topbydomain`, {
      domain,
      lang: language,
      days: String(daysInPast),
    })) ?? [];
  }

  async getContextByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<TrendContext[]> {
    return (await this.request<TrendContext[]>(`${this.apiBase}/trends/contextbydomain`, {
      term,
      domain,
      lang: language,
      days: String(daysInPast),
    })) ?? [];
  }

  async getLifecycleByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<Lifecycle[]> {
    return (await this.request<Lifecycle[]>(`${this.apiBase}/trends/lifecyclebydomain`, {
      term,
      domain,
      lang: language,
      days: String(daysInPast),
    })) ?? [];
  }

  async getDomainComparison(domainA: string, domainB: string, language: string, daysInPast: number): Promise<DomainComparison[]> {
    return (await this.request<DomainComparison[]>(`${this.apiBase}/trends/comparedomains`, {
      domain_a: domainA,
      domain_b: domainB,
      lang: language,
      days: String(daysInPast),
    })) ?? [];
  }

  async getArticlesByTrend(root: string, term: string, date?: string, days?: number, offset?: number, limit?: number): Promise<AnalyzedArticle[]> {
    const params: Record<string, string> = { root, term };
    if (date) params.date = date;
    if (days !== undefined) params.days = String(days);
    if (offset !== undefined) params.offset = String(offset);
    if (limit !== undefined) params.limit = String(limit);
    return (await this.request<AnalyzedArticle[]>(`${this.apiBase}/articles`, params)) ?? [];
  }

  async getSentimentsByTrend(root: string, term: string, date?: string, days?: number): Promise<SentimentItem | null> {
    const params: Record<string, string> = { root, term };
    if (date) params.date = date;
    if (days !== undefined) params.days = String(days);
    return this.request<SentimentItem>(`${this.apiBase}/sentiments`, params);
  }
}
