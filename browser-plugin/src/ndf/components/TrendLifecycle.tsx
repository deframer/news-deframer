import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { TrendLifecycleChart } from './TrendLifecycleChart';

interface TrendLifecycleProps {
  domain: string;
  days: number;
}

export const TrendLifecycle = ({ domain, days }: TrendLifecycleProps) => {
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
          placeholder={t('trends.lifecycle_placeholder', "Enter term to analyze...")}
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button className="lifecycle-btn" onClick={handleSearch}>
          {t('trends.analyze', 'Analyze')}
        </button>
      </div>

      {activeTerm && <TrendLifecycleChart domain={domain} days={days} term={activeTerm} />}
    </div>
  );
};