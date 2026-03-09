export enum PageType {
  PORTAL = 'portal',
  ARTICLE = 'article',
}

export const classifyUrl = (url: URL, portalUrl?: string | null): PageType => {
  if (portalUrl) {
    const currentHost = url.host.replace(/^www\./, '');
    const currentUrl = (currentHost + url.pathname).replace(/\/$/, '');

    let targetUrl = portalUrl.replace(/^https?:\/\//, '').replace(/\/$/, '');
    targetUrl = targetUrl.replace(/^www\./, '');

    if (currentUrl === targetUrl) {
      return PageType.PORTAL;
    }
  }

  const pathSegments = url.pathname.split('/').filter((p) => p.length > 0);

  // If there are any path segments, treat as article
  if (pathSegments.length > 0) {
    return PageType.ARTICLE;
  }

  // Root domain (no path) is portal
  return PageType.PORTAL;
};