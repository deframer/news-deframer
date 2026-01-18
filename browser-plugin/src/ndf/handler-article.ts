import log from '../shared/logger';
import { NewsDeframerClient, AnalyzedItem } from './client';

const createArticleHtml = (item: AnalyzedItem): string => {
  const title = item.title_corrected || item.title_original || 'No title';
  const description = item.description_corrected || item.description_original || 'No description';
  const image =
    item.media && item.media.medium === 'image' && item.media.url
      ? `<div class="image-container"><img src="${item.media.url}" alt="${item.media.description || ''}" style="width: 100%; height: auto; border-radius: 8px 8px 0 0;"></div>`
      : '';

  const analysisDetails = `
    <div class="analysis-grid">
      <div><strong>Framing:</strong> ${item.framing?.toFixed(2)}</div>
      <div>${item.framing_reason || ''}</div>
      <div><strong>Clickbait:</strong> ${item.clickbait?.toFixed(2)}</div>
      <div>${item.clickbait_reason || ''}</div>
      <div><strong>Persuasive:</strong> ${item.persuasive?.toFixed(2)}</div>
      <div>${item.persuasive_reason || ''}</div>
      <div><strong>Hyper Stimulus:</strong> ${item.hyper_stimulus?.toFixed(2)}</div>
      <div>${item.hyper_stimulus_reason || ''}</div>
      <div><strong>Speculative:</strong> ${item.speculative?.toFixed(2)}</div>
      <div>${item.speculative_reason || ''}</div>
    </div>
    <h4>Overall Reason</h4>
    <p>${item.overall_reason || ''}</p>
  `;

  return `
    <html>
      <head>
        <title>News Deframer: ${title}</title>
        <style>
          body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; margin: 0; background-color: #f0f2f5; }
          .container { max-width: 800px; margin: 2em auto; padding: 2em; background-color: #fff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); position: relative; }
          .rating { position: absolute; top: 10px; right: 10px; background-color: rgba(0,0,0,0.7); color: #fff; padding: 5px 8px; border-radius: 4px; font-size: 0.9em; font-weight: bold; z-index: 1; }
          .card { margin-bottom: 2em; }
          .image-container { margin-bottom: 1em; }
          h1 { margin: 0 0 10px; font-size: 1.8em; }
          p { font-size: 1.1em; color: #333; }
          .analysis-section { margin-top: 2em; }
          .analysis-grid { display: grid; grid-template-columns: 1fr 2fr; gap: 10px; margin-bottom: 1em; }
          .original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid #ddd; }
        </style>
      </head>
      <body>
        <div class="container">
          <div class="rating">${item.rating.toFixed(2)}</div>
          <div class="card">
            ${image}
            <h1>${title}</h1>
            <p>${description}</p>
          </div>
          <div class="analysis-section">
            <h2>Analysis</h2>
            ${analysisDetails}
          </div>
          <div class="original-content">
            <h3>Original</h3>
            <h4>${item.title_original || ''}</h4>
            <p>${item.description_original || ''}</p>
          </div>
        </div>
      </body>
    </html>
  `;
};

export const handleArticle = async (client: NewsDeframerClient) => {
  log.info('Article page detected.');

  try {
    const item = await client.getItem(window.location.href);
    if (item) {
      window.stop();
      log.info('Successfully fetched item:', item);
      document.documentElement.innerHTML = createArticleHtml(item);
    } else {
      log.info('No item found for this URL.');
    }
  } catch (error) {
    log.error('Failed to fetch item:', error);
  }
};

