import { useTranslation } from 'react-i18next';
import { TrendComparisonMetric } from './TrendRepo';

interface TrendCompareProps {
  items: TrendComparisonMetric[];
  currentDomain: string;
  compareDomain: string;
}

const compareCss = `
  .compare-grid {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 16px;
    margin-top: 10px;
  }
  .compare-col {
    background: var(--card-bg, #fff);
    border: 1px solid var(--border-color, #eee);
    border-radius: 8px;
    padding: 12px;
    display: flex;
    flex-direction: column;
  }
  .col-header {
    font-size: 0.9em;
    font-weight: bold;
    text-align: center;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 2px solid transparent;
    color: var(--text-color);
  }
  .col-header.unique-a { border-color: var(--primary-color, #0056b3); }
  .col-header.intersect { border-color: var(--secondary-text, #6c757d); }
  .col-header.unique-b { border-color: var(--danger-color, #dc3545); }

  .compare-list {
    list-style: none;
    padding: 0;
    margin: 0;
  }
  .compare-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    border-bottom: 1px solid var(--border-color, #f5f5f5);
    font-size: 0.9em;
  }
  .compare-item:last-child { border-bottom: none; }
  .topic-name { font-weight: 500; }
  .topic-score {
    font-size: 0.8em;
    color: var(--secondary-text);
    background: var(--bg-color-secondary, #f5f5f5);
    padding: 2px 6px;
    border-radius: 4px;
  }
`;

export const TrendCompare = ({ items, currentDomain, compareDomain }: TrendCompareProps) => {
  const { t } = useTranslation();

  const uniqueA = items.filter(i => i.classification === 'BLINDSPOT_A');
  const intersect = items.filter(i => i.classification === 'INTERSECT');
  const uniqueB = items.filter(i => i.classification === 'BLINDSPOT_B');

  const renderList = (list: TrendComparisonMetric[], scoreKey: 'score_a' | 'score_b' | 'both') => (
    <ul className="compare-list">
      {list.map((item, idx) => (
        <li key={`${item.trend_topic}-${idx}`} className="compare-item">
          <span className="topic-name">{item.trend_topic}</span>
          <span className="topic-score">
            {scoreKey === 'both'
              ? `${item.score_a} / ${item.score_b}`
              : item[scoreKey]}
          </span>
        </li>
      ))}
      {list.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>None</li>}
    </ul>
  );

  return (
    <>
      <style>{compareCss}</style>
      <div className="compare-grid">
        {/* Column A: Unique to Current Domain */}
        <div className="compare-col">
          <div className="col-header unique-a">
            Our Topics
          </div>
          {renderList(uniqueA, 'score_a')}
        </div>

        {/* Column B: Blindspots (Unique to Compare Domain) */}
        <div className="compare-col">
          <div className="col-header unique-b">
            Their Topics
          </div>
          {renderList(uniqueB, 'score_b')}
        </div>

        {/* Column Intersect: Shared */}
        <div className="compare-col">
          <div className="col-header intersect">
            Shared
          </div>
          {renderList(intersect, 'both')}
        </div>
      </div>
    </>
  );
};