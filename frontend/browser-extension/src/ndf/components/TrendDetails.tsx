import { Dispatch, SetStateAction } from 'react';
import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';
import { ArticleList } from './ArticleList';
import { TrendContextChart } from './TrendContext';
import { TrendLifecycleChart } from './TrendLifecycleChart';
import { TrendSentiments } from './TrendSentiments';

interface TrendDetailsProps {
  term: string;
  domain: DomainEntry;
  days: number;
  showBorder?: boolean;
  activeTab: 'lifecycle' | 'context' | 'articles' | 'sentiments';
  setActiveTab: Dispatch<SetStateAction<'lifecycle' | 'context' | 'articles' | 'sentiments'>>;
}

export const TrendDetails = ({ term, domain, days, showBorder = true, activeTab, setActiveTab }: TrendDetailsProps) => {
  const { t } = useTranslation();

  return (
    <div className={`trend-details-container ${!showBorder ? 'no-border' : ''}`}>
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
        <button
          className={`sub-tab-btn ${activeTab === 'sentiments' ? 'active' : ''}`}
          onClick={() => setActiveTab('sentiments')}
        >
          {t('trends.sentiments', 'Sentiments')}
        </button>
      </div>
      {activeTab === 'context' && <TrendContextChart topic={term} days={days} domain={domain} />}
      {activeTab === 'lifecycle' && <TrendLifecycleChart key={`${term}-${domain.domain}-${days}`} domain={domain} days={days} term={term} />}
      {activeTab === 'articles' && <ArticleList key={`${term}-${domain.domain}-${days}`} term={term} domain={domain} days={days} hideTitle={true} />}
      {activeTab === 'sentiments' && <TrendSentiments term={term} domain={domain} days={days} />}
    </div>
  );
};
