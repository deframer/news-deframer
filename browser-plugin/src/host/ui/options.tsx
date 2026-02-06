import '../../shared/i18n';

import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { getDomain } from 'tldts';

import { classifyUrl, PageType } from '../../ndf/utils/url-classifier';
import { invalidateDomainCache } from '../../shared/domain-cache';
import log from '../../shared/logger';
import {
  DEFAULT_BACKEND_URL,
  getSettings,
  Settings,
} from '../../shared/settings';
import { getThemeCss, globalStyles, Theme } from '../../shared/theme';
import { ProxyResponse } from '../../shared/types';
import { ToggleSwitch } from './ToggleSwitch';

type Status = 'idle' | 'loading' | 'success' | 'error';

type SettingsWithSearch = Settings & { searchEngineUrl: string };

// This file is the component, index.tsx is the entry point

export const Options = () => {
  const { t, i18n } = useTranslation();
  const [settings, setSettings] = useState<SettingsWithSearch>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
    enabled: false,
    theme: 'system',
    searchEngineUrl: 'https://search.brave.com',
  });
  const [status, setStatus] = useState<Status>('idle');
  const [loaded, setLoaded] = useState(false);
  const [domains, setDomains] = useState<string[]>([]);
  const [lang, setLang] = useState<string>('default');

  // Load settings on mount (only from chrome.storage, no network calls)
  useEffect(() => {
    log.debug('Loading settings...');
    // Load settings AND language
    Promise.all([getSettings(), chrome.storage.local.get('ndf_language')]).then(([loadedSettings, storageResult]) => {
      log.debug('Settings loaded:', loadedSettings);
      const settingsWithSearch = loadedSettings as SettingsWithSearch;
      // Default searchEngineUrl if missing
      if (!settingsWithSearch.searchEngineUrl) {
        settingsWithSearch.searchEngineUrl = 'https://search.brave.com';
      }
      setSettings(settingsWithSearch);
      const storage = storageResult as { ndf_language?: string };
      setLang(storage.ndf_language || 'default');
      setLoaded(true);

      // Show "Checking..." status and test connection if extension is enabled
      // Delay slightly to ensure dialog window is fully open
      if (loadedSettings.enabled) {
        log.debug('Extension enabled, dialog window should be fully open now, starting connection check in 100ms');
        setTimeout(() => {
          log.debug('Starting connection test...');
          testConnection(loadedSettings);
        }, 100);
      }
      // Don't auto-test connection - only when user clicks Test button
    }).catch((error) => {
      log.error('Failed to load settings:', error);
    });
  }, []);

  // Autosave when settings change (skip initial load) - NO REFRESH for regular changes
  useEffect(() => {
    if (!loaded) return;

    const handler = setTimeout(async () => {
      await chrome.storage.local.set(settings);
    }, 500);

    return () => clearTimeout(handler);
  }, [settings, loaded]);

  // Save immediately and refresh tab when specific user interactions require it
  const saveAndRefresh = async (newSettings: Settings, checkDomain = true, overrideDomains?: string[], skipScripting = false) => {
    await chrome.storage.local.set(newSettings);
    const [tab] = await chrome.tabs.query({
      active: true,
      currentWindow: true,
    });
    if (tab?.id && tab.url?.match(/^https?:\/\//)) {
      let hasNdf = false;
      if (!skipScripting) {
        try {
          const results = await chrome.scripting.executeScript({
            target: { tabId: tab.id },
            func: () => !!document.getElementById('ndf-root')
          });
          hasNdf = !!(results && results[0] && results[0].result);
        } catch {
          // ignore
        }
      }
      if (hasNdf) {
        chrome.tabs.reload(tab.id);
        return;
      }
      // If not has NDF, check if classified as portal or article
      const url = new URL(tab.url);
      const pageType = classifyUrl(url);
      const useDomains = overrideDomains || domains;
      if (pageType === PageType.PORTAL || pageType === PageType.ARTICLE) {
        if (!checkDomain) {
          chrome.tabs.reload(tab.id);
        } else {
          const siteHost = url.host;
          const rootDomain = getDomain(siteHost.replace(/:\d+$/, ''));
          if (useDomains.includes(siteHost) || (rootDomain && useDomains.includes(rootDomain))) {
            chrome.tabs.reload(tab.id);
          }
        }
      }
    }
  };

  // Update theme on change
  useEffect(() => {
    // Inject shared theme styles
    const styleId = 'ndf-theme-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }
    style.textContent = getThemeCss(settings.theme as Theme) + globalStyles;
  }, [settings.theme]);

  // Update language on change
  useEffect(() => {
    if (lang === 'default') {
      const detected = navigator.language.split('-')[0];
      if (['de', 'en'].includes(detected)) {
        i18n.changeLanguage(detected);
      } else {
        i18n.changeLanguage('en');
      }
    } else {
      i18n.changeLanguage(lang);
    }
  }, [lang, i18n]);

  const handleLangChange = async (val: string) => {
    setLang(val);
    // Save to chrome.storage.local so content scripts can read it
    await chrome.storage.local.set({ ndf_language: val });
    saveAndRefresh(settings);
  };


  const testConnection = async (currentSettings: Settings) => {
    setStatus('loading');
    let loadedDomains: string[] = [];
    try {
      const headers: HeadersInit = {};
      if (currentSettings.username && currentSettings.password) {
        headers['Authorization'] =
          'Basic ' +
          btoa(`${currentSettings.username}:${currentSettings.password}`);
      }

      const url =
        currentSettings.backendUrl.replace(/\/$/, '') + '/api/domains';

      const response = await new Promise<ProxyResponse>((resolve, reject) => {
        chrome.runtime.sendMessage(
          { type: 'PROXY_REQ', url, headers, timeout: 5000 },
          (res: ProxyResponse) => {
            if (chrome.runtime.lastError) {
              reject(new Error(chrome.runtime.lastError.message));
              return;
            }
            resolve(res);
          },
        );
      });

      if (response && response.ok) {
        setStatus('success');
        await invalidateDomainCache();
        // Assume response.data is the domains array or { domains: [...] }, possibly as string
        if (response.data) {
          let data: unknown = response.data;
          if (typeof data === 'string') {
            try {
              data = JSON.parse(data);
            } catch {
              // ignore
            }
          }
          const domainList = Array.isArray(data) ? data as string[] : (data as Record<string, unknown>)?.domains || (data as Record<string, unknown>)?.data;
          if (Array.isArray(domainList)) {
            loadedDomains = domainList;
            setDomains(domainList);
          }
        }
        return { connected: true, domains: loadedDomains };
      } else {
        setStatus('error');
        return { connected: false, domains: [] };
      }
    } catch {
      setStatus('error');
      return { connected: false, domains: [] };
    }
  };

  const handleTestClick = () => {
    if (status === 'loading') {
      return; // Ignore click if test is already in progress
    }
    testConnection(settings);
  };

  const handleEnableToggle = async (enabled: boolean) => {
    const newSettings = { ...settings, enabled };
    setSettings(newSettings);

    // Persist immediately so the reloaded tab sees the new state
    await chrome.storage.local.set(newSettings);

    if (enabled) {
      // On successful enable, test the connection and refresh only if connected and domain matches
      const { connected, domains: loadedDomains } = await testConnection(newSettings);
      if (connected) {
        await saveAndRefresh(newSettings, true, loadedDomains, true);
      }
    } else {
      // On disable, invalidate cache, reload tab if on supported domain, and close popup
      setStatus('idle');
      await invalidateDomainCache();
      const [tab] = await chrome.tabs.query({
        active: true,
        currentWindow: true,
      });
      if (tab?.id && tab.url?.match(/^https?:\/\//)) {
        const url = new URL(tab.url);
        const siteHost = url.host;
        const rootDomain = getDomain(siteHost.replace(/:\d+$/, ''));
        if (domains.includes(siteHost) || (rootDomain && domains.includes(rootDomain))) {
          chrome.tabs.reload(tab.id);
        }
      }
      window.close();
    }
  };

  if (!loaded) return <div style={{ padding: '20px' }}>{t('options.loading')}</div>;

  const isConnected = status === 'success';
  const isError = status === 'error';
  const isLoading = status === 'loading';

  return (
    <div className="options-container">
      <div className="header">
        <h2>News Deframer</h2>
        {/* Status Badge */}
        {status !== 'idle' && (
          <div
            className={`status-badge ${isConnected ? 'connected' : isError ? 'error' : isLoading ? 'testing' : ''}`}
          >
            <span className="status-dot" />
            {isConnected ? t('options.status_connected') : isError ? t('options.status_error') : t('options.status_checking')}
          </div>
        )}
      </div>

      <div className="content-grid">
        <div className="settings-column">
          <div className={`card ${!settings.enabled ? 'disabled' : ''}`}>
            <h3 className="section-title">
              {t('options.section_connection')}
            </h3>

            <div className="form-group">
              <label className="input-label">
                {t('options.label_server_url')}
              </label>
              <input
                type="text"
                className="text-input"
                value={settings.backendUrl}
                onChange={(e) =>
                  setSettings({ ...settings, backendUrl: e.target.value })
                }
                disabled={!settings.enabled}
                placeholder="http://localhost:8080"
              />
            </div>

            <div className="form-group">
              <label className="input-label">
                {t('options.label_username')}{' '}
                <span className="optional-text">
                  {t('options.label_optional')}
                </span>
              </label>
              <input
                type="text"
                className="text-input"
                value={settings.username}
                onChange={(e) =>
                  setSettings({ ...settings, username: e.target.value })
                }
                disabled={!settings.enabled}
              />
            </div>

            <div className="form-group-last">
              <label className="input-label">
                {t('options.label_password')}{' '}
                <span className="optional-text">
                  {t('options.label_optional')}
                </span>
              </label>
              <input
                type="password"
                className="text-input"
                value={settings.password}
                onChange={(e) =>
                  setSettings({ ...settings, password: e.target.value })
                }
                disabled={!settings.enabled}
              />
            </div>

            <button
              className="action-button"
              onClick={handleTestClick}
              disabled={!settings.enabled || isLoading}
            >
              {isLoading ? t('options.btn_testing') : t('options.btn_test_connection')}
            </button>
          </div>

          {/* Search Engine Settings */}
          <div className={`card ${!settings.enabled ? 'disabled' : ''}`}>
            <div className="form-group-last">
              <label className="input-label">
                {t('options.label_search_engine')}
              </label>
              <input
                type="text"
                className="text-input"
                value={settings.searchEngineUrl}
                onChange={(e) =>
                  setSettings({ ...settings, searchEngineUrl: e.target.value })
                }
                disabled={!settings.enabled}
              />
            </div>
          </div>
        </div>

        <div className="settings-column">
          {/* Language Selection */}
          <div className="card">
            <h3 className="section-title">
              {t('options.language_label')}
            </h3>
            <select
              className="select-input"
              value={lang}
              onChange={(e) => handleLangChange(e.target.value)}
            >
              <option value="default">{t('options.default')}</option>
              <option value="en">English</option>
              <option value="de">Deutsch</option>
            </select>
          </div>

          {/* Theme Selection */}
          <div className="card">
            <h3 className="section-title">
              {t('options.theme_label')}
            </h3>
            <div className="theme-toggle-group">
               {(['light', 'dark', 'system'] as const).map((themeMode) => (
                 <button
                   key={themeMode}
                   className={`theme-button ${settings.theme === themeMode ? 'active' : ''}`}
                    onClick={() => {
                      const newSettings = { ...settings, theme: themeMode };
                      setSettings(newSettings);
                      saveAndRefresh(newSettings);
                    }}
                >
                  {themeMode === 'system' ? t('options.default') : t(`options.theme_${themeMode}`)}
                </button>
              ))}
            </div>
          </div>
          {/* General Settings */}
          <div className="card">
             <h3 className="section-title">
              {t('options.section_general')}
            </h3>
            <ToggleSwitch
              id="enable-extension"
              label={t('options.label_enable_extension')}
              checked={settings.enabled}
              onChange={handleEnableToggle}
            />
          </div>

          {/* Project Link */}
          <div className="card">
            <h3 className="section-title">
              {t('options.section_project')}
            </h3>
            <a
              href="https://deframer.github.io/"
              target="_blank"
              rel="noopener noreferrer"
              style={{ color: 'var(--accent-color)', textDecoration: 'none', fontWeight: 500 }}
            >
              {t('footer.github_link')}
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};
