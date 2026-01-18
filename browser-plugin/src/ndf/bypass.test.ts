import {
  deleteBypassForDomain,
  deleteBypassForUrl,
  hasBypassForDomain,
  hasBypassForUrl,
  setBypassForDomain,
  setBypassForUrl,
} from './bypass';

// Mock localStorage
const localStorageMock = (() => {
  let store: { [key: string]: string } = {};
  return {
    getItem: (key: string) => store[key] || null,
    setItem: (key: string, value: string) => {
      store[key] = value.toString();
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('Bypass Management', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  describe('URL Bypasses', () => {
    const testUrl = 'https://example.com/article-1';
    const anotherUrl = 'https://example.com/article-2';

    it('should not have a bypass for a URL initially', () => {
      expect(hasBypassForUrl(testUrl)).toBe(false);
    });

    it('should set and check a bypass for a URL', () => {
      setBypassForUrl(testUrl);
      expect(hasBypassForUrl(testUrl)).toBe(true);
    });

    it('should handle multiple URLs', () => {
      setBypassForUrl(testUrl);
      setBypassForUrl(anotherUrl);
      expect(hasBypassForUrl(testUrl)).toBe(true);
      expect(hasBypassForUrl(anotherUrl)).toBe(true);
    });

    it('should delete a bypass for a URL', () => {
      setBypassForUrl(testUrl);
      setBypassForUrl(anotherUrl);
      deleteBypassForUrl(testUrl);
      expect(hasBypassForUrl(testUrl)).toBe(false);
      expect(hasBypassForUrl(anotherUrl)).toBe(true);
    });
  });

  describe('Domain Bypasses', () => {
    const testDomain = 'example.com';
    const anotherDomain = 'anotherexample.com';

    it('should not have a bypass for a domain initially', () => {
      expect(hasBypassForDomain(testDomain)).toBe(false);
    });

    it('should set and check a bypass for a domain', () => {
      setBypassForDomain(testDomain);
      expect(hasBypassForDomain(testDomain)).toBe(true);
    });

    it('should delete a bypass for a domain', () => {
      setBypassForDomain(testDomain);
      setBypassForDomain(anotherDomain);
      deleteBypassForDomain(testDomain);
      expect(hasBypassForDomain(testDomain)).toBe(false);
      expect(hasBypassForDomain(anotherDomain)).toBe(true);
    });
  });
});