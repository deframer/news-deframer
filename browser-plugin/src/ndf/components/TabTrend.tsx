import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { Footer } from './Footer';
import { TrendCompare } from './TrendCompare';
import { TrendLifecycle } from './TrendLifecycle';
import { TrendComparisonMetric, TrendRepo } from './TrendRepo';
import { TrendTopTagCloud } from './TrendTopTagCloud';

const tabTrendCss = `
  .trend-container {
    display: flex;
    flex-direction: column;
    background: var(--card-bg, #fff);
    min-height: 400px;
  }
  .trend-footer { display: none; }
  @media (max-width: 799px) {
    .trend-container {
      height: calc(100vh - 160px);
      min-height: 0;
    }
    .trend-footer { display: block; margin-top: 20px; }
  }

  .trend-header {
    background: var(--card-bg, #fff);
    flex-shrink: 0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
    z-index: 5;
    position: relative;
  }

  /* 1. Top Navigation Tabs */
  .nav-tabs {
    display: flex;
    border-bottom: 1px solid var(--border-color, #eee);
    background: var(--card-bg, #fff);
    overflow-x: auto;
    scrollbar-width: none;
  }
  .nav-tabs::-webkit-scrollbar { display: none; }
  .nav-tab {
    flex: 1;
    padding: 14px 10px;
    border: none;
    background: none;
    cursor: pointer;
    font-size: 0.95em;
    font-weight: 500;
    color: var(--secondary-text);
    border-bottom: 3px solid transparent;
    transition: all 0.2s;
    text-align: center;
    white-space: nowrap;
    min-width: fit-content;
  }
  .nav-tab:hover {
    color: var(--primary-color, #0056b3);
    background-color: var(--hover-bg, rgba(0,0,0,0.02));
  }
  .nav-tab.active {
    color: var(--primary-color, #0056b3);
    border-bottom-color: var(--primary-color, #0056b3);
    font-weight: 600;
  }
  @media (max-width: 799px) {
    .nav-tab {
      padding: 12px 5px;
      font-size: 0.9em;
    }
  }

  /* 2. Filter Bar (Time Selection) */
  .filter-bar {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color, #f0f0f0);
    display: flex;
    align-items: center;
    justify-content: flex-end;
    background: var(--bg-color-secondary, #fafafa);
    gap: 10px;
  }
  .filter-label {
    font-size: 0.85em;
    font-weight: 600;
    color: var(--secondary-text);
  }
  .time-selector {
    display: flex;
    gap: 2px;
    background: var(--border-color, #e0e0e0);
    padding: 2px;
    border-radius: 6px;
  }
  .time-btn {
    border: none;
    background: transparent;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 0.8em;
    cursor: pointer;
    color: var(--secondary-text);
    font-weight: 500;
    transition: all 0.1s;
  }
  .time-btn:hover {
    color: var(--text-color);
  }
  .time-btn.active {
    background: var(--card-bg, #fff);
    color: var(--text-color);
    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
    font-weight: 600;
  }

  /* 3. Content Area */
  .trend-content {
    padding: 16px;
    flex: 1;
    overflow-y: auto;
  }
  .domain-select {
    padding: 6px 10px;
    border-radius: 6px;
    border: 1px solid var(--border-color, #ddd);
    background: var(--card-bg, #fff);
    color: var(--text-color);
    font-size: 0.9em;
  }
`;

const TIME_RANGES = [
  { id: '24h', days: 1, label: 'trends.time_ranges.last_24h' },
  { id: '7d', days: 7, label: 'trends.time_ranges.last_7d' },
  { id: '30d', days: 30, label: 'trends.time_ranges.last_30d' },
  { id: '90d', days: 90, label: 'trends.time_ranges.last_90d' },
  { id: '365d', days: 365, label: 'trends.time_ranges.last_365d' },
];

export interface TrendItem {
  word: string;
  rank: number;
  count: number;
  utility: number;
  outlierRatio: number;
}

