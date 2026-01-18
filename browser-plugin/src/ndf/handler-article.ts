import log from '../shared/logger';
import { AnalyzedItem, NewsDeframerClient } from './client';

const formatRatingPercent = (rating: number | undefined): number => Math.round((rating || 0.0) * 100);

const getColorForRating = (percentage: number): string => {
  if (percentage < 34) return '#28a745'; // Green for low/neutral
  if (percentage < 67) return '#ffc107'; // Yellow for medium
  return '#dc3545'; // Red for high/bad
};

const createArticleHtml = (item: AnalyzedItem): string => {
  const title = item.title_corrected || item.title_original || 'No title';
  const description = item.description_corrected || item.description_original || 'No description';
  const imageUrl = item.media && item.media.medium === 'image' && item.media.url ? item.media.url : '';

  const metrics = [
    { id: 'framing', label: 'Framing', value: formatRatingPercent(item.framing), reason: item.framing_reason || 'No reason provided.' },
    { id: 'clickbait', label: 'Clickbait', value: formatRatingPercent(item.clickbait), reason: item.clickbait_reason || 'No reason provided.' },
    { id: 'persuasive', label: 'Persuasive', value: formatRatingPercent(item.persuasive), reason: item.persuasive_reason || 'No reason provided.' },
    { id: 'hyper_stimulus', label: 'Hyper Stimulus', value: formatRatingPercent(item.hyper_stimulus), reason: item.hyper_stimulus_reason || 'No reason provided.' },
    { id: 'speculative', label: 'Speculative', value: formatRatingPercent(item.speculative), reason: item.speculative_reason || 'No reason provided.' },
  ];

  const overallValue = formatRatingPercent(item.rating);
  const overallReason = item.overall_reason || 'No overall reason provided.';

  return `
    <html>
      <head>
        <title>News Deframer: ${title}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <style>
          body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; margin: 0; background-color: #f0f2f5; color: #333; }
          .container { max-width: 800px; margin: 2em auto; background-color: #fff; border-radius: 12px; box-shadow: 0 6px 12px rgba(0,0,0,0.1); }
          .image-container img { width: 100%; height: auto; border-radius: 12px 12px 0 0; display: block; }
          .main-content { padding: 1.5em; }
          h1 { margin: 0 0 10px; font-size: 2em; }
          .description { font-size: 1.1em; color: #555; margin-bottom: 2em; }
          
          .analysis-section { border-top: 2px solid #eee; padding-top: 1.5em; margin-top: 1.5em; }

          /* MOBILE FIRST (Default) */
          .metric-item {
            display: block;
            margin-bottom: 1.5em;
          }
          
          .metric-label {
            font-weight: bold;
            font-size: 1.1em;
            margin-bottom: 5px; 
            display: block;
          }

          .metric-details {
            display: block;
            width: 100%;
          }

          .bar-container {
            background-color: #e9ecef;
            border-radius: 5px;
            height: 10px;
            width: 100%;
            margin-bottom: 5px; 
          }
          
          .bar { height: 100%; border-radius: 5px; }

          .reason {
            margin: 0;
            font-size: 0.95em;
            color: #666;
            line-height: 1.4;
          }

          .metric-separator { width: 95%; margin: 2em auto; border: 0; border-top: 1px solid #eee; }
          .original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid #ddd; }

          /* DESKTOP LAYOUT (> 800px) */
          @media (min-width: 800px) {
            .metric-item {
              display: grid;
              grid-template-columns: 140px 1fr;
              gap: 0 20px;
              align-items: start;
            }

            .metric-label {
              margin-bottom: 0;
              padding-top: 0;
            }
          }
        </style>
      </head>
      <body>
        <div class="container">
          ${imageUrl ? `<div class="image-container"><img src="${imageUrl}" alt="${title}"></div>` : ''}
          <div class="main-content">
            <h1>${title}</h1>
            <p class="description">${description}</p>
            <div class="analysis-section">
              ${metrics.map(m => `
                <div class="metric-item">
                  <div class="metric-label">${m.label}</div>
                  <div class="metric-details">
                    <div class="bar-container" role="meter" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${m.value}" aria-label="${m.label} rating: ${m.value}%" aria-describedby="${m.id}-reason">
                      <div class="bar" style="width: ${m.value}%; background-color: ${getColorForRating(m.value)};"></div>
                    </div>
                    <p id="${m.id}-reason" class="reason">${m.reason}</p>
                  </div>
                </div>
              `).join('')}
              
              <hr class="metric-separator">

              <div class="metric-item">
                <div class="metric-label">Overall Rating</div>
                <div class="metric-details">
                  <div class="bar-container" role="meter" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${overallValue}" aria-label="Overall rating: ${overallValue}%" aria-describedby="overall-reason">
                    <div class="bar" style="width: ${overallValue}%; background-color: ${getColorForRating(overallValue)};"></div>
                  </div>
                  <p id="overall-reason" class="reason">${overallReason}</p>
                </div>
              </div>
            </div>
            <div class="original-content">
              <h3>Original</h3>
              <h4>${item.title_original || ''}</h4>
              <p>${item.description_original || ''}</p>
            </div>
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
      log.info('Successfully fetched item.');
      document.documentElement.innerHTML = createArticleHtml(item);
    } else {
      log.info('No item found for this URL.');
    }
  } catch (error) {
    log.error('Failed to fetch item:', error);
  }
};
