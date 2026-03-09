import '@browser/i18n';

import { Footer } from '@frontend-shared/components/Footer';
import { TabPortal } from '@frontend-shared/components/TabPortal';
import { TabTrend } from '@frontend-shared/components/TabTrend';
import log from '@frontend-shared/logger';
import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import { AnalyzedItem, DomainEntry } from '../ndf-api';

interface PortalPageProps {
  items: AnalyzedItem[];
  domain: DomainEntry;
  availableDomains: DomainEntry[];
  searchEngineUrl: string;
  api: NewsDeframerApi;
}

export const PortalPage = ({ items, domain, availableDomains, searchEngineUrl, api }: PortalPageProps) => {
  const { t } = useTranslation();
  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;
  const [activeTab, setActiveTab] = useState<'portal' | 'trends'>('portal');

  const bypassAndReload = () => {
    window.scrollTo(0, 0);
    log.debug('Bypassing for this session and reloading.');
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
  };

  return (
    <>
      <div className="portal-container">
        <div className="header">
          <h1>
            <span className="title-desktop">{t('portal.title')}</span>
            <span className="title-mobile">{t('portal.title')}</span>
          </h1>
          <button id="btn-hide" className="btn" onClick={bypassAndReload}>{t('portal.hide')}</button>
        </div>

        <div className="tabs">
          <button
            className={`tab-btn ${activeTab === 'portal' ? 'active' : ''}`}
            onClick={() => setActiveTab('portal')}
          >
            { rootDomain }
          </button>
          <button
            className={`tab-btn ${activeTab === 'trends' ? 'active' : ''}`}
            onClick={() => setActiveTab('trends')}
          >
            {t('portal.trends')}
          </button>
        </div>

        {activeTab === 'portal' ? <TabPortal items={items} /> : <TabTrend api={api} domain={domain} availableDomains={availableDomains} searchEngineUrl={searchEngineUrl} />}

        <div className={`footer-container ${activeTab}`}>
          <Footer />
        </div>
      </div>
    </>
  );
};
