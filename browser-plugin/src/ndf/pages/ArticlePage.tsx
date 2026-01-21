import '../../shared/i18n';

import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { MetaData } from '../components/MetaData';
import { RatingBar } from '../components/RatingBar';

// All styles are encapsulated here. A Shadow DOM will prevent them from leaking.
const articlePageCss = `
  :host {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    display: block;
    background-color: var(--bg-color);
    color: var(--text-color);
    min-height: 100vh;
  }

  /* All styles from createArticleHtml are moved here */
  .page-header { background-color: var(--header-bg); color: var(--text-color); padding: 8px 1.5em; text-align: left; border-bottom: 1px solid var(--border-color); z-index: 1001; display: flex; align-items: center; justify-content: space-between; }
  .btn-back { display: inline-flex; align-items: center; gap: 6px; padding: 8px 12px; border: 1px solid var(--accent-color); border-radius: 8px; background-color: var(--accent-color); color: var(--accent-text); font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
  .btn-back:hover { background-color: var(--accent-hover); }
  .btn-hide { padding: 8px 12px; border: 1px solid var(--btn-border); border-radius: 8px; background-color: var(--btn-bg); color: var(--btn-text); font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; }
  .btn-hide:hover { background-color: var(--btn-hover-bg); }
  .container { max-width: 800px; margin: 0; background-color: var(--card-bg); border-radius: 0; box-shadow: none; padding-bottom: 20px; color: var(--text-color); }
  .image-container img { width: 100%; height: auto; border-radius: 0; display: block; }
  .main-content { padding: 1.5em; }
  h1 { margin: 0 0 10px; font-size: 2em; color: var(--text-color); }
  .description { font-size: 1.1em; color: var(--secondary-text); margin-bottom: 0.5em; }
  .analysis-section { border-top: 2px solid var(--border-color); padding-top: 1.5em; margin-top: 1.5em; }
  .metric-item { display: block; margin-bottom: 1.5em; }
  .metric-label { font-weight: bold; font-size: 1.1em; margin-bottom: 5px; display: block; color: var(--text-color); }
  .metric-details { display: block; width: 100%; }
  .reason { margin: 5px 0 0; font-size: 0.95em; color: var(--secondary-text); line-height: 1.4; }
  .original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid var(--border-color); }
  .action-buttons { display: flex; justify-content: space-around; gap: 10px; padding: 1.5em; background-color: var(--card-bg); border-top: 1px solid var(--border-color); }
  .btn { flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px; padding: 10px; border: 1px solid var(--btn-border); border-radius: 8px; background-color: var(--btn-bg); color: var(--btn-text); font-size: 0.95em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
  .btn:hover { background-color: var(--btn-hover-bg); }
  .btn-primary { background-color: var(--accent-color); color: var(--accent-text); border-color: var(--accent-color); }
  .btn-primary:hover { background-color: var(--accent-hover); }

  @media (min-width: 800px) {
    :host { padding: 2em 0; }
    .page-header { border-radius: 12px 12px 0 0; }
    .container { margin: 0 auto; border-radius: 12px; box-shadow: var(--card-shadow); }
  }

  @media (max-width: 799px) {
    :host { padding-top: 55px; padding-bottom: 80px; }
    .page-header { position: fixed; top: 0; left: 0; width: 100%; box-sizing: border-box; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
    .action-buttons { position: fixed; bottom: 0; left: 0; width: 100%; z-index: 1000; box-shadow: 0 -2px 10px rgba(0,0,0,0.1); margin-top: 0; box-sizing: border-box; }
  }


`;

interface ArticlePageProps {
  item: AnalyzedItem;
}

export const ArticlePage = ({ item }: ArticlePageProps) => {
  const { t } = useTranslation();
  const [showDetails, setShowDetails] = useState(false);
  const [showOriginal, setShowOriginal] = useState(false);

  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;
  const title = item.title_corrected || item.title_original || t('article.no_title');
  const description = item.description_corrected || item.description_original || t('article.no_description');
  const imageUrl = item.media?.medium === 'image' ? item.media.url : '';

  const metrics = [
    { id: 'framing', label: t('metrics.framing'), value: item.framing, reason: item.framing_reason },
    { id: 'clickbait', label: t('metrics.clickbait'), value: item.clickbait, reason: item.clickbait_reason },
    { id: 'persuasive', label: t('metrics.persuasive'), value: item.persuasive, reason: item.persuasive_reason },
    { id: 'hyper_stimulus', label: t('metrics.hyper_stimulus'), value: item.hyper_stimulus, reason: item.hyper_stimulus_reason },
    { id: 'speculative', label: t('metrics.speculative'), value: item.speculative, reason: item.speculative_reason },
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
          <a href="/" className="btn-back" title={t('article.back_tooltip', { domain: rootDomain })}>{t('article.back')}</a>
          <button onClick={bypassAndReload} className="btn-hide">{t('article.hide')}</button>
        </header>

        {imageUrl && (
          <div className="image-container">
            <img src={imageUrl} alt={title} />
          </div>
        )}

        <div className="main-content">
          <h1>{title}</h1>
          <p className="description">{description}</p>
          <MetaData pubDate={(item as AnalyzedItem & { pubDate?: string | Date }).pubDate} />

          <div className="analysis-section">
            <div className="metric-item">
                <RatingBar value={item.rating} label={t('metrics.overall_rating')} reason={item.overall_reason} />
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
                 <h3>{t('article.original_section')}</h3>
                 <h4>{item.title_original || ''}</h4>
                 <p>{item.description_original || ''}</p>
               </div>
            )}
          </div>

        </div>

        <div className="action-buttons">
          <button onClick={() => setShowOriginal(true)} className="btn btn-primary">{t('article.btn_original_title')}</button>
          <button onClick={() => { setShowDetails(true); setShowOriginal(true); }} className="btn">{t('article.btn_details')}</button>
          <button onClick={bypassAndReload} className="btn">{t('article.btn_view_original')}</button>
        </div>

        <Footer />
      </div>
    </>
  );
};
