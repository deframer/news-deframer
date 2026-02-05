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
  domain: string;
}

export const TrendCompare = ({ items, baseItems, compareDomain, availableDomains, onSelectDomain, domain }: TrendCompareProps) => {
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
      <div className="compare-grid">
        {/* Column A: Unique to Current Domain */}
        <div className="compare-col">
          <div className="col-header unique-a">
            {t('trends.compare.trending_on', { domain })}
          </div>
          {compareDomain ? renderList(uniqueA, 'score_a') : renderBaseList(baseItems)}
        </div>

        {/* Column B: Blindspots (Unique to Compare Domain) */}
        <div className="compare-col">
          <div className="col-header unique-b">
            <span style={{ marginRight: '5px' }}>{t('trends.compare.their_topics', 'Trending')}</span>
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