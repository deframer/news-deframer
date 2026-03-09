import {
  AnalyzedArticle,
  AnalyzedItem,
  DomainComparison,
  DomainEntry,
  Lifecycle,
  NewsDeframerApi,
  TrendContext,
  TrendMetric,
} from '@frontend-shared/ndf-api-interfaces';
import { Settings } from '@frontend-shared/settings';

import { getCachedDomains, invalidateDomainCache, setCachedDomains } from '../domain-cache';

export type {
  AnalyzedArticle,
  AnalyzedItem,
  DomainComparison,
  DomainEntry,
  Lifecycle,
  NewsDeframerApi,
  TrendContext,
  TrendMetric,
} from '@frontend-shared/ndf-api-interfaces';

export class NewsDeframerClient implements NewsDeframerApi {
  constructor(private config: Settings) {}

  private async proxyRequest<T>(endpoint: string, params: Record<string, string>): Promise<T | null> {
    const headers: Record<string, string> = {};
    if (this.config.username && this.config.password) {
      headers.Authorization = 'Basic ' + btoa(`${this.config.username}:${this.config.password}`);
    }

    const url = new URL(this.config.backendUrl.replace(/\/$/, '') + endpoint);
    Object.entries(params).forEach(([key, value]) => url.searchParams.append(key, value));

    if (typeof chrome !== 'undefined' && chrome.runtime && chrome.runtime.sendMessage) {
      return new Promise<T | null>((resolve, reject) => {
        chrome.runtime.sendMessage({ type: 'PROXY_REQ', url: url.toString(), headers }, (response) => {
          if (chrome.runtime.lastError) {
            reject(new Error(chrome.runtime.lastError.message));
          } else if (response && response.ok) {
            resolve(response.data as T);
          } else if (response && response.status === 404) {
            resolve(null);
          } else {
            reject(new Error(response?.error || `API Error: ${response?.status}`));
          }
        });
      });
    }

    throw new Error('Extension context missing: Cannot proxy request');
  }

  async getDomains(): Promise<DomainEntry[]> {
    const cached = await getCachedDomains();
    if (cached) {
      return cached;
    }

    try {
      const result = await this.proxyRequest<DomainEntry[]>('/api/domains', {});
      const domains = result ?? [];

      if (domains.length > 0) {
        await setCachedDomains(domains);
      }

      return domains;
    } catch (error) {
      await invalidateDomainCache();
      throw error;
    }
  }

  async getItem(url: string): Promise<AnalyzedItem | null> {
    return this.proxyRequest<AnalyzedItem>('/api/item', { url });
  }

  async getSite(root: string, maxScore?: number): Promise<AnalyzedItem[]> {
    const params: Record<string, string> = { root };
    if (maxScore !== undefined) {
      params.max_score = maxScore.toString();
    }
    const result = await this.proxyRequest<AnalyzedItem[]>('/api/site', params);
    return result ?? [];
  }

  async getTopTrendByDomain(domain: string, language: string, daysInPast: number): Promise<TrendMetric[]> {
    const params: Record<string, string> = {
      domain,
      lang: language,
      days: daysInPast.toString(),
    };
    const result = await this.proxyRequest<TrendMetric[]>('/api/trends/topbydomain', params);
    return result ?? [];
  }

  async getContextByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<TrendContext[]> {
    const params: Record<string, string> = {
      term,
      domain,
      lang: language,
      days: daysInPast.toString(),
    };
    const result = await this.proxyRequest<TrendContext[]>('/api/trends/contextbydomain', params);
    return result ?? [];
  }

  async getLifecycleByDomain(term: string, domain: string, language: string, daysInPast: number): Promise<Lifecycle[]> {
    const params: Record<string, string> = {
      term,
      domain,
      lang: language,
      days: daysInPast.toString(),
    };
    const result = await this.proxyRequest<Lifecycle[]>('/api/trends/lifecyclebydomain', params);
    return result ?? [];
  }

  async getDomainComparison(domainA: string, domainB: string, language: string, daysInPast: number): Promise<DomainComparison[]> {
    const params: Record<string, string> = {
      domain_a: domainA,
      domain_b: domainB,
      lang: language,
      days: daysInPast.toString(),
    };
    const result = await this.proxyRequest<DomainComparison[]>('/api/trends/comparedomains', params);
    return result ?? [];
  }

  async getArticlesByTrend(root: string, term: string, date?: string, days?: number, offset?: number, limit?: number): Promise<AnalyzedArticle[]> {
    const params: Record<string, string> = {
      root,
      term,
    };
    if (date) {
      params.date = date;
    }
    if (days !== undefined && days !== null) {
      params.days = days.toString();
    }
    if (offset !== undefined) {
      params.offset = offset.toString();
    }
    if (limit !== undefined) {
      params.limit = limit.toString();
    }
    const result = await this.proxyRequest<AnalyzedArticle[]>('/api/articles', params);
    return result ?? [];
  }
}
