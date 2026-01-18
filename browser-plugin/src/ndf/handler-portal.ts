import { getDomain } from 'tldts';

import log from '../shared/logger';
import { AnalyzedItem,NewsDeframerClient } from './client';

const formatRating = (rating: number | undefined): string => (rating || 0.0).toFixed(2);

const createTilesHtml = (items: AnalyzedItem[], rootDomain: string): string => {
  const tiles = items
    .map((item) => {
      const title = item.title_corrected || item.title_original || 'No title';
      const description = item.description_corrected || item.description_original || 'No description';
      const image =
        item.media && item.media.medium === 'image' && item.media.url
          ? `<div class="image-container"><img src="${item.media.url}" alt="${item.media.description || ''}" style="width: 100%; height: auto;"></div>`
          : '';

      return `
        <a href="${item.url}" class="tile-link">
          <div class="tile">
            <div class="rating">${formatRating(item.rating)}</div>
            ${image}
            <div class="content">
              <h3>${title}</h3>
              <p>${description}</p>
            </div>
          </div>
        </a>
      `;
    })
    .join('');

  return `
    <html>
      <head>
        <title>News Deframer: ${rootDomain}</title>
        <style>
          body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; margin: 0; background-color: #f0f2f5; }
          .container { padding: 2em; }
          .header { text-align: center; margin-bottom: 2em; }
          .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }
          .tile-link { text-decoration: none; color: inherit; }
          .tile { background-color: #fff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; height: 100%; position: relative; }
          .tile:hover { transform: translateY(-5px); box-shadow: 0 8px 12px rgba(0,0,0,0.15); }
          .rating { position: absolute; top: 10px; right: 10px; background-color: rgba(0,0,0,0.7); color: #fff; padding: 5px 8px; border-radius: 4px; font-size: 0.9em; font-weight: bold; z-index: 1; }
          .image-container { padding: 15px 15px 0; }
          .content { padding: 15px; }
          h3 { margin: 0 0 10px; font-size: 1.1em; }
          p { font-size: 0.9em; color: #666; }
        </style>
      </head>
      <body>
        <div class="container">
          <div class="header">
            <h1>News Deframer: ${rootDomain}</h1>
          </div>
          <div class="grid">
            ${tiles}
          </div>
        </div>
      </body>
    </html>
  `;
};

export const handlePortal = async (client: NewsDeframerClient) => {
  log.info('Portal page detected.');
  const rootDomain = getDomain(window.location.hostname);

  if (!rootDomain) {
    log.error('Could not determine root domain.');
    return;
  }

  try {
    const items = await client.getSite(rootDomain);
    if (items.length > 0) {
      window.stop();
      log.info(`Successfully fetched ${items.length} items for ${rootDomain}.`);
      document.documentElement.innerHTML = createTilesHtml(items, rootDomain);
    } else {
      log.info(`No items found for ${rootDomain}.`);
    }
  } catch (error) {
    log.error(`Failed to fetch items for ${rootDomain}:`, error);
  }
};

