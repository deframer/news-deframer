const BYPASSED_URLS_KEY = '__news_deframer_bypassed_urls';
const BYPASSED_DOMAINS_KEY = '__news_deframer_bypassed_domains';

// Helper to get and parse data safely from localStorage
const getStoredSet = (key: string): Set<string> => {
  try {
    const storedValue = localStorage.getItem(key);
    if (storedValue) {
      const items = JSON.parse(storedValue);
      if (Array.isArray(items)) {
        return new Set(items);
      }
    }
  } catch (e) {
    console.error('Error reading bypass data from localStorage', e);
  }
  return new Set();
};

// Helper to stringify and set data safely to localStorage
const setStoredSet = (key: string, data: Set<string>): void => {
  try {
    localStorage.setItem(key, JSON.stringify(Array.from(data)));
  } catch (e) {
    console.error('Error writing bypass data to localStorage', e);
  }
};

export const hasBypassForUrl = (url: string): boolean => {
  return getStoredSet(BYPASSED_URLS_KEY).has(url);
};

export const hasBypassForDomain = (domain: string): boolean => {
  return getStoredSet(BYPASSED_DOMAINS_KEY).has(domain);
};

export const setBypassForUrl = (url: string): void => {
  const urls = getStoredSet(BYPASSED_URLS_KEY);
  urls.add(url);
  setStoredSet(BYPASSED_URLS_KEY, urls);
};

export const setBypassForDomain = (domain: string): void => {
  const domains = getStoredSet(BYPASSED_DOMAINS_KEY);
  domains.add(domain);
  setStoredSet(BYPASSED_DOMAINS_KEY, domains);
};

export const deleteBypassForUrl = (url: string): void => {
  const urls = getStoredSet(BYPASSED_URLS_KEY);
  urls.delete(url);
  setStoredSet(BYPASSED_URLS_KEY, urls);
};

export const deleteBypassForDomain = (domain: string): void => {
  const domains = getStoredSet(BYPASSED_DOMAINS_KEY);
  domains.delete(domain);
  setStoredSet(BYPASSED_DOMAINS_KEY, domains);
};