export const TabTrend = ({ domain }: { domain: string }) => {
  const { t } = useTranslation();
  const [viewMode, setViewMode] = useState<'cloud' | 'compare' | 'lifecycle'>('cloud');
  const [timeRange, setTimeRange] = useState('7d');
  const [items, setItems] = useState<TrendItem[]>([]);
  const [compareItems, setCompareItems] = useState<TrendComparisonMetric[]>([]);

  // Use getDomain to normalize the passed domain for the exclusion list (which contains root domains)
  const rootDomainForFilter = getDomain(domain) || domain;
  const availableDomains = TrendRepo.getAvailableDomains(rootDomainForFilter).filter(d => d.id !== rootDomainForFilter);
  const [compareDomain, setCompareDomain] = useState<string | null>(availableDomains[0]?.id || null);

  const currentDays = TIME_RANGES.find(r => r.id === timeRange)?.days || 7;
  log.info(`trend analysis - current domain: ${domain}, days: ${currentDays}`);

  useEffect(() => {
    const fetchData = async () => {
      // Always fetch base trends so we have them for the "Our Topics" column in compare mode
      const data = await TrendRepo.getTrends(domain, currentDays);
      const mappedItems: TrendItem[] = data.map((d, index) => ({
        word: d.trend_topic,
        rank: index + 1,
        count: d.frequency,
        utility: d.utility,
        outlierRatio: d.outlier_ratio
      }));
      setItems(mappedItems);

      if (viewMode === 'compare' && compareDomain) {
        const data = await TrendRepo.getTrendComparison(domain, compareDomain, currentDays);
        setCompareItems(data);
      }
    };
    fetchData();
  }, [domain, timeRange, viewMode, compareDomain, currentDays]);

  return (
    <div className="trend-container">
      <style>{tabTrendCss}</style>

      <div style={{ padding: '8px', textAlign: 'center', fontSize: '0.85em', color: 'var(--secondary-text)', background: 'var(--bg-color-secondary, #f5f5f5)' }}>
        This is dummy data
      </div>

      <div className="trend-header">
        {/* 1. Main Navigation Tabs */}
        <div className="nav-tabs">
          <button
            className={`nav-tab ${viewMode === 'cloud' ? 'active' : ''}`}
            onClick={() => setViewMode('cloud')}
          >
            {t('trends.cloud', 'Cloud')}
          </button>
          <button
            className={`nav-tab ${viewMode === 'compare' ? 'active' : ''}`}
            onClick={() => setViewMode('compare')}
          >
            {t('trends.compare_view', 'Compare')}
          </button>
          <button
            className={`nav-tab ${viewMode === 'lifecycle' ? 'active' : ''}`}
            onClick={() => setViewMode('lifecycle')}
          >
            {t('trends.lifecycle', 'Lifecycle')}
          </button>
        </div>

        {/* 2. Filter Bar (Time Selection) */}
        <div className="filter-bar">
          <span className="filter-label">{t('trends.time_label', 'Time:')}</span>
          <div className="time-selector">
            {TIME_RANGES.map((range) => (
              <button
                key={range.id}
                className={`time-btn ${timeRange === range.id ? 'active' : ''}`}
                onClick={() => setTimeRange(range.id)}
              >
                {t(range.label)}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* 3. Content Area */}
      <div className="trend-content">
        {viewMode === 'cloud' && <TrendTopTagCloud items={items} days={currentDays} domain={domain} />}

        {viewMode === 'compare' && (
          <TrendCompare
            items={compareItems}
            baseItems={items}
            compareDomain={compareDomain}
            availableDomains={availableDomains}
            onSelectDomain={setCompareDomain}
          />
        )}

        {viewMode === 'lifecycle' && <TrendLifecycle domain={domain} days={currentDays} />}

        <div className="trend-footer">
          <Footer />
        </div>
      </div>
    </div>
  );
};