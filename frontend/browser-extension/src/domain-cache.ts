const CACHE_KEY = 'domain_cache';
const CACHE_DURATION = 3 * 60 * 1000; // 3 minutes

export interface CachedDomainEntry {
  domain: string;
  language: string;
}

interface DomainCache {
  domains: CachedDomainEntry[];
  timestamp: number;
}

export const getCachedDomains = (): Promise<CachedDomainEntry[] | null> => {
  return new Promise((resolve) => {
    chrome.storage.local.get([CACHE_KEY], (result) => {
      const cached = result[CACHE_KEY] as DomainCache | undefined;
      if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        resolve(cached.domains);
      } else {
        resolve(null);
      }
    });
  });
};

export const setCachedDomains = (domains: CachedDomainEntry[]): Promise<void> => {
  return new Promise((resolve) => {
    const cache: DomainCache = {
      domains,
      timestamp: Date.now(),
    };
    chrome.storage.local.set({ [CACHE_KEY]: cache }, () => {
      resolve();
    });
  });
};

export const invalidateDomainCache = (): Promise<void> => {
  return new Promise((resolve) => {
    const cache: DomainCache = {
      domains: [],
      timestamp: 0,
    };
    chrome.storage.local.set({ [CACHE_KEY]: cache }, () => {
      resolve();
    });
  });
};
