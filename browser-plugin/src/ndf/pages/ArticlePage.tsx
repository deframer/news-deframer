import '../../shared/i18n';

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { MetaData } from '../components/MetaData';
import { RatingBar } from '../components/RatingBar';
import { stripHtml } from '../utils/html-utils';

interface ArticlePageProps {
  item: AnalyzedItem;
}

export const ArticlePage = ({ item }: ArticlePageProps) => {
  const { t } = useTranslation();
  const [showDetails, setShowDetails] = useState(false);
  const [showOriginal, setShowOriginal] = useState(false);

  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;
  const title = item.title_corrected || stripHtml(item.title_original) || t('article.no_title');
  const description = item.description_corrected || stripHtml(item.description_original)|| t('article.no_description');
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
    log.debug('Bypassing for this session and reloading.');
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
  };

  return (
    <>
      <div className="article-container">
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
                 <h4>{stripHtml(item.title_original || '')}</h4>
                 <p>{stripHtml(item.description_original || '')}</p>
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
