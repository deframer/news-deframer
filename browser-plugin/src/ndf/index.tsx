import { useEffect,useState } from 'react';
import { createRoot } from 'react-dom/client';
import { getDomain } from 'tldts';

import i18n from '../shared/i18n';
import log from '../shared/logger';
import { getSettings, Settings } from '../shared/settings';
import { getThemeCss, globalStyles, Theme } from '../shared/theme';
import { AnalyzedItem, DomainEntry, NewsDeframerClient } from './client';
import { Spinner } from './components/Spinner';
import { purgeOverlayElements, resetElementAttributes, shouldKeepNode } from './domGuard';
import { ArticlePage } from './pages/ArticlePage';
import { PortalPage } from './pages/PortalPage';
import { ndfStyles } from './styles';
import { classifyUrl, PageType } from './utils/url-classifier';

const App = ({ theme }: { theme: string }) => {
  const [pageType, setPageType] = useState<PageType | null>(null);
  const [data, setData] = useState<
    AnalyzedItem | AnalyzedItem[] | null
  >(null);
  const [currentDomain, setCurrentDomain] = useState<DomainEntry | null>(null);
  const [availableDomains, setAvailableDomains] = useState<DomainEntry[]>([]);
  const [searchEngineUrl, setSearchEngineUrl] = useState<string>('https://search.brave.com');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleStorageChange = (changes: { [key: string]: chrome.storage.StorageChange }, areaName: string) => {
      if (areaName === 'local' && changes.searchEngineUrl) {
        setSearchEngineUrl(String(changes.searchEngineUrl.newValue));
      }
    };
    chrome.storage.onChanged.addListener(handleStorageChange);
    return () => {
      chrome.storage.onChanged.removeListener(handleStorageChange);
    };
  }, []);

  useEffect(() => {
    // Inject global styles into the Light DOM to fix body margin and background "frame"
    const styleId = 'ndf-global-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }

    const themeCss = getThemeCss(theme as Theme);
    style.textContent = themeCss + globalStyles;
  }, [theme]);

  useEffect(() => {
    const init = async () => {
      // For debugging the spinner, uncomment the following line
      // await new Promise(resolve => setTimeout(resolve, 2000));
      try {
        // We don't need to check settings here, start() already did.
        const settings = await getSettings();
        if (settings.searchEngineUrl) {
          setSearchEngineUrl(settings.searchEngineUrl);
        }
        const client = new NewsDeframerClient(settings);
        const allDomains = await client.getDomains();
        setAvailableDomains(allDomains);

        const siteHost = window.location.host;
        const rootDomain = getDomain(window.location.hostname.replace(/:\d+$/, ''));

        // Check if the site's full host or root domain is registered in the backend.
        let matchedDomain: DomainEntry | null = null;

        const exactMatch = allDomains.find((d) => d.domain === siteHost);
        if (exactMatch) {
          matchedDomain = exactMatch;
        } else if (rootDomain) {
          const found = allDomains.find((d) => d.domain === rootDomain);
          if (found) {
            matchedDomain = found;
          }
        }

        if (matchedDomain) {
          setCurrentDomain(matchedDomain);
        }

        const type = classifyUrl(new URL(window.location.href), matchedDomain?.portal_url);
        setPageType(type);

        if (type === PageType.PORTAL) {

          if (matchedDomain) {
            // Fetch the content using the matched domain from the list.
            const items = await client.getSite(matchedDomain.domain);
            if (items.length > 0) {
              setData(items);
            } else {
              // This can happen if the matched domain has no items.
              log.error(`Attempted to get site data for matched domain "${matchedDomain.domain}", but received no items.`);
              throw new Error(`No items found for this portal (${matchedDomain.domain}).`);
            }
          } else {
            // If neither the full host nor root domain is in the list, the site is not supported.
            const checkedDomains = [siteHost];
            if (rootDomain) checkedDomains.push(rootDomain);
            log.error(`Domain authorization failed: checked [${checkedDomains.join(', ')}] but none found in allDomains=[${allDomains.map(d => d.domain).join(', ')}]`);
            throw new Error(`Domain "${checkedDomains.join('" or "')}" is not configured in the Deframer backend.`);
          }
        } else if (type === PageType.ARTICLE) {
          const articleUrl = window.location.href;
          log.debug(`Attempting to fetch article data for URL: "${articleUrl}"`);
          const item = await client.getItem(articleUrl);
          if (item) {
            log.debug(`Successfully retrieved article data for URL: "${articleUrl}"`);
            setData(item);
          } else {
            log.error(`No article data found for URL: "${articleUrl}" - this URL may not be in the backend database`);
            throw new Error(`No item found for this URL: ${articleUrl}`);
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
    chrome.runtime.sendMessage({ type: 'SET_BYPASS', value: true }, () => {
      window.location.reload();
    });
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
          <style>{getThemeCss(theme as Theme) + globalStyles + ndfStyles}</style>
          <ArticlePage item={data as AnalyzedItem} />
        </>
      );
    case PageType.PORTAL:
      // data is AnalyzedItem[]
      return (
        <>
          <style>{getThemeCss(theme as Theme) + globalStyles + ndfStyles}</style>
          <PortalPage items={data as AnalyzedItem[]} domain={currentDomain!} availableDomains={availableDomains} searchEngineUrl={searchEngineUrl} />
        </>
      );
    default:
      return null;
  }
};

export const start = async (providedSettings?: Settings) => {
  if (sessionStorage.getItem('__ndf-bypass')) {
    sessionStorage.removeItem('__ndf-bypass');
    log.debug('Bypass detected. Not starting NDF.');
    return;
  }

  try {
    const settings = providedSettings || await getSettings();
    if (!settings.enabled) {
      log.debug('NDF is disabled.');
      return;
    }

    // Apply language setting from storage
    const storage = (await chrome.storage.local.get('ndf_language')) as { ndf_language?: string };
    const lang = storage.ndf_language || 'default';
    if (lang !== 'default') {
      await i18n.changeLanguage(lang);
    }
    else {
      const detected = navigator.language.split('-')[0];
      if (['de', 'en'].includes(detected)) {
        await i18n.changeLanguage(detected);
      } else {
        await i18n.changeLanguage('en');
      }
    }

    const type = classifyUrl(new URL(window.location.href));
    if (type === PageType.PORTAL || type === PageType.ARTICLE) {
      log.debug(`Page detected as ${type}.`);
      // window.stop(); // we rely on the DOM guard instead

      const html = document.documentElement;
      if (html) {
        resetElementAttributes(html);
        const domain = getDomain(window.location.hostname) || window.location.hostname;
        html.innerHTML = `<head><title>NDF • ${domain}</title></head><body></body>`;
      }

      // Ensure body exists and is accessible
      let body = document.body;
      if (!body) {
        body = document.createElement('body');
        if (html) {
          html.appendChild(body);
        } else {
          document.appendChild(body);
        }
      }
      resetElementAttributes(body);

      const rootEl = document.createElement('div');
      rootEl.id = 'ndf-root';
      body.appendChild(rootEl);

      // Watch for any Light DOM changes. If a late script tries to add
      // elements or attributes, we remove them so the old page never returns.
      const observer = new MutationObserver((mutations) => {
        let needsCleanup = false;

        mutations.forEach((mutation) => {
          if (
            mutation.type === 'attributes' &&
            (mutation.target === html || mutation.target === document.body)
          ) {
            needsCleanup = true;
            return;
          }

          // If the mutation happens inside our root or styles, ignore it
          if (rootEl.contains(mutation.target)) return;
          if (
            mutation.target instanceof Element &&
            mutation.target.id === 'ndf-global-styles'
          )
            return;

          mutation.addedNodes.forEach((node) => {
            // Allow our critical elements
            if (shouldKeepNode(node, { rootEl })) {
              return;
            }

            if (
              node.parentNode &&
              !(node instanceof Element &&
              (node.id === 'ndf-root' || node.id === 'ndf-global-styles'))
            ) {
              node.parentNode.removeChild(node);
            }
          });
        });

        if (needsCleanup) {
          resetElementAttributes(html);
          resetElementAttributes(document.body);
          purgeOverlayElements();
        }
      });

      observer.observe(document.documentElement, {
        childList: true,
        subtree: true,
        attributes: true,
      });

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
