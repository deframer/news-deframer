import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';
import { TrendDetails } from './TrendDetails';

interface TrendSearchProps {
  domain: DomainEntry;
  days: number;
}

export const TrendSearch = ({ domain, days }: TrendSearchProps) => {
  const { t } = useTranslation();
  const [term, setTerm] = useState('');
  const [activeTerm, setActiveTerm] = useState('');

  const handleSearch = () => {
    if (!term.trim()) return;
    setActiveTerm(term);
  };

  return (
    <div>
      <div className="lifecycle-controls">
        <input
          type="text"
          className="lifecycle-input"
          placeholder={t('trends.search_placeholder', "Enter term to analyze...")}
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button className="lifecycle-btn" onClick={handleSearch}>
          {t('trends.analyze', 'Search')}
        </button>
      </div>

      {activeTerm && <TrendDetails domain={domain} days={days} term={activeTerm} showBorder={false} />}
    </div>
  );
};