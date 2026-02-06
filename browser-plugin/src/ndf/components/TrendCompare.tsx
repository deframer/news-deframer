import { useTranslation } from 'react-i18next';

import { DomainEntry } from '../client';
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
  domain: DomainEntry;
  searchEngineUrl: string;
}

const OpenIcon = ({ onClick, term, domain }: { onClick: () => void; term: string; domain: string }) => (
  <div className="open-icon-wrapper">
    <button
      onClick={(e) => {
        e.stopPropagation();
        onClick();
      }}
      className="open-icon-btn"
      type="button"
    >
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
        <polyline points="15 3 21 3 21 9"></polyline>
        <line x1="10" y1="14" x2="21" y2="3"></line>
      </svg>
    </button>
    <div className="icon-tooltip">
      {term} <span style={{ opacity: 0.6, margin: '0 2px' }}>•</span> {domain}
    </div>
  </div>
);

export const TrendCompare = ({ items, baseItems, compareDomain, availableDomains, onSelectDomain, domain, searchEngineUrl }: TrendCompareProps) => {
  const { t } = useTranslation();

  const handleSearch = (term: string, searchDomain: string) => {
    const query = encodeURIComponent(`${term} site:${searchDomain}`);
    const baseUrl = searchEngineUrl.replace(/\/$/, '');
    window.open(`${baseUrl}/search?q=${query}`, '_blank');
  };

  // If comparing, use the comparison data. If not, uniqueA is just the base list (top 10)
  const uniqueA = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_A') : [];
  const intersect = compareDomain ? items.filter(i => i.classification === 'INTERSECT') : [];
  const uniqueB = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_B') : [];

  const renderList = (list: TrendComparisonMetric[], scoreKey: 'score_a' | 'score_b' | 'both') => (
    <ul className="compare-list">
      {list.map((item, idx) => (
        <li key={`${item.trend_topic}-${idx}`} className="compare-item">
          <span className="topic-name">{item.trend_topic}</span>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <span className="topic-score">
              {scoreKey === 'both'
                ? `${item.score_a} / ${item.score_b}`
                : item[scoreKey]}
            </span>
            {scoreKey === 'score_a' && <OpenIcon onClick={() => handleSearch(item.trend_topic, domain.domain)} term={item.trend_topic} domain={domain.domain} />}
            {scoreKey === 'score_b' && compareDomain && <OpenIcon onClick={() => handleSearch(item.trend_topic, compareDomain)} term={item.trend_topic} domain={compareDomain} />}
            {scoreKey === 'both' && (
              <>
                <OpenIcon onClick={() => handleSearch(item.trend_topic, domain.domain)} term={item.trend_topic} domain={domain.domain} />
                {compareDomain && <OpenIcon onClick={() => handleSearch(item.trend_topic, compareDomain)} term={item.trend_topic} domain={compareDomain} />}
              </>
            )}
          </div>
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
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <span className="topic-score">
              {item.outlierRatio.toFixed(1)}x
            </span>
            <OpenIcon onClick={() => handleSearch(item.word, domain.domain)} term={item.word} domain={domain.domain} />
          </div>
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
            {t('trends.compare.trending_on', { domain: domain.domain })}
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