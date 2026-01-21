import { useEffect, useState } from 'react';

import { invalidateDomainCache } from '../../shared/domain-cache';
import {
  DEFAULT_BACKEND_URL,
  getSettings,
  Settings,
} from '../../shared/settings';
import log from '../../shared/logger';
import { getThemeCss } from '../../shared/theme';
import { ProxyResponse } from '../../shared/types';
import { classifyUrl, PageType } from '../../ndf/utils/url-classifier';
import { getDomain } from 'tldts';
import { Footer } from './Footer';
import { ToggleSwitch } from './ToggleSwitch';

type Status = 'idle' | 'loading' | 'success' | 'error';

// This file is the component, index.tsx is the entry point

export const Options = () => {
  const [settings, setSettings] = useState<Settings>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
    enabled: true,
    theme: 'system',
  });
  const [status, setStatus] = useState<Status>('idle');
  const [loaded, setLoaded] = useState(false);
  const [isDark, setIsDark] = useState(false);
  const [domains, setDomains] = useState<string[]>([]);

  // Load settings on mount (only from chrome.storage, no network calls)
  useEffect(() => {
    log.debug('Loading settings...');
    getSettings().then((loadedSettings) => {
      log.debug('Settings loaded:', loadedSettings);
      setSettings(loadedSettings);
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
        } catch (e) {
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
    let dark: boolean;
    if (settings.theme === 'system') {
      dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    } else {
      dark = settings.theme === 'dark';
    }
    setIsDark(dark);

    // Inject theme styles
    const styleId = 'theme-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }
    style.textContent = getThemeCss(settings.theme);
  }, [settings.theme]);



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
          let data = response.data;
          if (typeof data === 'string') {
            try {
              data = JSON.parse(data);
            } catch (e) {
              // ignore
            }
          }
          const domainList = Array.isArray(data) ? data : (data as any)?.domains || (data as any)?.data;
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

  if (!loaded) return <div style={{ padding: '20px' }}>Loading...</div>;

  const isConnected = status === 'success';
  const isError = status === 'error';
  const isLoading = status === 'loading';

  return (
    <div
      style={{
        padding: '24px',
        width: '780px', // Force width to be wide
        margin: '0 auto',
        color: 'var(--text-color)', // Font is now inherited from body
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
          borderBottom: '1px solid var(--border-color)',
          paddingBottom: '16px',
        }}
      >
        <h2 style={{ margin: 0, fontSize: '20px' }}>News Deframer</h2>
        {/* Status Badge */}
        {status !== 'idle' && (
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              fontSize: '12px',
              fontWeight: 600,
              color: isConnected
                ? 'var(--success-color)'
                : isError
                ? 'var(--danger-color)'
                : 'var(--secondary-text)',
              backgroundColor: 'var(--bg-color)', // Simple bg for now
              border: `1px solid ${
                isConnected
                  ? 'var(--success-color)'
                  : isError
                  ? 'var(--danger-color)'
                  : 'var(--border-color)'
              }`,
              padding: '4px 8px',
              borderRadius: '12px',
              boxShadow: isDark && (isConnected || isError) ? `0 0 8px ${
                isConnected ? 'var(--success-color)' : 'var(--danger-color)'
              }` : 'none',
            }}
          >
            <span
              style={{
                display: 'block',
                width: '6px',
                height: '6px',
                borderRadius: '50%',
                backgroundColor: 'currentColor',
                marginRight: '6px',
              }}
            />
            {isConnected ? 'Connected' : isError ? 'Error' : 'Checking...'}
          </div>
        )}
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '20px', alignItems: 'start', marginBottom: '24px' }}>
        <div
          style={{
            backgroundColor: 'var(--card-bg)',
            border: '1px solid var(--border-color)',
            borderRadius: '8px',
            padding: '16px',
            boxShadow: 'var(--card-shadow)',
            opacity: settings.enabled ? 1 : 0.6,
            transition: 'opacity 0.2s',
            pointerEvents: settings.enabled ? 'auto' : 'none',
          }}
        >
          <h3
            style={{
              marginTop: 0,
              marginBottom: '16px',
              fontSize: '14px',
              textTransform: 'uppercase',
              letterSpacing: '0.05em',
              color: 'var(--secondary-text)',
            }}
          >
            Connection
          </h3>

        <div style={{ marginBottom: '16px' }}>
          <label
            style={{
              display: 'block',
              marginBottom: '6px',
              fontSize: '14px',
              fontWeight: 500,
            }}
          >
            Server URL
          </label>
          <input
            type="text"
            value={settings.backendUrl}
            onChange={(e) =>
              setSettings({ ...settings, backendUrl: e.target.value })
            }
            disabled={!settings.enabled}
            style={{
              width: '100%',
              padding: '10px',
              borderRadius: '6px',
              border: '1px solid var(--btn-border)',
              fontSize: '14px',
              backgroundColor: settings.enabled ? 'var(--bg-color)' : 'var(--rating-bg)',
              color: 'var(--text-color)',
            }}
            placeholder="http://localhost:8080"
          />
        </div>

          <div style={{ marginBottom: '16px' }}>
          <label
            style={{
              display: 'block',
              marginBottom: '6px',
              fontSize: '14px',
              fontWeight: 500,
            }}
          >
            Username{' '}
            <span style={{ fontWeight: 400, color: 'var(--secondary-text)' }}>
              (Optional)
            </span>
          </label>
          <input
            type="text"
            value={settings.username}
            onChange={(e) =>
              setSettings({ ...settings, username: e.target.value })
            }
            disabled={!settings.enabled}
            style={{
              width: '100%',
              padding: '10px',
              borderRadius: '6px',
              border: '1px solid var(--btn-border)',
              fontSize: '14px',
              backgroundColor: settings.enabled ? 'var(--bg-color)' : 'var(--rating-bg)',
              color: 'var(--text-color)',
            }}
          />
        </div>

        <div style={{ marginBottom: '20px' }}>
          <label
            style={{
              display: 'block',
              marginBottom: '6px',
              fontSize: '14px',
              fontWeight: 500,
            }}
          >
            Password{' '}
            <span style={{ fontWeight: 400, color: 'var(--secondary-text)' }}>
              (Optional)
            </span>
          </label>
          <input
            type="password"
            value={settings.password}
            onChange={(e) =>
              setSettings({ ...settings, password: e.target.value })
            }
            disabled={!settings.enabled}
            style={{
              width: '100%',
              padding: '10px',
              borderRadius: '6px',
              border: '1px solid var(--btn-border)',
              fontSize: '14px',
              backgroundColor: settings.enabled ? 'var(--bg-color)' : 'var(--rating-bg)',
              color: 'var(--text-color)',
            }}
          />
        </div>

          <button
            onClick={handleTestClick}
            disabled={!settings.enabled || isLoading}
            style={{
              width: '100%',
              padding: '10px',
              backgroundColor: settings.enabled ? 'var(--btn-bg)' : 'var(--rating-bg)',
              border: '1px solid var(--btn-border)',
              borderRadius: '6px',
              cursor: !settings.enabled || isLoading ? 'not-allowed' : 'pointer',
              fontSize: '14px',
              fontWeight: 500,
              color: settings.enabled ? 'var(--btn-text)' : 'var(--secondary-text)',
              transition: 'background-color 0.2s, color 0.2s',
            }}
          >
            {isLoading ? 'Testing...' : 'Test Connection'}
          </button>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          {/* Theme Selection */}
          <div
            style={{
              backgroundColor: 'var(--card-bg)',
              border: '1px solid var(--border-color)',
              borderRadius: '8px',
              padding: '16px',
              boxShadow: 'var(--card-shadow)',
            }}
          >
            <h3
              style={{
                marginTop: 0,
                marginBottom: '16px',
                fontSize: '14px',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                color: 'var(--secondary-text)',
              }}
            >
              Theme
            </h3>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
               {(['light', 'dark', 'system'] as const).map((theme) => (
                 <button
                   key={theme}
                    onClick={() => {
                      const newSettings = { ...settings, theme };
                      setSettings(newSettings);
                      saveAndRefresh(newSettings);
                    }}
                   style={{
                    padding: '10px',
                    border: `1px solid ${
                      settings.theme === theme
                        ? 'var(--accent-color)'
                        : 'var(--btn-border)'
                    }`,
                    borderRadius: '6px',
                    backgroundColor:
                      settings.theme === theme ? 'var(--accent-color)' : 'var(--btn-bg)',
                    color:
                      settings.theme === theme
                        ? 'var(--accent-text)'
                        : 'var(--btn-text)',
                    cursor: 'pointer',
                    fontSize: '14px',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    gap: '8px',
                  }}
                >
                  {theme === 'light' && '☀️ Light Mode'}
                  {theme === 'dark' && '🌙 Dark Mode'}
                  {theme === 'system' && '💻 System Default'}
                </button>
              ))}
            </div>
          </div>
          {/* General Settings */}
          <div
            style={{
              backgroundColor: 'var(--card-bg)',
              border: '1px solid var(--border-color)',
              borderRadius: '8px',
              padding: '16px',
              boxShadow: 'var(--card-shadow)',
            }}
          >
             <h3
              style={{
                marginTop: 0,
                marginBottom: '16px',
                fontSize: '14px',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                color: 'var(--secondary-text)',
              }}
            >
              General
            </h3>
            <ToggleSwitch
              id="enable-extension"
              label="Enable Extension"
              checked={settings.enabled}
              onChange={handleEnableToggle}
            />
          </div>
        </div>
      </div>
      <Footer />
    </div>
  );
};
