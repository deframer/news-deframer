import { Fragment, ReactNode, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { DomainComparison, DomainEntry, NewsDeframerClient, TrendMetric } from '../client';
import { ArticleList } from './ArticleList';

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
  const [selected, setSelected] = useState<{ term: string; domain: string; column: 'A' | 'B' | 'Intersect' } | null>(null);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const checkMobile = () => setIsMobile(window.matchMedia('(max-width: 799px)').matches);
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

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

  useEffect(() => {
    setSelected(null);
  }, [domain.domain, compareDomain]);

  const handleSearch = (term: string, searchDomain: string) => {
    const query = encodeURIComponent(`${term} site:${searchDomain}`);
    const baseUrl = searchEngineUrl.replace(/\/$/, '');
    window.open(`${baseUrl}/search?q=${query}`, '_blank');
  };

  const handleTrendClick = (term: string, trendDomain: string, column: 'A' | 'B' | 'Intersect') => {
    if (selected?.term === term && selected?.domain === trendDomain) {
      setSelected(null);
    } else {
      setSelected({ term, domain: trendDomain, column });
    }
  };

  const renderArticleList = () => {
    if (!selected) return null;
    return (
      <ArticleList
        term={selected.term}
        domain={{ domain: selected.domain, language: domain.language }}
        days={days}
        titleOverride={`${t('trends.articles', 'Articles')} / ${selected.domain} / ${selected.term}`}
      />
    );
  };

  // Helper to render a single trend item, handling different column types
  const renderTrendItem = (item: TrendMetric | DomainComparison, columnType: 'A' | 'B' | 'Intersect') => {
    const isDomainComparison = (item as DomainComparison).classification !== undefined;
    const trendTopic = isDomainComparison ? (item as DomainComparison).trend_topic : (item as TrendMetric).trend_topic;
    
    const isActive = selected?.term === trendTopic && selected?.column === columnType;
    
    let scoreA: number | string = '';
    let scoreB: number | string = '';
    const currentDomainName = domain.domain;
    const comparisonDomainName = compareDomain || '';

    if (isDomainComparison) {
      const dcItem = item as DomainComparison;
      scoreA = dcItem.score_a;
      scoreB = dcItem.score_b;
    } else { // TrendMetric for baseItems (Column A when not comparing)
      scoreA = (item as TrendMetric).outlier_ratio.toFixed(1) + 'x';
      scoreB = ''; 
    }

    let displayDomainName = '';
    let displayScore: number | string = ''; 
    let clickDomain = '';

    if (columnType === 'A') { 
      displayDomainName = currentDomainName;
      displayScore = scoreA;
      clickDomain = currentDomainName;
    } else if (columnType === 'B') { 
      displayDomainName = comparisonDomainName;
      displayScore = scoreB; 
      clickDomain = comparisonDomainName;
    } else { 
      // For Intersect, these are handled by individual chips, no need to set here
      // but ensuring they are explicitly empty for safety
      displayDomainName = ''; 
      displayScore = '';      
      clickDomain = '';       
    }


    // Determine the main click handler for the li element
    const liClickHandler = columnType !== 'Intersect' 
      ? () => handleTrendClick(trendTopic, clickDomain, columnType) 
      : undefined; // Intersect clicks are handled by individual chips

    return (
      <li 
        key={`${trendTopic}-${isDomainComparison ? (item as DomainComparison).classification : 'base'}`}
        className={`compare-item ${isActive ? 'selected' : ''} ${columnType === 'Intersect' ? 'shared' : ''}`}
        onClick={liClickHandler}
      >
        <div className="item-content">
          <span className="topic-name">{trendTopic}</span>
          <div className="item-sources">
            {columnType !== 'Intersect' ? (
              <div 
                className={`source-chip ${isActive ? 'active' : ''}`}
                onClick={liClickHandler} 
              >
                <span className="source-name">{displayDomainName}</span>
                <span className="source-score">{displayScore}</span>
              </div>
            ) : (
              /* Intersect column: two clickable chips */
              <>
                <div 
                  className={`source-chip ${isActive && selected?.domain === currentDomainName ? 'active' : ''}`}
                  onClick={() => handleTrendClick(trendTopic, currentDomainName, columnType)}
                >
                  <span className="source-name">{currentDomainName}</span>
                  <span className="source-score">{scoreA}</span>
                </div>
                {compareDomain && (
                  <div 
                    className={`source-chip ${isActive && selected?.domain === comparisonDomainName ? 'active' : ''}`}
                    onClick={() => handleTrendClick(trendTopic, comparisonDomainName, columnType)}
                  >
                    <span className="source-name">{comparisonDomainName}</span>
                    <span className="source-score">{scoreB}</span>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
        <div className="item-actions">
          {(columnType === 'A' || columnType === 'Intersect') && (
            <OpenIcon onClick={() => handleSearch(trendTopic, currentDomainName)} tooltip={<>{trendTopic} {BULLET_DELIMITER} {currentDomainName}</>} />
          )}
          {compareDomain && (columnType === 'B' || columnType === 'Intersect') && (
            <OpenIcon onClick={() => handleSearch(trendTopic, comparisonDomainName)} tooltip={<>{trendTopic} {BULLET_DELIMITER} {comparisonDomainName}</>} />
          )}
        </div>
      </li>
    );
  };

  const uniqueA = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_A') : [];
  const intersect = compareDomain ? items.filter(i => i.classification === 'INTERSECT') : [];
  const uniqueB = compareDomain ? items.filter(i => i.classification === 'BLINDSPOT_B') : [];

  const renderList = (list: DomainComparison[], column: 'A' | 'B' | 'Intersect') => (
    <>
      <ul className="compare-list">
        {list.map((item, idx) => (
          <Fragment key={`${item.trend_topic}-${idx}-frag`}>
            {renderTrendItem(item, column)}
          </Fragment>
        ))}
        {list.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>{t('trends.none', 'None')}</li>}
      </ul>
      {isMobile && selected?.column === column && renderArticleList()}
    </>
  );

  const renderBaseList = () => (
    <>
      <ul className="compare-list">
        {baseItems.slice(0, 10).map((item) => (
          <Fragment key={`${item.trend_topic}-base-frag`}>
            {renderTrendItem(item, 'A')}
          </Fragment>
        ))}
        {baseItems.length === 0 && <li className="compare-item" style={{color: 'var(--secondary-text)', fontStyle: 'italic'}}>{t('trends.none', 'None')}</li>}
      </ul>
      {isMobile && selected?.column === 'A' && renderArticleList()}
    </>
  );

  if (loading) {
    return <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="spinner-small" />
    </div>;
  }

  if (availableDomains.length === 0) {
    return (
      <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--secondary-text)', textAlign: 'center', padding: '20px' }}>
        {t('trends.compare.no_domains', 'There are no domains available for the selected language "{{language}}".', { language: domain.language })}
      </div>
    );
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
          {compareDomain ? renderList(uniqueA, 'A') : renderBaseList()}
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
          {renderList(uniqueB, 'B')}
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
          {renderList(intersect, 'Intersect')}
        </div>
      </div>
      {!isMobile && renderArticleList()}
    </>
  );
};
