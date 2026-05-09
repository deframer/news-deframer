import { encode as encodeBase64 } from 'base-64';

import {
  getArticlesByTrend,
  getContextByDomain,
  getDomainComparison,
  getDomains,
  getItem,
  getLifecycleByDomain,
  getSentimentsByTrend,
  getSite,
  getTopTrendByDomain,
  type AnalyzedArticle,
  type AnalyzedItem,
  type AnalyzedSiteItem,
  type DomainComparison,
  type DomainEntry,
  type Lifecycle,
  type SentimentItem,
  type SentimentScores,
  type ThinkResult,
  type TrendContext,
  type TrendMetric,
} from './generated/newsDeframerClient.gen';
import { logger } from './logger';
import { Settings } from './settingsService';

export type { AnalyzedArticle, AnalyzedItem, AnalyzedSiteItem, DomainComparison, DomainEntry, Lifecycle, SentimentItem, SentimentScores, ThinkResult, TrendContext, TrendMetric };

type GeneratedResponse<T> = {
  data: T;
  status: number;
};

export class NewsDeframerClient {
  constructor(private config: Settings) {}

  private getBackendBaseUrl(): string {
    return (this.config.backendUrl || '').trim().replace(/\/$/, '');
  }

  private createFetch(): typeof globalThis.fetch {
    return (input, init) => {
      const path = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
      const url = /^https?:\/\//i.test(path) ? path : `${this.getBackendBaseUrl()}${path}`;
      const headers = new Headers(init?.headers);

      if (!headers.has('Accept')) {
        headers.set('Accept', 'application/json');
      }

      if (this.config.username && this.config.password) {
        headers.set('Authorization', `Basic ${encodeBase64(`${this.config.username}:${this.config.password}`)}`);
      }

      return fetch(url, { ...init, headers });
    };
  }

  private async unwrap<T>(operation: Promise<GeneratedResponse<T>>, fallback: T | null): Promise<T | null> {
    const result = await operation;

    if (result.status === 404) {
      return fallback;
    }

    if (result.status < 200 || result.status >= 300) {
      logger.error('API non-ok response', { status: result.status });
      throw new Error(`API Error: ${result.status}`);
    }

    return result.data;
  }

  async getDomains(): Promise<DomainEntry[]> {
    const response = await this.unwrap<DomainEntry[]>(getDomains(undefined, this.createFetch()), null);
    if (response === null) {
      throw new Error('Domains endpoint returned 404');
    }
    return response;
  }

  async getItem(url: string): Promise<AnalyzedItem | null> {
    return this.unwrap<AnalyzedItem>(getItem({ url }, undefined, this.createFetch()), null);
  }

  async getSite(root: string, maxScore?: number): Promise<AnalyzedItem[]> {
    const response = await this.unwrap<AnalyzedSiteItem[]>(getSite({ root, max_score: maxScore }, undefined, this.createFetch()), []);
    return response as AnalyzedItem[];
  }

  async getTopTrendByDomain(domain: string, language: string, daysInPast: number): Promise<TrendMetric[]> {
    logger.info('Fetching top trends by domain', { domain, language, daysInPast });
    return (await this.unwrap<TrendMetric[]>(getTopTrendByDomain({ domain, lang: language, days: daysInPast }, undefined, this.createFetch()), [])) ?? [];
  }

  async getContextByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<TrendContext[]> {
    return (await this.unwrap<TrendContext[]>(getContextByDomain({ term, domain, lang: language, days: daysInPast }, undefined, this.createFetch()), [])) ?? [];
  }

  async getLifecycleByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<Lifecycle[]> {
    return (await this.unwrap<Lifecycle[]>(getLifecycleByDomain({ term, domain, lang: language, days: daysInPast }, undefined, this.createFetch()), [])) ?? [];
  }

  async getDomainComparison(domainA: string, domainB: string, language: string, daysInPast: number): Promise<DomainComparison[]> {
    return (await this.unwrap<DomainComparison[]>(getDomainComparison({ domain_a: domainA, domain_b: domainB, lang: language, days: daysInPast }, undefined, this.createFetch()), [])) ?? [];
  }

  async getArticlesByTrend(root: string, term: string, date?: string, days?: number, offset?: number, limit?: number): Promise<AnalyzedArticle[]> {
    return (await this.unwrap<AnalyzedArticle[]>(getArticlesByTrend({ root, term, date, days, offset, limit }, undefined, this.createFetch()), [])) ?? [];
  }

  async getSentimentsByTrend(root: string, term: string, date?: string, days?: number): Promise<SentimentItem | null> {
    return this.unwrap<SentimentItem>(getSentimentsByTrend({ root, term, date, days }, undefined, this.createFetch()), null);
  }
}
