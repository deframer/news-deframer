import '../../shared/i18n';

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { Footer } from '../components/Footer';
import { TabPortal } from '../components/TabPortal';
import { TabTrend } from '../components/TabTrend';

const portalPageCss = `
  :host {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    display: block;
    background-color: var(--bg-color);
    color: var(--text-color);
    min-height: 100vh;
  }
  .container { padding: 2em; }
  .header {
    background-color: var(--header-bg);
    color: var(--text-color);
    padding: 1em;
    text-align: center;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    margin-bottom: 2em;
  }
  h1 { margin: 0; font-size: 1.5em; }
  .title-mobile { display: none; }
  .btn {
    padding: 10px 16px;
    border: 1px solid var(--btn-border);
    border-radius: 8px;
    background-color: var(--btn-bg);
    color: var(--btn-text);
    font-size: 0.95em;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
    text-decoration: none;
    flex-shrink: 0;
  }
  .btn:hover { background-color: var(--btn-hover-bg); }

  #btn-hide {
    position: absolute;
    right: 1.5em;
    top: 50%;
    transform: translateY(-50%);
  }

  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }

  .tabs { display: flex; justify-content: center; margin-bottom: 20px; border-bottom: 1px solid var(--border-color, #ccc); }
  .tab-btn { padding: 12px 24px; background: none; border: none; border-bottom: 3px solid transparent; cursor: pointer; font-size: 1.1em; color: var(--text-color); opacity: 0.6; transition: all 0.2s; }
  .tab-btn:hover { opacity: 1; background-color: rgba(0,0,0,0.05); }
  .tab-btn.active { border-bottom-color: var(--btn-bg); opacity: 1; font-weight: 600; }

  @media (max-width: 799px) {
    :host { padding-top: 80px; }
    .container { padding: 1em; padding-top: 0; }
    .header { position: fixed; top: 0; left: 0; width: 100%; z-index: 1000; margin-bottom: 0; justify-content: space-between; padding: 10px 1em; }
    .title-desktop { display: none; }
    .title-mobile { display: inline; }
    #btn-hide { position: static; transform: none; }
  }

`;

interface PortalPageProps {
  items: AnalyzedItem[];
}

export const PortalPage = ({ items }: PortalPageProps) => {
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
      <style>{portalPageCss}</style>
      <div className="container">
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

        {activeTab === 'portal' ? <TabPortal items={items} /> : <TabTrend />}

        <Footer />
      </div>
    </>
  );
};
