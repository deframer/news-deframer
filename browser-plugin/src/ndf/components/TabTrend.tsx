import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { TrendComparisonMetric, TrendRepo } from './TrendRepo';
import { Footer } from './Footer';
import { TrendItem, TrendTop } from './TrendTop';
import { TrendCompare } from './TrendCompare';
import { TrendTopTagCloud } from './TrendTopTagCloud';
import { TrendLifecycle } from './TrendLifecycle';

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
  { id: '24h', label: 'Last 24h' },
  { id: '7d', label: 'Last 7 days' },
  { id: '30d', label: 'Last 30 days' },
  { id: '90d', label: 'Last 90 days' },
  { id: '365d', label: 'Last 365 days' },
];

export const TabTrend = () => {
  const { t } = useTranslation();
  const [viewMode, setViewMode] = useState<'list' | 'cloud' | 'compare' | 'lifecycle'>('list');
  const [timeRange, setTimeRange] = useState('7d');
  const [items, setItems] = useState<TrendItem[]>([]);
  const [compareItems, setCompareItems] = useState<TrendComparisonMetric[]>([]);

  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;
  const availableDomains = TrendRepo.getAvailableDomains().filter(d => d.id !== rootDomain);
  const [compareDomain, setCompareDomain] = useState<string | null>(availableDomains[0]?.id || null);

  log.info(`trend analysis - current domain: ${rootDomain}, range: ${timeRange}`);

  useEffect(() => {
    const fetchData = async () => {
      // Always fetch base trends so we have them for the "Our Topics" column in compare mode
      const data = await TrendRepo.getTrends(rootDomain /*, timeRange */);
      const mappedItems: TrendItem[] = data.map((d, index) => ({
        word: d.trend_topic,
        rank: index + 1,
        count: d.frequency,
        utility: d.utility,
        outlierRatio: d.outlier_ratio
      }));
      setItems(mappedItems);

      if (viewMode === 'compare' && compareDomain) {
        const data = await TrendRepo.getTrendComparison(rootDomain, compareDomain, timeRange);
        setCompareItems(data);
      }
    };
    fetchData();
  }, [rootDomain, timeRange, viewMode, compareDomain]);

  return (
    <div className="trend-container">
      <style>{tabTrendCss}</style>

      <div className="trend-header">
        {/* 1. Main Navigation Tabs */}
        <div className="nav-tabs">
          <button
            className={`nav-tab ${viewMode === 'list' ? 'active' : ''}`}
            onClick={() => setViewMode('list')}
          >
            {t('trends.list') || 'List'}
          </button>
          <button
            className={`nav-tab ${viewMode === 'cloud' ? 'active' : ''}`}
            onClick={() => setViewMode('cloud')}
          >
            {t('trends.cloud') || 'Cloud'}
          </button>
          <button
            className={`nav-tab ${viewMode === 'compare' ? 'active' : ''}`}
            onClick={() => setViewMode('compare')}
          >
            {t('trends.compare') || 'Compare'}
          </button>
          <button
            className={`nav-tab ${viewMode === 'lifecycle' ? 'active' : ''}`}
            onClick={() => setViewMode('lifecycle')}
          >
            {t('trends.lifecycle') || 'Lifecycle'}
          </button>
        </div>

        {/* 2. Filter Bar (Time Selection) */}
        <div className="filter-bar">
          <span className="filter-label">Time:</span>
          <div className="time-selector">
            {TIME_RANGES.map((range) => (
              <button
                key={range.id}
                className={`time-btn ${timeRange === range.id ? 'active' : ''}`}
                onClick={() => setTimeRange(range.id)}
              >
                {range.id}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* 3. Content Area */}
      <div className="trend-content">
        {viewMode === 'list' && <TrendTop items={items} />}

        {viewMode === 'cloud' && <TrendTopTagCloud items={items} />}

        {viewMode === 'compare' && (
          <TrendCompare
            items={compareItems}
            baseItems={items}
            currentDomain={rootDomain}
            compareDomain={compareDomain}
            availableDomains={availableDomains}
            onSelectDomain={setCompareDomain}
          />
        )}

        {viewMode === 'lifecycle' && <TrendLifecycle domain={rootDomain} timeRange={timeRange} />}

        <div className="trend-footer">
          <Footer />
        </div>
      </div>
    </div>
  );
};