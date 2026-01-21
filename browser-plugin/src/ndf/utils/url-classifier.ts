export enum PageType {
  PORTAL = 'portal',
  ARTICLE = 'article',
}

// TODO: Consider adding server-side endpoint to determine if a URL should be treated as portal vs article
// Currently using simple heuristic: root domain = portal, any path = article
export const classifyUrl = (url: URL): PageType => {
  const pathSegments = url.pathname.split('/').filter((p) => p.length > 0);

  // If there are any path segments, treat as article
  if (pathSegments.length > 0) {
    return PageType.ARTICLE;
  }

  // Root domain (no path) is portal
  return PageType.PORTAL;
};