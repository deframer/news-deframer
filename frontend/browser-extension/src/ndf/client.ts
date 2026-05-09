import { getCachedDomains, invalidateDomainCache, setCachedDomains } from '../shared/domain-cache';
import { CONNECTION_TIMEOUT_MS, removeStallDomainsFromSelection, Settings } from '../shared/settings';
import {
  type AnalyzedArticle,
  type AnalyzedItem,
  type AnalyzedSiteItem,
  type DomainComparison,
  type DomainEntry,
  getArticlesByTrend,
  getContextByDomain,
  getDomainComparison,
  getDomains as getGeneratedDomains,
  getItem,
  getLifecycleByDomain,
  getSentimentsByTrend,
  getSite,
  getTopTrendByDomain,
  type Lifecycle,
  type SentimentItem,
  type SentimentScores,
  type ThinkResult,
  type TrendContext,
  type TrendMetric,
} from './generated/newsDeframerClient.gen';

export type { AnalyzedArticle, AnalyzedItem, DomainComparison, DomainEntry, Lifecycle, SentimentItem, SentimentScores, ThinkResult, TrendContext, TrendMetric };

type GeneratedResponse<T> = {
  data: T;
  status: number;
};

type ProxyResponse = {
  ok: boolean;
  status?: number;
  data?: unknown;
  error?: string;
};

export class NewsDeframerClient {
  constructor(private config: Settings) {}

  private getBackendBaseUrl(): string {
    return this.config.backendUrl.replace(/\/$/, '');
  }

  private createProxyFetch(): typeof globalThis.fetch {
    return (input, init) => new Promise<Response>((resolve, reject) => {
      if (typeof chrome === 'undefined' || !chrome.runtime?.sendMessage) {
        reject(new Error('Extension context missing: Cannot proxy request'));
        return;
      }

      const path = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
      const url = /^https?:\/\//i.test(path) ? path : `${this.getBackendBaseUrl()}${path}`;
      const headers = new Headers(init?.headers);

      if (!headers.has('Accept')) {
        headers.set('Accept', 'application/json');
      }

      if (this.config.username && this.config.password) {
        headers.set('Authorization', `Basic ${btoa(`${this.config.username}:${this.config.password}`)}`);
      }

      const timeoutId = setTimeout(() => {
        reject(new Error(`Request timed out after ${CONNECTION_TIMEOUT_MS}ms`));
      }, CONNECTION_TIMEOUT_MS + 1000);

      chrome.runtime.sendMessage(
        { type: 'PROXY_REQ', url, headers: Object.fromEntries(headers.entries()), timeout: CONNECTION_TIMEOUT_MS },
        (proxyResponse?: ProxyResponse) => {
          clearTimeout(timeoutId);

          if (chrome.runtime.lastError) {
            reject(new Error(chrome.runtime.lastError.message));
            return;
          }

          if (!proxyResponse) {
            reject(new Error(`GET ${url} failed: Unknown error`));
            return;
          }

          const status = proxyResponse.status ?? (proxyResponse.ok ? 200 : 500);
          const body = proxyResponse.ok ? proxyResponse.data : proxyResponse.data ?? proxyResponse.error ?? null;

          resolve(new Response(JSON.stringify(body), {
            status,
            headers: { 'Content-Type': 'application/json' },
          }));
        },
      );
    });
  }

  private async unwrap<T>(operation: Promise<GeneratedResponse<T>>, fallback: T | null): Promise<T | null> {
    const result = await operation;

    if (result.status === 404) {
      return fallback;
    }

    if (result.status < 200 || result.status >= 300) {
      throw new Error(`API Error: ${result.status}`);
    }

    return result.data;
  }

  async getDomains(bypassCache = false): Promise<DomainEntry[]> {
    if (!bypassCache) {
      const cached = await getCachedDomains();
      if (cached) {
        return cached;
      }
    }

    try {
      const domains = (await this.unwrap<DomainEntry[]>(getGeneratedDomains(undefined, this.createProxyFetch()), [])) ?? [];

      if (domains.length > 0) {
        await setCachedDomains(domains);
        await removeStallDomainsFromSelection(domains.map((domain) => domain.domain));
      } else {
        await invalidateDomainCache();
      }

      return domains;
    } catch (error) {
      await invalidateDomainCache();
      throw error;
    }
  }

  async getItem(url: string): Promise<AnalyzedItem | null> {
    return this.unwrap<AnalyzedItem>(getItem({ url }, undefined, this.createProxyFetch()), null);
  }

  async getSite(root: string, maxScore?: number): Promise<AnalyzedItem[]> {
    const result = await this.unwrap<AnalyzedSiteItem[]>(getSite({ root, max_score: maxScore }, undefined, this.createProxyFetch()), []);
    return (result ?? []) as AnalyzedItem[];
  }

  async getTopTrendByDomain(domain: string, language: string, daysInPast: number): Promise<TrendMetric[]> {
    return (await this.unwrap<TrendMetric[]>(getTopTrendByDomain({ domain, lang: language, days: daysInPast }, undefined, this.createProxyFetch()), [])) ?? [];
  }

  async getContextByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<TrendContext[]> {
    return (await this.unwrap<TrendContext[]>(getContextByDomain({ term, domain, lang: language, days: daysInPast }, undefined, this.createProxyFetch()), [])) ?? [];
  }

  async getLifecycleByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<Lifecycle[]> {
    return (await this.unwrap<Lifecycle[]>(getLifecycleByDomain({ term, domain, lang: language, days: daysInPast }, undefined, this.createProxyFetch()), [])) ?? [];
  }

  async getDomainComparison(domainA: string, domainB: string, language: string, daysInPast: number): Promise<DomainComparison[]> {
    return (await this.unwrap<DomainComparison[]>(getDomainComparison({ domain_a: domainA, domain_b: domainB, lang: language, days: daysInPast }, undefined, this.createProxyFetch()), [])) ?? [];
  }

  async getArticlesByTrend(root: string, term: string, date?: string, days?: number, offset?: number, limit?: number): Promise<AnalyzedArticle[]> {
    return (await this.unwrap<AnalyzedArticle[]>(getArticlesByTrend({ root, term, date, days, offset, limit }, undefined, this.createProxyFetch()), [])) ?? [];
  }

  async getSentimentsByTrend(root: string, term: string, date?: string, days?: number): Promise<SentimentItem | null> {
    return this.unwrap<SentimentItem>(getSentimentsByTrend({ root, term, date, days }, undefined, this.createProxyFetch()), null);
  }
}
