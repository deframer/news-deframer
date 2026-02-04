import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { TrendContext } from './TrendContext';

export interface TrendItem {
  word: string;
  rank: number;
  count: number;
  utility: number;
  outlierRatio: number;
}

interface TrendTopProps {
  items: TrendItem[];
}

const trendTopCss = `
  .trend-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  .trend-item {
    display: flex;
    align-items: center;
    padding: 12px;
    border-bottom: 1px solid var(--border-color, #eee);
    transition: background-color 0.2s;
    cursor: pointer;
  }
  .trend-item:hover {
    background-color: var(--hover-bg, rgba(0,0,0,0.02));
  }
  .trend-rank {
    font-weight: bold;
    margin-right: 15px;
    width: 24px;
    text-align: center;
    color: var(--secondary-text);
    font-size: 0.9em;
  }
  .trend-word {
    flex: 1;
    font-size: 1em;
    color: var(--text-color);
    font-weight: 500;
  }
  .trend-stats {
    display: flex;
    gap: 15px;
    align-items: center;
  }
  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    min-width: 45px;
    position: relative;
  }
  .stat-tooltip {
    visibility: hidden;
    opacity: 0;
    background-color: var(--tooltip-bg, rgba(0,0,0,0.8));
    color: var(--tooltip-text, #fff);
    text-align: center;
    border-radius: 6px;
    padding: 6px 10px;
    position: absolute;
    z-index: 20;
    bottom: 100%;
    left: 50%;
    transform: translateX(-50%);
    margin-bottom: 6px;
    font-size: 12px;
    white-space: nowrap;
    pointer-events: none;
    transition: opacity 0.2s;
    font-weight: normal;
    text-transform: none;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
  }
  .stat-item:hover .stat-tooltip {
    visibility: visible;
    opacity: 1;
  }
  .stat-label {
    font-size: 0.65em;
    text-transform: uppercase;
    color: var(--secondary-text);
  }
  .stat-value {
    font-size: 0.85em;
    font-weight: 600;
    color: var(--text-color);
  }
  .trend-item-container {
    display: flex;
    flex-direction: column;
    border-bottom: 1px solid var(--border-color, #eee);
  }
  .trend-item-container .trend-item {
    border-bottom: none;
  }
  .expand-icon {
    margin-left: 12px;
    color: var(--secondary-text);
    transition: transform 0.2s ease;
    display: flex;
    align-items: center;
    opacity: 0.5;
  }
  .trend-item:hover .expand-icon {
    opacity: 1;
  }
  .expand-icon.expanded {
    transform: rotate(180deg);
  }
`;

export const TrendTop = ({ items }: TrendTopProps) => {
  // Ensure we only show top 10
  const sortedItems = [...items].sort((a, b) => a.rank - b.rank).slice(0, 10);
  const [expandedTopic, setExpandedTopic] = useState<string | null>(null);

  const toggleExpand = (word: string) => {
    setExpandedTopic(expandedTopic === word ? null : word);
  };

  return (
    <>
      <style>{trendTopCss}</style>
      <ul className="trend-list">
        {sortedItems.map((item) => (
          <li key={item.word} className="trend-item-container">
            <div className="trend-item" onClick={() => toggleExpand(item.word)}>
              <span className="trend-rank">{item.rank}</span>
              <span className="trend-word">{item.word}</span>
              <div className="trend-stats">
                <div className="stat-item">
                  <span className="stat-label">Trend</span>
                  <span className="stat-value" style={{ color: item.outlierRatio > 1.5 ? 'var(--primary-color, #0056b3)' : 'inherit' }}>
                    {item.outlierRatio.toFixed(1)}x
                  </span>
                  <div className="stat-tooltip">Burstiness (Outlier Ratio)</div>
                </div>
                <div className="stat-item">
                  <span className="stat-label">Vol</span>
                  <span className="stat-value">{item.count}</span>
                  <div className="stat-tooltip">Frequency (Interest)</div>
                </div>
              </div>
              <div className={`expand-icon ${expandedTopic === item.word ? 'expanded' : ''}`}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </div>
            </div>
            {expandedTopic === item.word && (
              <TrendContext topic={item.word} />
            )}
          </li>
        ))}
      </ul>
    </>
  );
};