export enum PageType {
  PORTAL = 'portal',
  ARTICLE = 'article',
}

export const classifyUrl = (url: URL): PageType => {
  const pathSegments = url.pathname.split('/').filter((p) => p.length > 0);

  if (pathSegments.length > 1) {
    return PageType.ARTICLE;
  }

  if (pathSegments.length === 1) {
    const segment = pathSegments[0].toLowerCase();
    if ((segment.endsWith('.html') || segment.endsWith('.htm')) && !segment.startsWith('index.')) {
      return PageType.ARTICLE;
    }
  }

  return PageType.PORTAL;
};