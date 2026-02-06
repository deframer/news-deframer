import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';
import { TrendContextMetric, TrendRepo } from './TrendRepo';

interface TrendContextProps {
  topic: string;
  className?: string;
  days?: number;
  domain?: DomainEntry;
}

export const TrendContextChart = ({ topic, className, days, domain }: TrendContextProps) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<TrendContextMetric[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    setLoading(true);
    TrendRepo.getTrendContext(topic, domain?.domain, days).then((data) => {
      if (mounted) {
        setItems(data);
        setLoading(false);
      }
    });
    return () => { mounted = false; };
  }, [topic, days, domain]);

  if (loading) return (
    <div className={`trend-context ${className || ''}`} style={{ minHeight: '80px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="spinner-small" />
    </div>
  );
  if (items.length === 0) return null;

  return (
    <div className={`trend-context ${className || ''}`}>
      <div className="context-header">
        {t('trends.context_header', 'Context: How is "{{topic}}" being described?', { topic })}
      </div>
      <div className="context-list">
        {items.map((item) => (
          <span key={item.context_word} className="context-chip">
            {item.context_word}
            <div className="chip-tooltip">
              {item.type} - {t('trends.frequency_label', 'Frequency')}: {item.frequency}
            </div>
          </span>
        ))}
      </div>
    </div>
  );
};