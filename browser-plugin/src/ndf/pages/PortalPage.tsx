import '../../shared/i18n';

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { TabPortal } from '../components/TabPortal';
import { TabTrend } from '../components/TabTrend';

interface PortalPageProps {
  items: AnalyzedItem[];
  domain: string;
}

export const PortalPage = ({ items, domain }: PortalPageProps) => {
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

        {activeTab === 'portal' ? <TabPortal items={items} /> : <TabTrend domain={domain} />}

        <div className={`footer-container ${activeTab}`}>
          <Footer />
        </div>
      </div>
    </>
  );
};
