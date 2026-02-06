import { getCachedDomains, invalidateDomainCache,setCachedDomains } from '../shared/domain-cache';
import { Settings } from '../shared/settings';

// --- Type Definitions based on Go backend models ---

export interface DomainEntry {
  domain: string;
  language: string;
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
  media?: MediaContent;
  rating: number;
  pubDate?: string;
}

// --- API Client ---

export class NewsDeframerClient {
  constructor(private config: Settings) {}

  private async proxyRequest<T>(endpoint: string, params: Record<string, string>): Promise<T | null> {
    const headers: Record<string, string> = {};
    if (this.config.username && this.config.password) {
      headers['Authorization'] = 'Basic ' + btoa(`${this.config.username}:${this.config.password}`);
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
            resolve(null); // Not found is not an error
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
}
