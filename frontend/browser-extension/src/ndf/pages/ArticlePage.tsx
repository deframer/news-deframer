import '../../shared/i18n';

import { useId, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import type { AnalysisOutput, EmotionVector } from '../../shared/sentiments';
import { sentimentsToLabels } from '../../shared/sentiments';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { MetaData } from '../components/MetaData';
import { RatingBar } from '../components/RatingBar';
import { Sentiments } from '../components/Sentiments';
import { SentimentsInterpretation } from '../components/SentimentsInterpretation';
import { SentimentsToggle } from '../components/SentimentsToggle';
import { stripHtml } from '../utils/html-utils';

interface ArticlePageProps {
  item: AnalyzedItem;
}

export const ArticlePage = ({ item }: ArticlePageProps) => {
  const { t, i18n } = useTranslation();
  const tooltipId = useId();
  const [activeTab, setActiveTab] = useState<'article' | 'sentiments'>('article');
  type DetailMode = 'closed' | 'original' | 'details';
  const [mode, setMode] = useState<DetailMode>('closed');

  const toggleMode = (next: Exclude<DetailMode, 'closed'>) => {
    setMode((current) => current === next ? 'closed' : next);
  };
  const [sentimentType, setSentimentType] = useState<'sentiments' | 'sentiments_deframed'>('sentiments');
  const [tooltipContent, setTooltipContent] = useState('');
  const [isTooltipOpen, setIsTooltipOpen] = useState(false);
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 });

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

  const metricValues = useMemo(() => {
    return item[sentimentType] || { valence: 5, arousal: 5, dominance: 5, joy: 1, anger: 1, sadness: 1, fear: 1, disgust: 1 };
  }, [item, sentimentType]);

  const interpretation = useMemo((): AnalysisOutput | null => {
    const ev: EmotionVector = {
      valence: metricValues.valence,
      arousal: metricValues.arousal,
      dominance: metricValues.dominance,
      joy: metricValues.joy,
      anger: metricValues.anger,
      sadness: metricValues.sadness,
      fear: metricValues.fear,
      disgust: metricValues.disgust,
    };
    const lang = i18n.language || 'en';
    return sentimentsToLabels(ev, lang) as AnalysisOutput;
  }, [metricValues, i18n.language]);

  const showTooltip = (target: HTMLElement, description: string) => {
    const rect = target.getBoundingClientRect();
    const preferredX = rect.left + rect.width / 2;
    const tooltipWidth = Math.min(280, window.innerWidth * 0.6, window.innerWidth - 24);
    const minX = tooltipWidth / 2 + 12;
    const maxX = window.innerWidth - tooltipWidth / 2 - 12;
    const clampedX = Math.max(minX, Math.min(preferredX, maxX));
    setTooltipPos({ x: clampedX, y: rect.top - 10 });
    setTooltipContent(description);
    setIsTooltipOpen(true);
  };

  const hideTooltip = () => {
    setIsTooltipOpen(false);
  };

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
          <div style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
            <MetaData
              pubDate={(item as AnalyzedItem & { pubDate?: string | Date }).pubDate}
              author={item.authors?.join(', ')}
            />
          </div>

           <div className="tabs">
             <button
               className={`tab-btn ${activeTab === 'article' ? 'active' : ''}`}
               onClick={() => setActiveTab('article')}
             >
               {t('article.tab_article')}
             </button>
             <button
               className={`tab-btn ${activeTab === 'sentiments' ? 'active' : ''}`}
               onClick={() => setActiveTab('sentiments')}
             >
               {t('article.tab_sentiments')}
             </button>
           </div>

           {activeTab === 'article' ? (
             <div className="analysis-section">
               <div className="metric-item">
                   <RatingBar value={item.rating} label={t('metrics.overall_rating')} reason={item.overall_reason} />
               </div>

                {mode === 'details' && (
                  <div id="metrics-content">
                    {metrics.map(m => (
                      <div className="metric-item" key={m.id}>
                        <RatingBar value={m.value} label={m.label} reason={m.reason} id={m.id} />
                      </div>
                    ))}
                  </div>
                )}

                {(mode === 'original' || mode === 'details') && (
                   <div id="original-content" className="original-content">
                     <h3>{t('article.original_section')}</h3>
                     <h4>{stripHtml(item.title_original || '')}</h4>
                     <p>{stripHtml(item.description_original || '')}</p>
                   </div>
                )}
             </div>
            ) : (
              <div className="sentiment-panel sentiment-layout-vertical">
                <SentimentsToggle
                  sentimentType={sentimentType}
                  onTypeChange={setSentimentType}
                  hasDeframed={!!item.sentiments_deframed}
                />
                <section className="sentiment-columns-layout">
                  <Sentiments
                    metricValues={metricValues}
                    tooltipId={tooltipId}
                    onShowTooltip={showTooltip}
                    onHideTooltip={hideTooltip}
                  />
                  <div className="sentiment-column-right">
                    {interpretation && (
                      <SentimentsInterpretation interpretation={interpretation} />
                    )}
                  </div>
                </section>
                {isTooltipOpen && (
                  <div
                    id={tooltipId}
                    role="tooltip"
                    className="article-header-floating-tooltip"
                    style={{ left: tooltipPos.x, top: tooltipPos.y }}
                  >
                    {tooltipContent}
                  </div>
                )}
              </div>
            )}

        </div>

        <div className="action-buttons">
          <button onClick={() => toggleMode('original')} className={`btn ${mode === 'original' ? 'btn-primary' : ''}`}>{t('article.btn_original_title')}</button>
          <button onClick={() => toggleMode('details')} className={`btn ${mode === 'details' ? 'btn-primary' : ''}`}>{t('article.btn_details')}</button>
          <button onClick={bypassAndReload} className="btn">{t('article.btn_view_original')}</button>
        </div>

        <Footer />
      </div>
    </>
  );
};
