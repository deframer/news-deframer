import log from '../shared/logger';
import { NewsDeframerClient } from './client';

export const handleArticle = async (client: NewsDeframerClient) => {
  log.info('Article page detected.');
  document.body.style.border = '15px solid blue'; // Visual proof for article
  try {
    const item = await client.getItem(window.location.href);
    if (item) {
      log.info('Successfully fetched item:', item);
    } else {
      log.info('No item found for this URL.');
    }
  } catch (error) {
    log.error('Failed to fetch item:', error);
  }
};
