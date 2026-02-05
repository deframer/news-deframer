import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { TrendContextChart } from './TrendContext';
import { TrendLifecycleChart } from './TrendLifecycleChart';

interface TrendDetailsProps {
  term: string;
  domain: string;
  days: number;
  showBorder?: boolean;
}

export const TrendDetails = ({ term, domain, days, showBorder = true }: TrendDetailsProps) => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState<'lifecycle' | 'context'>('context');

  return (
    <div className={`trend-details-container ${!showBorder ? 'no-border' : ''}`}>
      <div className="sub-tabs">
        <button
          className={`sub-tab-btn ${activeTab === 'context' ? 'active' : ''}`}
          onClick={() => setActiveTab('context')}
        >
          {t('trends.context', 'Context')}
        </button>
        <button
          className={`sub-tab-btn ${activeTab === 'lifecycle' ? 'active' : ''}`}
          onClick={() => setActiveTab('lifecycle')}
        >
          {t('trends.lifecycle', 'Lifecycle')}
        </button>
      </div>
      {activeTab === 'context' && <TrendContextChart topic={term} days={days} domain={domain} />}
      {activeTab === 'lifecycle' && <TrendLifecycleChart domain={domain} days={days} term={term} />}
    </div>
  );
};