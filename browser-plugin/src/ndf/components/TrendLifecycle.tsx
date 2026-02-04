import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { TrendLifecycleMetric, TrendRepo } from './TrendRepo';

interface TrendLifecycleProps {
  domain: string;
  timeRange: string;
}

const lifecycleCss = `
  .lifecycle-controls {
    display: flex;
    gap: 10px;
    margin-bottom: 20px;
  }
  .lifecycle-input {
    flex: 1;
    padding: 8px 12px;
    border: 1px solid var(--border-color, #ddd);
    border-radius: 6px;
    font-size: 0.95em;
    background: var(--card-bg, #fff);
    color: var(--text-color);
  }
  .lifecycle-btn {
    padding: 8px 16px;
    background: var(--primary-color, #0056b3);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 500;
    transition: background-color 0.2s;
  }
  .lifecycle-btn:hover {
    background: var(--primary-color-dark, #004494);
  }
  .chart-container {
    height: 300px;
    display: flex;
    align-items: flex-end;
    gap: 4px;
    padding: 20px 0 70px 0;
    border-bottom: 1px solid var(--border-color, #ccc);
    position: relative;
    width: 100%;
  }
  @media (max-width: 799px) {
    .chart-container {
      gap: 1px;
    }
  }
  .chart-bar-wrapper {
    flex: 1;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    position: relative;
    outline: none;
  }
  .chart-bar-wrapper:hover,
  .chart-bar-wrapper:focus-visible {
    z-index: 20;
  }
  .chart-bar-wrapper:focus-visible .chart-bar {
    box-shadow: 0 0 0 2px var(--text-color);
    opacity: 0.8;
  }
  .chart-bar {
    width: 100%;
    border-radius: 4px 4px 0 0;
    transition: height 0.5s ease;
    min-height: 65px;
    position: relative;
  }
  .chart-bar:hover {
    opacity: 0.8;
  }
  .chart-bar.lateral {
    background-color: var(--badge-bg, #e9ecef);
    border-top: 3px dotted var(--secondary-text);
  }
  .trend-icon {
    position: absolute;
    bottom: 100%;
    left: 50%;
    transform: translateX(-50%);
    font-size: 0.8em;
    line-height: 1;
    margin-bottom: 4px;
    font-weight: bold;
  }
  .bar-label {
    font-size: 0.85em;
    font-weight: bold;
    text-align: center;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    position: absolute;
    bottom: 8px;
    left: 0;
    width: 100%;
    pointer-events: none;
    z-index: 2;
  }
  @media (max-width: 799px) {
    .bar-label {
      left: 50%;
      bottom: 25px;
      transform: translateX(-50%) rotate(-90deg);
      transform-origin: center;
      text-align: center;
      width: 80px;
    }
  }
  .bar-tooltip {
    position: absolute;
    bottom: 35px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0,0,0,0.8);
    color: white;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 0.75em;
    white-space: nowrap;
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.2s;
    z-index: 10;
  }
  @media (max-width: 799px) {
    .bar-tooltip {
      bottom: 70px;
    }
  }
  .chart-bar-wrapper:hover .bar-tooltip,
  .chart-bar-wrapper:focus-visible .bar-tooltip {
    opacity: 1;
  }
`;

export const TrendLifecycle = ({ domain, timeRange }: TrendLifecycleProps) => {
  const { t } = useTranslation();
  const [term, setTerm] = useState('');
  const [data, setData] = useState<TrendLifecycleMetric[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = async () => {
    if (!term.trim()) return;
    setLoading(true);
    const result = await TrendRepo.getTrendLifecycle(term, domain, timeRange);
    // Sort by date ascending for the chart
    const sorted = [...result].sort((a, b) => new Date(a.time_slice).getTime() - new Date(b.time_slice).getTime());
    setData(sorted);
    setLoading(false);
  };

  const maxFreq = data.length > 0 ? Math.max(...data.map(d => d.frequency)) : 0;

  return (
    <div>
      <style>{lifecycleCss}</style>
      <div className="lifecycle-controls">
        <input
          type="text"
          className="lifecycle-input"
          placeholder={t('trends.lifecycle_placeholder') || "Enter term to analyze..."}
          value={term}
          onChange={(e) => setTerm(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
        />
        <button className="lifecycle-btn" onClick={handleSearch} disabled={loading}>
          {loading ? '...' : (t('trends.analyze') || 'Analyze')}
        </button>
      </div>

      {data.length > 0 && (
        <div className="chart-container">
          {data.map((item, idx) => {
            const heightPercent = maxFreq > 0 ? (item.frequency / maxFreq) * 100 : 0;
            const dateLabel = new Date(item.time_slice).toLocaleDateString(undefined, { month: 'short', day: 'numeric' });

            const textShadow = '0 0 3px rgba(0,0,0,0.7)';
            let style: React.CSSProperties = { height: `${heightPercent}%` };
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
                aria-label={`${dateLabel}: Frequency ${item.frequency}, Velocity ${item.velocity > 0 ? '+' : ''}${item.velocity}`}
              >
                <div
                  className={barClass}
                  style={style}
                >
                  {icon}
                  <div className="bar-tooltip">
                    {dateLabel}<br/>
                    Freq: {item.frequency}<br/>
                    Vel: {item.velocity > 0 ? '+' : ''}{item.velocity}
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
          No data to display. Try searching for "trump".
        </div>
      )}
    </div>
  );
};