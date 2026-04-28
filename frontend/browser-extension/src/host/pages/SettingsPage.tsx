import '../../shared/i18n';

import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { DomainEntry, NewsDeframerClient } from '../../ndf/client';
import { DEFAULT_BACKEND_URL, getSettings, Settings } from '../../shared/settings';
import { getThemeCss, globalStyles, Theme } from '../../shared/theme';
import { SettingsAbout } from '../components/SettingsAbout';
import { SettingsDomains } from '../components/SettingsDomains';
import { SettingsForm } from '../components/SettingsForm';
import { HostStatus, StatusBadge } from '../components/StatusBadge';
import { testConnection } from '../lib/connection';

export const SettingsPage = () => {
  const { t, i18n } = useTranslation();
  const version = chrome.runtime.getManifest().version;
  const [settings, setSettings] = useState<Settings>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
    enabled: false,
    theme: 'system',
    searchEngineUrl: 'https://search.brave.com',
    selectedDomains: [],
  });
  const [status, setStatus] = useState<HostStatus>('idle');
  const [loaded, setLoaded] = useState(false);
  const [lang, setLang] = useState('default');
  const [activeTab, setActiveTab] = useState<'settings' | 'domains' | 'about'>('domains');
  const [domains, setDomains] = useState<DomainEntry[]>([]);
  const [domainsLoading, setDomainsLoading] = useState(false);
  const [domainsUnavailable, setDomainsUnavailable] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const persistSettings = useCallback(
    async (nextSettings: Settings) => {
      await chrome.storage.local.set(nextSettings);
      await chrome.storage.local.set({ ndf_language: lang });
    },
    [lang]
  );

  const refreshConnection = useCallback(
    async (nextSettings: Settings): Promise<boolean> => {
      if (!nextSettings.enabled) {
        setStatus('idle');
        setErrorMessage(null);
        setDomains([]);
        setDomainsUnavailable(false);
        return true;
      }

      const trimmedUrl = nextSettings.backendUrl?.trim() ?? '';
      if (!trimmedUrl) {
        setStatus('error');
        setErrorMessage(i18n.t('options.invalid_server_url', 'Invalid server URL'));
        setDomains([]);
        setDomainsUnavailable(true);
        return false;
      }

      setStatus('loading');
      setErrorMessage(null);
      try {
        const result = await testConnection(nextSettings);
        setStatus(result.connected ? 'success' : 'error');
        setErrorMessage(result.connected ? null : result.errorMessage || 'Connection failed');
        if (result.connected) {
          setDomains(result.domains);
          setDomainsUnavailable(false);
        } else {
          setDomains([]);
          setDomainsUnavailable(true);
        }
        return result.connected;
      } catch (err: Error) {
        setStatus('error');
        setErrorMessage(err.message);
        setDomains([]);
        setDomainsUnavailable(true);
        return false;
      }
    },
    []
  );

  const loadDomains = useCallback(async (settingsToUse: Settings) => {
    setDomainsLoading(true);
    try {
      const client = new NewsDeframerClient(settingsToUse);
      const nextDomains = await client.getDomains();
      setDomains(nextDomains);
      setDomainsUnavailable(false);
    } catch {
      setDomains([]);
      setDomainsUnavailable(true);
    } finally {
      setDomainsLoading(false);
    }
  }, []);

  useEffect(() => {
    Promise.all([getSettings(), chrome.storage.local.get('ndf_language')]).then(([loadedSettings, storageResult]) => {
      setSettings(loadedSettings);
      const storage = storageResult as { ndf_language?: string };
      setLang(storage.ndf_language || 'default');
      setLoaded(true);

      if (loadedSettings.enabled) {
        setStatus('loading');
        testConnection(loadedSettings)
          .then((result) => {
            setStatus(result.connected ? 'success' : 'error');
            setErrorMessage(result.connected ? null : result.errorMessage || 'Connection failed');
            if (result.connected) {
              setDomains(result.domains);
              setDomainsUnavailable(false);
            } else {
              setDomains([]);
              setDomainsUnavailable(true);
            }
          })
          .catch((err: Error) => {
            setStatus('error');
            setErrorMessage(err.message);
            setDomains([]);
            setDomainsUnavailable(true);
          });
      } else {
        loadDomains(loadedSettings).catch(() => undefined);
      }
    });
  }, [loadDomains]);

  useEffect(() => {
    const styleId = 'ndf-theme-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }
    style.textContent = getThemeCss(settings.theme as Theme) + globalStyles;
  }, [settings.theme]);

  useEffect(() => {
    if (lang === 'default') {
      const detected = navigator.language.split('-')[0];
      i18n.changeLanguage(['de', 'en'].includes(detected) ? detected : 'en');
      return;
    }
    i18n.changeLanguage(lang);
  }, [lang, i18n]);

  const handleLanguageChange = async (value: string) => {
    setLang(value);
  };

  const handleSelectedDomainsChange = async (selectedDomains: string[]) => {
    const nextSettings = { ...settings, selectedDomains };
    setSettings(nextSettings);
    await chrome.storage.local.set({ selectedDomains });
  };

  const handleTestConnection = async (nextSettings: Settings = settings) => {
    const connected = await refreshConnection(nextSettings);
    if (connected) {
      await persistSettings(nextSettings);
    }
  };

  const handleApply = async () => {
    await persistSettings(settings);
    await refreshConnection(settings);
  };

  if (!loaded) {
    return <div className="settings-page-container">{t('options.loading')}</div>;
  }

  return (
    <main className="settings-page-container" aria-labelledby="settings-page-title">
      <div className="header page-header">
        <div>
          <p className="eyebrow">News Deframer - {version}</p>
          <h1 id="settings-page-title">{t('options.settings_title', 'Settings')}</h1>
        </div>
        <StatusBadge
          status={status}
          enabled={settings.enabled}
          labels={{
            connected: t('options.status_connected'),
            error: t('options.status_error'),
            checking: t('options.status_loading'),
            disabled: t('options.status_disabled', 'Disabled'),
          }}
        />
      </div>
      <div className="tabs settings-tabs" role="tablist" aria-label={t('options.settings_title', 'Settings')}>
        <button
          className={`tab-btn ${activeTab === 'domains' ? 'active' : ''}`}
          onClick={() => setActiveTab('domains')}
          type="button"
          role="tab"
          aria-selected={activeTab === 'domains'}
        >
          {t('options.tab_domains')}
        </button>
        <button
          className={`tab-btn ${activeTab === 'settings' ? 'active' : ''}`}
          onClick={() => setActiveTab('settings')}
          type="button"
          role="tab"
          aria-selected={activeTab === 'settings'}
        >
          {t('options.tab_settings')}
        </button>
        <button
          className={`tab-btn ${activeTab === 'about' ? 'active' : ''}`}
          onClick={() => setActiveTab('about')}
          type="button"
          role="tab"
          aria-selected={activeTab === 'about'}
        >
          {t('options.tab_about')}
        </button>
      </div>
      {activeTab === 'settings' ? (
        <SettingsForm
          settings={settings}
          lang={lang}
          status={status}
          errorMessage={errorMessage}
          onSettingsChange={setSettings}
          onLanguageChange={handleLanguageChange}
          onTestConnection={handleTestConnection}
        />
      ) : activeTab === 'about' ? (
        <SettingsAbout />
      ) : (
        <SettingsDomains domains={domains} domainsLoading={domainsLoading} domainsUnavailable={domainsUnavailable} settings={settings} onSelectedDomainsChange={handleSelectedDomainsChange} />
      )}
      {activeTab === 'settings' ? (
        <div className="settings-actions">
          <button className="action-button action-button-enabled settings-close-button" onClick={handleApply} type="button">
            {t('options.apply', 'Apply')}
          </button>
        </div>
      ) : null}
    </main>
  );
};
