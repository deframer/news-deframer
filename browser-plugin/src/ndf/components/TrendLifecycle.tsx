import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { TrendLifecycleMetric, TrendRepo } from './TrendRepo';

interface TrendLifecycleProps {
  domain: string;
  days: number;
}

export const TrendLifecycle = ({ domain, days }: TrendLifecycleProps) => {
  const { t } = useTranslation();
  const [term, setTerm] = useState('');
  const [data, setData] = useState<TrendLifecycleMetric[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    if (!term.trim()) return;
    setLoading(true);
    const result = await TrendRepo.getTrendLifecycle(term, domain, days);
    // Sort by date ascending for the chart
    const sorted = [...result].sort((a, b) => new Date(a.time_slice).getTime() - new Date(b.time_slice).getTime());
    setData(sorted);
    setLoading(false);
  };

  const maxFreq = data.length > 0 ? Math.max(...data.map(d => d.frequency)) : 0;

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
        <button className="lifecycle-btn" onClick={handleSearch} disabled={loading}>
          {loading ? '...' : t('trends.analyze', 'Analyze')}
        </button>
      </div>

      {data.length > 0 && (
        <div className="chart-container">
          {data.map((item, idx) => {
            const heightPercent = maxFreq > 0 ? (item.frequency / maxFreq) * 100 : 0;
            const dateLabel = new Date(item.time_slice).toLocaleDateString(undefined, { month: 'short', day: 'numeric' });

            const textShadow = '0 0 3px rgba(0,0,0,0.7)';
            const style: React.CSSProperties = { height: `${heightPercent}%` };
            let icon = null;
            let barClass = 'chart-bar';
            // Unified label style: below the bar with a -45 degree angle
            const labelStyle: React.CSSProperties = {
              color: 'var(--text-color)',
            };

            if (item.velocity > 0) {
              style.backgroundColor = 'var(--success-color, #198754)';
              icon = <span className="trend-icon" style={{ color: 'var(--success-color, #198754)' }}>▲</span>;
              labelStyle.color = '#ffffff';
              labelStyle.textShadow = textShadow;
            } else if (item.velocity < 0) {
              style.backgroundColor = 'var(--danger-color, #b02a37)';
              icon = <span className="trend-icon" style={{ color: 'var(--danger-color, #b02a37)' }}>▼</span>;
              labelStyle.color = '#ffffff';
              labelStyle.textShadow = textShadow;
            } else {
              barClass += ' lateral';
            }

            return (
              <div
                key={item.time_slice}
                className="chart-bar-wrapper"
                tabIndex={0}
                role="img"
                aria-label={t('trends.lifecycle_aria_label', '{{date}}: Frequency {{frequency}}, Velocity {{velocity}}', {
                  date: dateLabel,
                  frequency: item.frequency,
                  velocity: (item.velocity > 0 ? '+' : '') + item.velocity
                })}
              >
                <div
                  className={barClass}
                  style={style}
                >
                  {icon}
                  <div className="bar-tooltip">
                    {dateLabel}<br/>
                    {t('trends.freq', 'Freq')}: {item.frequency}<br/>
                    {t('trends.vel', 'Vel')}: {item.velocity > 0 ? '+' : ''}{item.velocity}
                  </div>
                </div>
                {(data.length < 15 || idx % Math.ceil(data.length / 10) === 0) && (
                  <div className="bar-label" style={labelStyle}>{dateLabel}</div>
                )}
              </div>
            );
          })}
        </div>
      )}

      {data.length === 0 && !loading && term && (
        <div style={{ textAlign: 'center', color: 'var(--secondary-text)', marginTop: '2em' }}>
          {t('trends.lifecycle_no_data', 'No data to display. Try searching for "trump".')}
        </div>
      )}
    </div>
  );
};