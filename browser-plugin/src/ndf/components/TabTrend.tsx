import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import log from '../../shared/logger';
import { DomainEntry } from '../client';
import { Footer } from './Footer';
import { TrendCompare } from './TrendCompare';
import { TrendSearch } from './TrendSearch';
import { TrendTagCloud } from './TrendTagCloud';

const TIME_RANGES = [
  { id: '24h', days: 1, label: 'trends.time_ranges.last_24h' },
  { id: '7d', days: 7, label: 'trends.time_ranges.last_7d' },
  { id: '30d', days: 30, label: 'trends.time_ranges.last_30d' },
  { id: '90d', days: 90, label: 'trends.time_ranges.last_90d' },
  { id: '365d', days: 365, label: 'trends.time_ranges.last_365d' },
];

export const TabTrend = ({ domain, availableDomains, searchEngineUrl }: { domain: DomainEntry; availableDomains: DomainEntry[]; searchEngineUrl: string }) => {
  const { t } = useTranslation();
  const [viewMode, setViewMode] = useState<'cloud' | 'compare' | 'lifecycle'>('cloud');
  const [timeRange, setTimeRange] = useState('7d');

  const domainOptions = availableDomains.filter(d => d.domain !== domain.domain && d.language === domain.language).map(d => ({ id: d.domain, name: d.domain }));
  const [compareDomain, setCompareDomain] = useState<string | null>(domainOptions[0]?.id || null);

  const currentDays = TIME_RANGES.find(r => r.id === timeRange)?.days || 7;
  log.info(`trend analysis - current domain: ${domain.domain}, days: ${currentDays}`);

  return (
    <div className="trend-container">
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
            {t('trends.search', 'Search')}
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
        {viewMode === 'cloud' && <TrendTagCloud domain={domain} days={currentDays} />}

        {viewMode === 'compare' && (
          <TrendCompare
            days={currentDays}
            baseItems={[]}
            compareDomain={compareDomain}
            availableDomains={domainOptions}
            onSelectDomain={setCompareDomain}
            domain={domain}
            searchEngineUrl={searchEngineUrl}
          />
        )}

        {viewMode === 'lifecycle' && <TrendSearch domain={domain} days={currentDays} />}

        <div className="trend-footer">
          <Footer />
        </div>
      </div>
    </div>
  );
};