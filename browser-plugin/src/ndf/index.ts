import log from '../shared/logger';
import { Settings } from '../shared/settings';
import { NewsDeframerClient } from './client';
import { handleArticle } from './handler-article';
import { handlePortal } from './handler-portal';

export default {
  start: async (config: Settings) => {
    try {
      const client = new NewsDeframerClient(config);
      const domains = await client.getDomains();
      const currentUrl = new URL(window.location.href);

      const isMonitored = domains.some((d) => currentUrl.hostname === d || currentUrl.hostname.endsWith('.' + d));

      if (!isMonitored) {
        return; // Not a monitored domain, do nothing.
      }

      // Determine if it's a portal or an article page.
      const pathSegments = currentUrl.pathname.split('/').filter((p) => p.length > 0);

      if (pathSegments.length <= 1) {
        handlePortal();
      } else {
        handleArticle();
      }
    } catch (e) {
      log.error('NDF Error:', e);
    }
  },
};
