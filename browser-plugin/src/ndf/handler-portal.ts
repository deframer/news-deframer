import { getDomain } from 'tldts';

import log from '../shared/logger';
import { AnalyzedItem,NewsDeframerClient } from './client';
import { formatRatingPercent, getRatingColors } from './ratings';

const createTilesHtml = (items: AnalyzedItem[], rootDomain: string): string => {
  const tiles = items
    .map((item) => {
      const title = item.title_corrected || item.title_original || 'No title';
      const description = item.description_corrected || item.description_original || 'No description';
      const ratingValue = formatRatingPercent(item.rating);
      const ratingColors = getRatingColors(ratingValue);
      const overallReason = item.overall_reason || 'No reason provided.';

      const barHtml = `
        <div class="bar-container" role="meter" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${ratingValue}" aria-label="Overall rating: ${ratingValue}%. Reason: ${overallReason}">
          <div class="bar" style="width: ${ratingValue}%; background-color: ${ratingColors.bg};"></div>
          <div class="bar-overlay" style="color: ${ratingColors.text}; ${ratingColors.text === '#ffffff' ? 'text-shadow: 0 0 3px rgba(0,0,0,0.7);' : ''}">
            <span>${ratingValue}%</span>
          </div>
          <div class="tooltip-text">${overallReason}</div>
        </div>
      `;

      const imageHtml =
        item.media && item.media.medium === 'image' && item.media.url
          ? `<div class="image-container">
               <img src="${item.media.url}" alt="${item.media.description || ''}" style="width: 100%; height: auto; display: block;">
               ${barHtml}
             </div>`
          : barHtml;

      return `
        <a href="${item.url}" class="tile-link">
          <div class="tile">
            ${imageHtml}
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
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <style>
          body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            margin: 0;
            background-color: #f0f2f5;
            padding-top: 60px; /* Space for sticky header */
          }
          .container { padding: 2em; }
          .header {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            background-color: #fff;
            padding: 1em;
            z-index: 1000;
            text-align: center;
            margin-bottom: 0;
            box-shadow: 0 2px 5px rgba(0,0,0,0.1);
            box-sizing: border-box;
          }
          .header h1 {
            margin: 0;
            font-size: 1.5em;
          }
          .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }
          .tile-link { text-decoration: none; color: inherit; }
          .tile { background-color: #fff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; height: 100%; position: relative; }
          .tile:hover { transform: translateY(-5px); box-shadow: 0 8px 12px rgba(0,0,0,0.15); }
          .image-container { position: relative; }
          .content { padding: 15px; }
          h3 { margin: 0 0 10px; font-size: 1.1em; }
          p { font-size: 0.9em; color: #666; }
          .bar-container {
            background-color: #e9ecef;
            height: 30px;
            width: 100%;
            position: relative; /* For overlay and tooltip */
          }
          .image-container .bar-container {
            position: absolute;
            top: 0;
            left: 0;
          }
          .bar { height: 100%; }
          .bar-overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            padding-left: 10px;
            font-weight: bold;
            font-size: 0.9em;
            pointer-events: none;
          }
          .bar-container .tooltip-text {
            visibility: hidden;
            opacity: 0;
            width: 250px;
            background-color: rgba(0,0,0,0.85);
            color: #fff;
            text-align: left;
            border-radius: 6px;
            padding: 10px;
            position: absolute;
            z-index: 1;
            top: 100%;
            left: 50%;
            transform: translateX(-50%);
            transition: opacity 0.2s;
            font-size: 0.9em;
            line-height: 1.4;
            pointer-events: none;
          }
          .bar-container:hover .tooltip-text {
            visibility: visible;
            opacity: 1;
          }
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
  log.info('Portal page detected. Stopping window immediately.');
  window.stop();
  const rootDomain = getDomain(window.location.hostname);

  if (!rootDomain) {
    log.error('Could not determine root domain. Reloading with bypass.');
    sessionStorage.setItem('ndf-bypass', 'true');
    window.location.reload();
    return;
  }

  try {
    const items = await client.getSite(rootDomain);
    if (items.length > 0) {
      log.info(`Successfully fetched ${items.length} items for ${rootDomain}.`);
      document.documentElement.innerHTML = createTilesHtml(items, rootDomain);
    } else {
      log.info(`No items found for ${rootDomain}. Reloading with bypass.`);
      sessionStorage.setItem('ndf-bypass', 'true');
      window.location.reload();
    }
  } catch (error) {
    log.error(`Failed to fetch items for ${rootDomain}:`, error);
    sessionStorage.setItem('ndf-bypass', 'true');
    window.location.reload();
  }
};
