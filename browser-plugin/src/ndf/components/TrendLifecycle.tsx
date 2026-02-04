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
    height: 250px;
    display: flex;
    align-items: flex-end;
    gap: 4px;
    padding: 20px 0 30px 0;
    border-bottom: 1px solid var(--border-color, #ccc);
    position: relative;
  }
  .chart-bar-wrapper {
    flex: 1;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    position: relative;
  }
  .chart-bar {
    width: 100%;
    background-color: var(--primary-color, #0056b3);
    border-radius: 4px 4px 0 0;
    transition: height 0.5s ease;
    min-height: 2px;
    position: relative;
  }
  .chart-bar:hover {
    opacity: 0.8;
  }
  .bar-label {
    font-size: 0.7em;
    color: var(--secondary-text);
    text-align: center;
    margin-top: 5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transform: rotate(-45deg);
    transform-origin: top left;
    position: absolute;
    top: 100%;
    left: 50%;
    width: 100px;
  }
  .bar-tooltip {
    position: absolute;
    bottom: 100%;
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
    margin-bottom: 5px;
  }
  .chart-bar-wrapper:hover .bar-tooltip {
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

            // Color coding for velocity: Green for positive, Red for negative, Default for 0
            let barColor = 'var(--primary-color, #0056b3)';
            if (item.velocity > 0) barColor = '#28a745'; // Green
            if (item.velocity < 0) barColor = '#dc3545'; // Red

            return (
              <div key={item.time_slice} className="chart-bar-wrapper">
                <div
                  className="chart-bar"
                  style={{ height: `${heightPercent}%`, backgroundColor: barColor }}
                ></div>
                <div className="bar-tooltip">
                  {dateLabel}<br/>
                  Freq: {item.frequency}<br/>
                  Vel: {item.velocity > 0 ? '+' : ''}{item.velocity}
                </div>
                {/* Show label only for every nth item if too many items, or always if few */}
                {(data.length < 15 || idx % Math.ceil(data.length / 10) === 0) && (
                  <div className="bar-label">{dateLabel}</div>
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