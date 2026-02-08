import { classifyUrl, PageType } from './url-classifier';

describe('URL Classifier', () => {
  test('should classify root domain as PORTAL', () => {
    const url = new URL('https://www.example.com/');
    expect(classifyUrl(url)).toBe(PageType.PORTAL);
  });

  test('should classify single segment path as ARTICLE', () => {
    const url = new URL('https://www.example.com/politics');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should classify multi-segment path as ARTICLE', () => {
    const url = new URL('https://www.example.com/politics/2023/some-news');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should ignore trailing slashes', () => {
    const url = new URL('https://www.example.com/politics/');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should handle query parameters', () => {
    const url = new URL('https://www.example.com/politics?sort=newest');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should classify winfuture news link as ARTICLE', () => {
    const url = new URL('https://www.example.com/news,156261.html');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should classify index.html as ARTICLE', () => {
    const url = new URL('https://www.example.com/index.html');
    expect(classifyUrl(url)).toBe(PageType.ARTICLE);
  });

  test('should classify specific portal url as PORTAL', () => {
    const url = new URL('https://www.example.com/news');
    const portalUrl = 'example.com/news';
    expect(classifyUrl(url, portalUrl)).toBe(PageType.PORTAL);
  });

  test('should classify subpath of portal url as ARTICLE', () => {
    const url = new URL('https://www.example.com/news/article-123');
    const portalUrl = 'example.com/news';
    expect(classifyUrl(url, portalUrl)).toBe(PageType.ARTICLE);
  });
});