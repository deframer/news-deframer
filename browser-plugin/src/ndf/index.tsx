import React from 'react';
import { createRoot } from 'react-dom/client';
import { getDomain } from 'tldts';

import log from '../shared/logger';
import { getSettings } from '../shared/settings';
import { getThemeCss, globalStyles, Theme } from '../shared/theme';
import { AnalyzedItem, NewsDeframerClient } from './client';
import { Spinner } from './components/Spinner';
import { ArticlePage } from './pages/ArticlePage';
import { PortalPage } from './pages/PortalPage';
import { classifyUrl, PageType } from './utils/url-classifier';

const App = ({ theme }: { theme: Theme }) => {
  const [pageType, setPageType] = React.useState<PageType | null>(null);
  const [data, setData] = React.useState<
    AnalyzedItem | AnalyzedItem[] | null
  >(null);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    // Inject global styles into the Light DOM to fix body margin and background "frame"
    const styleId = 'ndf-global-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }

    const themeCss = getThemeCss(theme);
    style.textContent = themeCss + globalStyles;
  }, [theme]);

  React.useEffect(() => {
    const init = async () => {
      try {
        // We don't need to check settings here, start() already did.
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const type = classifyUrl(new URL(window.location.href));
        setPageType(type);

        if (type === PageType.PORTAL) {
          const rootDomain = getDomain(window.location.hostname);
          if (!rootDomain) throw new Error('Could not determine root domain.');
          const items = await client.getSite(rootDomain);
          if (items.length > 0) {
            setData(items);
          } else {
            throw new Error('No items found for this portal.');
          }
        } else if (type === PageType.ARTICLE) {
          const item = await client.getItem(window.location.href);
          if (item) {
            setData(item);
          } else {
            throw new Error('No item found for this URL.');
          }
        }
      } catch (err: unknown) {
        log.error('Failed to initialize NDF:', err);
        const message =
          err instanceof Error ? err.message : 'An unknown error occurred.';
        setError(message);
      }
    };
    init();
  }, []);

  if (error) {
    // If there was an error, bypass and reload to see the original page.
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
    return null; // Render nothing while we reload
  }

  if (!pageType || !data) {
    // Render a loading state, as the page is currently blank.
    return <Spinner />;
  }

  switch (pageType) {
    case PageType.ARTICLE:
      // data is AnalyzedItem
      return (
        <>
          <style>{getThemeCss(theme)}</style>
          <ArticlePage item={data as AnalyzedItem} />
        </>
      );
    case PageType.PORTAL:
      // data is AnalyzedItem[]
      return (
        <>
          <style>{getThemeCss(theme)}</style>
          <PortalPage items={data as AnalyzedItem[]} />
        </>
      );
    default:
      return null;
  }
};

const start = async () => {
  if (sessionStorage.getItem('__ndf-bypass')) {
    sessionStorage.removeItem('__ndf-bypass');
    log.info('Bypass detected. Not starting NDF.');
    return;
  }

  try {
    const settings = await getSettings();
    if (!settings.enabled) {
      log.info('NDF is disabled.');
      return;
    }

    const type = classifyUrl(new URL(window.location.href));
    if (type === PageType.PORTAL || type === PageType.ARTICLE) {
      log.info(`Page detected as ${type}. Stopping window immediately.`);
      window.stop();

      // Create a clean slate
      document.documentElement.innerHTML =
        '<head><title>News Deframer</title></head><body></body>';

      const rootEl = document.createElement('div');
      rootEl.id = 'ndf-root';
      document.body.appendChild(rootEl);

      const shadowRoot = rootEl.attachShadow({ mode: 'open' });
      const appContainer = document.createElement('div');
      shadowRoot.appendChild(appContainer);

      const root = createRoot(appContainer);
      root.render(<App theme={settings.theme} />);
    }
  } catch (e) {
    log.error('Could not start NDF', e);
  }
};

start();
