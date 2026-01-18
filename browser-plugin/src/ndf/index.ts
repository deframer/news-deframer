import { getCachedDomains } from '../shared/domain-cache';
import log from '../shared/logger';
import { Settings } from '../shared/settings';
import { NewsDeframerClient } from './client';
import { handleArticle } from './handler-article';
import { handlePortal } from './handler-portal';

export default {
  start: async (config: Settings) => {
    try {
      // Check for bypass flag (set if we stopped a page that turned out to be invalid)
      if (sessionStorage.getItem('ndf-bypass')) {
        log.info('NDF Bypass flag detected. Skipping Deframer for this load.');
        sessionStorage.removeItem('ndf-bypass');
        return;
      }

      const client = new NewsDeframerClient(config);
      
      // Phase 1: Check Local Cache (Fast Path)
      let domains = await getCachedDomains();
      const currentUrl = new URL(window.location.href);

      if (!domains) {
         // Phase 2: Cache Miss (Slow Path)
         // We do not stop loading here to avoid blocking random pages.
         // We fetch domains in the background.
         domains = await client.getDomains();
         
         if (domains.length === 0) {
           log.error('Could not get domains from backend, aborting.');
           return;
         }

         const isMonitored = domains.some((d) => currentUrl.hostname === d || currentUrl.hostname.endsWith('.' + d));

         if (isMonitored) {
           // We are on a monitored domain, but we missed the start.
           // Reload to trigger the Fast Path (Phase 1).
           log.info('Monitored domain detected after load. Reloading to apply Deframer...');
           window.location.reload();
         }
         return;
      }

      // Phase 1 Continue: Cache Hit
      const isMonitored = domains.some((d) => currentUrl.hostname === d || currentUrl.hostname.endsWith('.' + d));

      if (!isMonitored) {
        return; // Not a monitored domain, do nothing.
      }

      // Determine if it's a portal or an article page.
      const pathSegments = currentUrl.pathname.split('/').filter((p) => p.length > 0);

      if (pathSegments.length <= 1) {
        handlePortal(client);
      } else {
        handleArticle(client);
      }
    } catch (e) {
      log.error('NDF Error:', e);
    }
  },
};
