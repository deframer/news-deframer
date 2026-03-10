import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, TrendContext } from '../client';

interface TrendContextProps {
  topic: string;
  className?: string;
  days?: number;
  domain?: DomainEntry;
}

export const TrendContextChart = ({ topic, className, days, domain }: TrendContextProps) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<TrendContext[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    const fetchData = async () => {
      if (!domain) return;
      setLoading(true);
      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const data = await client.getContextByDomain(topic, domain.domain, domain.language, days || 7);
        if (mounted) {
          setItems(data);
        }
      } finally {
        if (mounted) setLoading(false);
      }
    };
    fetchData();
    return () => { mounted = false; };
  }, [topic, days, domain]);

  if (loading) return (
    <div className={`trend-context ${className || ''}`} style={{ minHeight: '80px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="spinner-small" />
    </div>
  );
  if (items.length === 0) return (
    <div className={`trend-context ${className || ''}`} style={{ minHeight: '80px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--secondary-text)' }}>
      {t('trends.context_no_data', 'No context data available for this topic.')}
    </div>
  );

  return (
    <div className={`trend-context ${className || ''}`}>
      <div className="context-header">
        {t('trends.context_header', 'Context: How is "{{topic}}" being described?', { topic })}
      </div>
      <div className="context-list">
        {items.map((item) => (
          <span key={item.context} className="context-chip">
            {item.context}
            <div className="chip-tooltip">
              {t('trends.frequency_label', 'Frequency')}: {item.frequency}
            </div>
          </span>
        ))}
      </div>
    </div>
  );
};