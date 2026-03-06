import { Dispatch, SetStateAction, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';
import { ArticleList } from './ArticleList';
import { TrendContextChart } from './TrendContext';
import { TrendLifecycleChart } from './TrendLifecycleChart';

interface TrendDetailsProps {
  term: string;
  domain: DomainEntry;
  days: number;
  showBorder?: boolean;
  activeTab: 'lifecycle' | 'context' | 'articles';
  setActiveTab: Dispatch<SetStateAction<'lifecycle' | 'context' | 'articles'>>;
}

export const TrendDetails = ({ term, domain, days, showBorder = true, activeTab, setActiveTab }: TrendDetailsProps) => {
  const { t } = useTranslation();
  const containerRef = useRef<HTMLDivElement>(null);

  const scrollTrendViewToBottom = () => {
    if (typeof window === 'undefined') return;
    const trendContent = containerRef.current?.closest('.trend-content');
    if (trendContent instanceof HTMLElement && trendContent.scrollHeight > trendContent.clientHeight + 1) {
      trendContent.scrollTo({ top: trendContent.scrollHeight, behavior: 'smooth' });
      return;
    }
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'smooth' });
  };

  useEffect(() => {
    if (activeTab !== 'articles') return;
    if (typeof window === 'undefined') return;
    if (window.matchMedia('(max-width: 799px)').matches) return;
    scrollTrendViewToBottom();
  }, [activeTab, term, days]);

  return (
    <div ref={containerRef} className={`trend-details-container ${!showBorder ? 'no-border' : ''}`}>
      <div className="sub-tabs">
        <button
          className={`sub-tab-btn ${activeTab === 'lifecycle' ? 'active' : ''}`}
          onClick={() => setActiveTab('lifecycle')}
        >
          {t('trends.lifecycle', 'Lifecycle')}
        </button>
        <button
          className={`sub-tab-btn ${activeTab === 'context' ? 'active' : ''}`}
          onClick={() => setActiveTab('context')}
        >
          {t('trends.context', 'Context')}
        </button>
        <button
          className={`sub-tab-btn ${activeTab === 'articles' ? 'active' : ''}`}
          onClick={() => setActiveTab('articles')}
        >
          {t('trends.articles', 'Articles')}
        </button>
      </div>
      {activeTab === 'context' && <TrendContextChart topic={term} days={days} domain={domain} />}
      {activeTab === 'lifecycle' && <TrendLifecycleChart domain={domain} days={days} term={term} />}
      {activeTab === 'articles' && <ArticleList term={term} domain={domain} days={days} hideTitle={true} />}
    </div>
  );
};
