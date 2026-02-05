import { useTranslation } from 'react-i18next';

import { TrendItem } from './TabTrend';
import { TrendComparisonMetric } from './TrendRepo';

interface DomainOption {
  id: string;
  name: string;
}

interface TrendCompareProps {
  items: TrendComparisonMetric[];
  baseItems: TrendItem[];
  compareDomain: string | null;
  availableDomains: DomainOption[];
  onSelectDomain: (domain: string) => void;
}

const compareCss = `
  .compare-grid {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 16px;
    margin-top: 10px;
  }
  @media (max-width: 799px) {
    .compare-grid {
      grid-template-columns: 1fr;
    }
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
  .col-header select {
    font-size: 1em;
    font-weight: bold;
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
  .header-select {
    padding: 4px 8px;
    border-radius: 4px;
    border: 1px solid var(--border-color, #ddd);
    background: var(--card-bg, #fff);
    color: var(--text-color);
    font-size: 0.9em;
    font-weight: bold;
    max-width: 100%;
  }
`;

export const TrendCompare = ({ items, baseItems, compareDomain, availableDomains, onSelectDomain }: TrendCompareProps) => {
  const { t } = useTranslation();

  // If comparing, use the comparison data. If not, uniqueA is just the base list (top 10)
  const uniqueA = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_A') : [];
  const intersect = compareDomain ? items.filter(i => i.classification === 'INTERSECT') : [];
  const uniqueB = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_B') : [];

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
      {list.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>{t('trends.none', 'None')}</li>}
    </ul>
  );

  const renderBaseList = (list: TrendItem[]) => (
    <ul className="compare-list">
      {list.slice(0, 10).map((item) => (
        <li key={item.word} className="compare-item">
          <span className="topic-name">{item.word}</span>
          <span className="topic-score">
            {item.outlierRatio.toFixed(1)}x
          </span>
        </li>
      ))}
    </ul>
  );

  return (
    <>
      <style>{compareCss}</style>
      <div className="compare-grid">
        {/* Column A: Unique to Current Domain */}
        <div className="compare-col">
          <div className="col-header unique-a">
            {t('trends.compare.our_topics', 'Our Topics')}
          </div>
          {compareDomain ? renderList(uniqueA, 'score_a') : renderBaseList(baseItems)}
        </div>

        {/* Column B: Blindspots (Unique to Compare Domain) */}
        <div className="compare-col">
          <div className="col-header unique-b">
            <span style={{ marginRight: '5px' }}>{t('trends.compare.their_topics', 'Their Topics')}</span>
            <select
              className="header-select"
              value={compareDomain || ""}
              onChange={(e) => onSelectDomain(e.target.value)}
            >
              {availableDomains.map(d => (
                <option key={d.id} value={d.id}>{d.name}</option>
              ))}
            </select>
          </div>
          {renderList(uniqueB, 'score_b')}
        </div>

        {/* Column Intersect: Shared */}
        <div className="compare-col">
          <div className="col-header intersect">
            {t('trends.compare.shared', 'Shared')}
          </div>
          {renderList(intersect, 'both')}
        </div>
      </div>
    </>
  );
};