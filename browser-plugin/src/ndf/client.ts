import { getCachedDomains, invalidateDomainCache,setCachedDomains } from '../shared/domain-cache';
import { Settings } from '../shared/settings';

// --- Type Definitions based on Go backend models ---

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
  updated_at?: string;
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

  async getDomains(): Promise<string[]> {
    const cached = await getCachedDomains();
    if (cached) {
      return cached;
    }

    try {
      const result = await this.proxyRequest<string[]>('/api/domains', {});
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
}
