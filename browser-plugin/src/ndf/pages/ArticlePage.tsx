import React, { useState } from 'react';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { RatingBar } from '../components/RatingBar';

// All styles are encapsulated here. A Shadow DOM will prevent them from leaking.
const articlePageCss = `
  /* All styles from createArticleHtml are moved here */
  .page-header { background-color: #fff; padding: 8px 1.5em; text-align: left; border-bottom: 1px solid #eee; z-index: 1001; display: flex; align-items: center; justify-content: space-between; }
  .btn-back { display: inline-flex; align-items: center; gap: 6px; padding: 8px 12px; border: 1px solid #007bff; border-radius: 8px; background-color: #007bff; color: #fff; font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
  .btn-back:hover { background-color: #0069d9; }
  .btn-hide { padding: 8px 12px; border: 1px solid #ddd; border-radius: 8px; background-color: #fff; color: #333; font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; }
  .btn-hide:hover { background-color: #f8f9fa; }
  .container { max-width: 800px; margin: 0; background-color: #fff; border-radius: 0; box-shadow: none; padding-bottom: 20px; }
  .image-container img { width: 100%; height: auto; border-radius: 0; display: block; }
  .main-content { padding: 1.5em; }
  h1 { margin: 0 0 10px; font-size: 2em; }
  .description { font-size: 1.1em; color: #555; margin-bottom: 2em; }
  .analysis-section { border-top: 2px solid #eee; padding-top: 1.5em; margin-top: 1.5em; }
  .metric-item { display: block; margin-bottom: 1.5em; }
  .metric-label { font-weight: bold; font-size: 1.1em; margin-bottom: 5px; display: block; }
  .metric-details { display: block; width: 100%; }
  .reason { margin: 5px 0 0; font-size: 0.95em; color: #666; line-height: 1.4; }
  .original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid #ddd; }
  .action-buttons { display: flex; justify-content: space-around; gap: 10px; padding: 1.5em; background-color: #fff; border-top: 1px solid #eee; }
  .btn { flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px; padding: 10px; border: 1px solid #ddd; border-radius: 8px; background-color: #fff; color: #333; font-size: 0.95em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
  .btn:hover { background-color: #f8f9fa; }
  .btn-primary { background-color: #007bff; color: #fff; border-color: #007bff; }
  .btn-primary:hover { background-color: #0069d9; }

  @media (min-width: 800px) {
    .page-header { border-radius: 12px 12px 0 0; }
    .container { margin: 2em auto; border-radius: 12px; box-shadow: 0 6px 12px rgba(0,0,0,0.1); }
    .metric-item { display: grid; grid-template-columns: 140px 1fr; gap: 0 20px; align-items: start; }
    .metric-label { margin-bottom: 0; padding-top: 4px; }
    .action-buttons { border-radius: 0 0 12px 12px; justify-content: center; }
    .action-buttons .btn { flex: 0 0 180px; }
  }

  @media (max-width: 799px) {
    body { padding-top: 55px; padding-bottom: 80px; }
    .page-header { position: fixed; top: 0; left: 0; width: 100%; box-sizing: border-box; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .action-buttons { position: fixed; bottom: 0; left: 0; width: 100%; z-index: 1000; box-shadow: 0 -2px 10px rgba(0,0,0,0.1); margin-top: 0; box-sizing: border-box; }
  }
`;

interface ArticlePageProps {
  item: AnalyzedItem;
}

export const ArticlePage = ({ item }: ArticlePageProps) => {
  const [showDetails, setShowDetails] = useState(false);
  const [showOriginal, setShowOriginal] = useState(false);

  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;
  const title = item.title_corrected || item.title_original || 'No title';
  const description = item.description_corrected || item.description_original || 'No description';
  const imageUrl = item.media?.medium === 'image' ? item.media.url : '';

  const metrics = [
    { id: 'framing', label: 'Framing', value: item.framing, reason: item.framing_reason },
    { id: 'clickbait', label: 'Clickbait', value: item.clickbait, reason: item.clickbait_reason },
    { id: 'persuasive', label: 'Persuasive', value: item.persuasive, reason: item.persuasive_reason },
    { id: 'hyper_stimulus', label: 'Hyper Stimulus', value: item.hyper_stimulus, reason: item.hyper_stimulus_reason },
    { id: 'speculative', label: 'Speculative', value: item.speculative, reason: item.speculative_reason },
  ];
  
  const bypassAndReload = () => {
    window.scrollTo(0, 0);
    log.info('Bypassing for this session and reloading.');
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
  };

  return (
    <>
      <style>{articlePageCss}</style>
      <div className="container">
        <header className="page-header">
          <a href="/" className="btn-back" title={`Go back to ${rootDomain} portal`}>Back</a>
          <button onClick={bypassAndReload} className="btn-hide">Hide</button>
        </header>

        {imageUrl && (
          <div className="image-container">
            <img src={imageUrl} alt={title} />
          </div>
        )}

        <div className="main-content">
          <h1>{title}</h1>
          <p className="description">{description}</p>
          
          <div className="analysis-section">
            <div className="metric-item">
                <RatingBar value={item.rating} label="Overall Rating" reason={item.overall_reason} />
            </div>

            {showDetails && (
              <div id="metrics-content">
                {metrics.map(m => (
                  <div className="metric-item" key={m.id}>
                    <RatingBar value={m.value} label={m.label} reason={m.reason} id={m.id} />
                  </div>
                ))}
              </div>
            )}

            {showOriginal && (
               <div id="original-content" className="original-content">
                 <h3>Original</h3>
                 <h4>{item.title_original || ''}</h4>
                 <p>{item.description_original || ''}</p>
               </div>
            )}
          </div>

          <Footer />
        </div>

        <div className="action-buttons">
          <button onClick={() => setShowOriginal(true)} className="btn btn-primary">Original Title</button>
          <button onClick={() => { setShowDetails(true); setShowOriginal(true); }} className="btn">Details</button>
          <button onClick={bypassAndReload} className="btn">View Original</button>
        </div>
      </div>
    </>
  );
};
