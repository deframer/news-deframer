import { CSSProperties, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { DomainEntry, Lifecycle, NewsDeframerClient } from '../client';
import { ArticleList } from './ArticleList';

interface TrendLifecycleChartProps {
  domain: DomainEntry;
  days: number;
  term: string;
}

export const TrendLifecycleChart = ({ domain, days, term }: TrendLifecycleChartProps) => {
  const { t, i18n } = useTranslation();
  const [data, setData] = useState<Lifecycle[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchData = async () => {
      if (!term) return;
      setLoading(true);
      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const result = await client.getLifecycleByDomain(term, domain.domain, domain.language, days);
        // Sort by date ascending for the chart
        const sorted = [...result].sort((a, b) => new Date(a.time_slice).getTime() - new Date(b.time_slice).getTime());
        setData(sorted);
      } catch (error) {
        console.error('Failed to fetch trend lifecycle data', error);
        setData([]);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [term, domain, days]);


  useEffect(() => {
    setSelectedDate(null);
  }, [term, domain.domain]);

  useEffect(() => {
    if (!selectedDate) return;
    if (typeof window === 'undefined') return;
    if (window.matchMedia('(max-width: 799px)').matches) return;
    const trendContent = chartRef.current?.closest('.trend-content');
    if (trendContent instanceof HTMLElement && trendContent.scrollHeight > trendContent.clientHeight + 1) {
      trendContent.scrollTo({ top: trendContent.scrollHeight, behavior: 'smooth' });
      return;
    }
    window.scrollTo({ top: document.documentElement.scrollHeight, behavior: 'smooth' });
  }, [selectedDate]);

  const maxFreq = data.length > 0 ? Math.max(...data.map(d => d.frequency)) : 0;

  if (loading) {
    return (
      <div className="chart-container" style={{ alignItems: 'center', justifyContent: 'center' }}>
        <div className="spinner-small" />
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="chart-container" style={{ alignItems: 'center', justifyContent: 'center', color: 'var(--secondary-text)' }}>
        {t('trends.lifecycle_no_data', 'No lifecycle data available for this topic.')}
      </div>
    );
  }

  return (
    <>
      <div ref={chartRef} className="chart-container">
        {data.map((item, idx) => {
          const heightPercent = maxFreq > 0 ? (item.frequency / maxFreq) * 100 : 0;
          const dateLabel = new Date(item.time_slice).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' });

          const isSelected = selectedDate === item.time_slice;
          const style: CSSProperties = { height: `${heightPercent}%` };
          let icon = null;
          let barClass = 'chart-bar';
          // Unified label style: below the bar with a -45 degree angle
          const labelStyle: CSSProperties = {
            color: 'var(--text-color)',
          };
  
          if (item.velocity > 0) {
            style.backgroundColor = 'var(--trend-up)';
            icon = <span className="trend-icon" style={{ color: 'var(--trend-up)' }}>▲</span>;
            labelStyle.color = 'var(--trend-text)';
          } else if (item.velocity < 0) {
            style.backgroundColor = 'var(--trend-down)';
            icon = <span className="trend-icon" style={{ color: 'var(--trend-down)' }}>▼</span>;
            labelStyle.color = 'var(--trend-text)';
          } else {
            style.backgroundColor = 'var(--trend-steady)';
            icon = <span className="trend-icon" style={{ color: 'var(--trend-steady)' }}>▶</span>;
            barClass += ' lateral';
            labelStyle.color = 'var(--trend-text)';
          }

          return (
            <div
              key={item.time_slice}
              className={`chart-bar-wrapper ${isSelected ? 'selected' : ''}`}
              onClick={() => setSelectedDate(isSelected ? null : item.time_slice)}
              tabIndex={0}
              role="img"
              aria-label={t('trends.search_aria_label', '{{date}}: Frequency {{frequency}}, Velocity {{velocity}}', {
                date: dateLabel,
                frequency: item.frequency,
                velocity: (item.velocity > 0 ? '+' : '') + item.velocity
              })}
            >
              <div className={barClass} style={style}>
                {icon}
                <div className="bar-tooltip">
                  {dateLabel}<br/>
                  {t('trends.freq', 'Freq')}: {item.frequency}<br/>
                  {t('trends.vel', 'Vel')}: {(item.velocity > 0 ? '+' : '') + item.velocity}
                </div>
              </div>
              {(data.length < 15 || idx % Math.ceil(data.length / 10) === 0) && (
                <div className="bar-label" style={labelStyle}>{dateLabel}</div>
              )}
            </div>
          );
        })}
      </div>
      {selectedDate && (
        <ArticleList
          term={term}
          domain={domain}
          date={new Date(selectedDate).toISOString().split('T')[0]} /* Explicitly format to YYYY-MM-DD for API */
          days={undefined}
          titleOverride={`${t('trends.articles', 'Articles')} / ${new Date(selectedDate).toLocaleDateString(i18n.language, { month: 'short', day: 'numeric' })} / ${term}`}
        />
      )}
    </>
  );
};
