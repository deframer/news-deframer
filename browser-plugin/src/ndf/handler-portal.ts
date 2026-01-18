import { getDomain } from 'tldts';

import log from '../shared/logger';
import { NewsDeframerClient } from './client';

export const handlePortal = async (client: NewsDeframerClient) => {
  log.info('Portal page detected.');
  document.body.style.border = '15px solid green'; // Visual proof for portal
  const rootDomain = getDomain(window.location.hostname);

  if (!rootDomain) {
    log.error('Could not determine root domain.');
    return;
  }

  try {
    const items = await client.getSite(rootDomain);
    if (items.length > 0) {
      log.info(`Successfully fetched ${items.length} items for ${rootDomain}:`, items);
    } else {
      log.info(`No items found for ${rootDomain}.`);
    }
  } catch (error) {
    log.error(`Failed to fetch items for ${rootDomain}:`, error);
  }
};
