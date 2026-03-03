import { ReactNode, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { DomainComparison, DomainEntry, NewsDeframerClient, TrendMetric } from '../client';

interface DomainOption {
  id: string;
  name: string;
}

interface TrendCompareProps {
  days: number;
  baseItems: TrendMetric[];
  compareDomain: string | null;
  availableDomains: DomainOption[];
  onSelectDomain: (domain: string) => void;
  domain: DomainEntry;
  searchEngineUrl: string;
}

const BULLET_DELIMITER = '•';

const OpenIcon = ({ onClick, tooltip }: { onClick: () => void; tooltip: ReactNode }) => (
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
    <div className="icon-tooltip">{tooltip}</div>
  </div>
);

export const TrendCompare = ({ baseItems, compareDomain, availableDomains, onSelectDomain, domain, searchEngineUrl, days }: TrendCompareProps) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<DomainComparison[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      if (!compareDomain) return;
      setLoading(true);
      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const data = await client.getDomainComparison(domain.domain, compareDomain, domain.language, days);
        setItems(data);
      } catch (error) {
        console.error('Failed to fetch trend comparison', error);
        setItems([]);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [domain, compareDomain, days]);

  const handleSearch = (term: string, searchDomain: string) => {
    const query = encodeURIComponent(`${term} site:${searchDomain}`);
    const baseUrl = searchEngineUrl.replace(/\/$/, '');
    window.open(`${baseUrl}/search?q=${query}`, '_blank');
  };

  // If comparing, use the comparison data. If not, uniqueA is just the base list (top 10)
  const uniqueA = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_A') : [];
  const intersect = compareDomain ? items.filter(i => i.classification === 'INTERSECT') : [];
  const uniqueB = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_B') : [];

  const renderList = (list: DomainComparison[], scoreKey: 'score_a' | 'score_b' | 'both') => (
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
            {scoreKey === 'score_a' && <OpenIcon onClick={() => handleSearch(item.trend_topic, domain.domain)} tooltip={<>{item.trend_topic} <span style={{ opacity: 0.6, margin: '0 2px' }}>{BULLET_DELIMITER}</span> {domain.domain}</>} />}
            {scoreKey === 'score_b' && compareDomain && <OpenIcon onClick={() => handleSearch(item.trend_topic, compareDomain)} tooltip={<>{item.trend_topic} <span style={{ opacity: 0.6, margin: '0 2px' }}>{BULLET_DELIMITER}</span> {compareDomain}</>} />}
            {scoreKey === 'both' && (
              <>
                <OpenIcon onClick={() => handleSearch(item.trend_topic, domain.domain)} tooltip={<>{item.trend_topic} <span style={{ opacity: 0.6, margin: '0 2px' }}>{BULLET_DELIMITER}</span> {domain.domain}</>} />
                {compareDomain && <OpenIcon onClick={() => handleSearch(item.trend_topic, compareDomain)} tooltip={<>{item.trend_topic} <span style={{ opacity: 0.6, margin: '0 2px' }}>{BULLET_DELIMITER}</span> {compareDomain}</>} />}
              </>
            )}
          </div>
        </li>
      ))}
      {list.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>{t('trends.none', 'None')}</li>}
    </ul>
  );

  const renderBaseList = (list: TrendMetric[]) => (
    <ul className="compare-list">
      {list.slice(0, 10).map((item) => (
        <li key={item.trend_topic} className="compare-item">
          <span className="topic-name">{item.trend_topic}</span>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <span className="topic-score">
              {item.outlier_ratio.toFixed(1)}x
            </span>
            <OpenIcon onClick={() => handleSearch(item.trend_topic, domain.domain)} tooltip={<>{item.trend_topic} <span style={{ opacity: 0.6, margin: '0 2px' }}>{BULLET_DELIMITER}</span> {domain.domain}</>} />
          </div>
        </li>
      ))}
      {list.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>{t('trends.none', 'None')}</li>}
    </ul>
  );

  if (loading) {
    return <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="spinner-small" />
    </div>;
  }

  return (
    <>
      <div className="compare-grid">
        {/* Column A: Unique to Current Domain */}
        <div className="compare-col">
          <div
            className="col-header unique-a"
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              height: '40px',
              paddingBottom: 0
            }}
          >
            <span>{t('trends.compare.trending_on', { domain: domain.domain })}</span>
          </div>
          {compareDomain ? renderList(uniqueA, 'score_a') : renderBaseList(baseItems)}
        </div>

        {/* Column B: Blindspots (Unique to Compare Domain) */}
        <div className="compare-col">
          <div
            className="col-header unique-b"
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              gap: '8px',
              height: '40px',
              paddingBottom: 0
            }}
          >
            <span>{t('trends.compare.their_topics', 'Trending')}</span>
            <select
              className="header-select"
              value={compareDomain || ""}
              onChange={(e) => onSelectDomain(e.target.value)}
            >
              {availableDomains.map(d => (
                <option key={d.id} value={d.id}>{d.name}</option>
              ))}
            </select>
            {compareDomain && (
              <OpenIcon onClick={() => window.open(`https://${compareDomain}`, '_blank')} tooltip={compareDomain} />
            )}
          </div>
          {renderList(uniqueB, 'score_b')}
        </div>

        {/* Column Intersect: Shared */}
        <div className="compare-col">
          <div
            className="col-header intersect"
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              height: '40px',
              paddingBottom: 0
            }}
          >
            <span>{t('trends.compare.shared', 'Shared')}</span>
          </div>
          {renderList(intersect, 'both')}
        </div>
      </div>
    </>
  );
};