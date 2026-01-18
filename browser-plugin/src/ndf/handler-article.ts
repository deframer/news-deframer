import log from '../shared/logger';
import { AnalyzedItem, NewsDeframerClient } from './client';

const formatRatingPercent = (rating: number | undefined): number => Math.round((rating || 0.0) * 100);

const getRatingColors = (percentage: number): { bg: string; text: string } => {
  if (percentage < 34) return { bg: '#198754', text: '#ffffff' }; // Accessible Green
  if (percentage < 67) return { bg: '#ffc107', text: '#000000' }; // Accessible Yellow
  return { bg: '#b02a37', text: '#ffffff' }; // Accessible Red
};

const createArticleHtml = (item: AnalyzedItem): string => {
  const title = item.title_corrected || item.title_original || 'No title';
  const description = item.description_corrected || item.description_original || 'No description';
  const imageUrl = item.media && item.media.medium === 'image' && item.media.url ? item.media.url : '';

  const metrics = [
    { id: 'framing', label: 'Framing', raw: item.framing, value: formatRatingPercent(item.framing), reason: item.framing_reason || 'No reason provided.' },
    { id: 'clickbait', label: 'Clickbait', raw: item.clickbait, value: formatRatingPercent(item.clickbait), reason: item.clickbait_reason || 'No reason provided.' },
    { id: 'persuasive', label: 'Persuasive', raw: item.persuasive, value: formatRatingPercent(item.persuasive), reason: item.persuasive_reason || 'No reason provided.' },
    { id: 'hyper_stimulus', label: 'Hyper Stimulus', raw: item.hyper_stimulus, value: formatRatingPercent(item.hyper_stimulus), reason: item.hyper_stimulus_reason || 'No reason provided.' },
    { id: 'speculative', label: 'Speculative', raw: item.speculative, value: formatRatingPercent(item.speculative), reason: item.speculative_reason || 'No reason provided.' },
  ];

  const overallRaw = item.rating;
  const overallValue = formatRatingPercent(item.rating);
  const overallReason = item.overall_reason || 'No overall reason provided.';
  const overallColors = getRatingColors(overallValue);

  return `
    <html>
      <head>
        <title>News Deframer: ${title}</title>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <style>
          body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; margin: 0; background-color: #f0f2f5; color: #333; }
          
          /* MOBILE FIRST (Full Bleed) */
          .container { 
            max-width: 800px; 
            margin: 0; 
            background-color: #fff; 
            border-radius: 0; 
            box-shadow: none; 
            padding-bottom: 20px; 
          }
          .image-container img { 
            width: 100%; 
            height: auto; 
            border-radius: 0; 
            display: block; 
          }
          
          .main-content { padding: 1.5em; }
          h1 { margin: 0 0 10px; font-size: 2em; }
          .description { font-size: 1.1em; color: #555; margin-bottom: 2em; }
          
          .analysis-section { border-top: 2px solid #eee; padding-top: 1.5em; margin-top: 1.5em; }

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
            height: 30px; 
            width: 100%;
            margin-bottom: 5px;
            position: relative; 
            overflow: hidden; 
          }
          
          .bar { 
            height: 100%; 
            border-radius: 5px; 
            transition: width 0.3s ease;
          }

          .bar-overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            display: flex;
            align-items: center;
            padding-left: 8px; 
            font-weight: bold;
            font-size: 0.95em;
            pointer-events: none; 
          }
          
          .reason {
            margin: 0;
            font-size: 0.95em;
            color: #666;
            line-height: 1.4;
          }

          .original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid #ddd; }

          /* BUTTONS STYLES */
          .action-buttons {
            display: flex;
            justify-content: space-around;
            gap: 10px;
            padding: 15px;
            background-color: #fff;
            border-top: 1px solid #eee;
            margin-top: 2em;
          }

          .btn {
            flex: 1;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 6px;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background-color: #fff;
            color: #333;
            font-size: 0.95em;
            font-weight: 500;
            cursor: pointer;
            transition: background-color 0.2s;
            text-decoration: none;
          }

          .btn:hover { background-color: #f8f9fa; }
          .btn-primary { background-color: #007bff; color: #fff; border-color: #007bff; }
          .btn-primary:hover { background-color: #0069d9; }

          /* DESKTOP LAYOUT (> 800px) */
          @media (min-width: 800px) {
            .container { 
              margin: 2em auto; 
              border-radius: 12px; 
              box-shadow: 0 6px 12px rgba(0,0,0,0.1); 
            }
            .image-container img { 
              border-radius: 12px 12px 0 0; 
            }

            .metric-item {
              display: grid;
              grid-template-columns: 140px 1fr;
              gap: 0 20px;
              align-items: start;
            }

            .metric-label {
              margin-bottom: 0;
              padding-top: 4px;
            }
            
            .action-buttons {
              border-radius: 0 0 12px 12px; 
            }
          }

          /* MOBILE STICKY BUTTONS (< 800px) */
          @media (max-width: 799px) {
            body { padding-bottom: 80px; } /* Prevent content from being hidden behind sticky bar */
            .action-buttons {
              position: fixed;
              bottom: 0;
              left: 0;
              width: 100%;
              z-index: 1000;
              box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
              margin-top: 0;
              box-sizing: border-box; /* Ensure padding doesn't overflow width */
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
              <div id="overall-rating-container" class="metric-item">
                <div class="metric-label">Overall Rating</div>
                <div class="metric-details">
                  <div class="bar-container" role="meter" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${overallValue}" aria-label="Overall rating: ${overallValue}%" aria-describedby="overall-reason">
                    <div class="bar" style="width: ${overallValue}%; background-color: ${overallColors.bg};"></div>
                    ${overallRaw !== undefined ? `
                      <div class="bar-overlay" style="color: ${overallColors.text};">
                        <span>${overallValue}%</span>
                      </div>
                    ` : ''}
                  </div>
                  <p id="overall-reason" class="reason">${overallReason}</p>
                </div>
              </div>

              <div id="metrics-content" style="display: none;">
                ${metrics.map(m => {
                  const colors = getRatingColors(m.value);
                  return `
                  <div class="metric-item">
                    <div class="metric-label">${m.label}</div>
                    <div class="metric-details">
                      <div class="bar-container" role="meter" aria-valuemin="0" aria-valuemax="100" aria-valuenow="${m.value}" aria-label="${m.label} rating: ${m.value}%" aria-describedby="${m.id}-reason">
                        <div class="bar" style="width: ${m.value}%; background-color: ${colors.bg};"></div>
                        ${m.raw !== undefined ? `
                          <div class="bar-overlay" style="color: ${colors.text};">
                            <span>${m.value}%</span>
                          </div>
                        ` : ''}
                      </div>
                      <p id="${m.id}-reason" class="reason">${m.reason}</p>
                    </div>
                  </div>
                `}).join('')}
              </div>
            </div>
            <div id="original-content" class="original-content" style="display: none;">
              <h3>Original</h3>
              <h4>${item.title_original || ''}</h4>
              <p>${item.description_original || ''}</p>
            </div>
          </div>
          
          <div class="action-buttons">
            <button id="btn-details" class="btn btn-primary">Details</button>
            <button id="btn-original" class="btn">Original Title</button>
            <a href="#" class="btn">View anyway</a>
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

      // Attach event listeners programmatically after HTML is injected
      const btnDetails = document.getElementById('btn-details');
      if (btnDetails) {
        btnDetails.addEventListener('click', () => {
          const metrics = document.getElementById('metrics-content');
          const original = document.getElementById('original-content');
          const overallRating = document.getElementById('overall-rating-container');
          
          if (metrics && original) {
            metrics.style.display = 'block';
            original.style.display = 'block';
            
            if (overallRating) {
              const rect = overallRating.getBoundingClientRect();
              const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
              const offset = 20; // Graceful offset
              window.scrollTo({ top: rect.top + scrollTop - offset, behavior: 'smooth' });
            }
          }
        });
      }

      const btnOriginal = document.getElementById('btn-original');
      if (btnOriginal) {
        btnOriginal.addEventListener('click', () => {
          const metrics = document.getElementById('metrics-content');
          const original = document.getElementById('original-content');
          const overallRating = document.getElementById('overall-rating-container');

          if (metrics && original) {
            metrics.style.display = 'none';
            original.style.display = 'block';
            
            if (overallRating) {
              const rect = overallRating.getBoundingClientRect();
              const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
              const offset = 20; // Graceful offset
              window.scrollTo({ top: rect.top + scrollTop - offset, behavior: 'smooth' });
            }
          }
        });
      }

    } else {
      log.info('No item found for this URL.');
    }
  } catch (error) {
    log.error('Failed to fetch item:', error);
  }
};
