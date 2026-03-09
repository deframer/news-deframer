import '@browser/i18n';

import { SettingsContent } from '@frontend-shared/settings/SettingsContent';
import { ConnectionStatus, StatusBadge } from '@frontend-shared/settings/StatusBadge';
import { getThemeCss, globalStyles, Theme } from '@frontend-shared/theme';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { testConnection } from '../lib/connection';
import { DEFAULT_BACKEND_URL, getSettings, Settings } from '../lib/settings-store';

export const SettingsPage = () => {
  const { t, i18n } = useTranslation();
  const [settings, setSettings] = useState<Settings>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
    enabled: false,
    theme: 'system',
    searchEngineUrl: 'https://search.brave.com',
  });
  const [status, setStatus] = useState<ConnectionStatus>('idle');
  const [loaded, setLoaded] = useState(false);
  const [lang, setLang] = useState('default');

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
          })
          .catch(() => {
            setStatus('error');
          });
      }
    });
  }, []);

  useEffect(() => {
    if (!loaded) return;
    const handler = setTimeout(async () => {
      await chrome.storage.local.set(settings);
    }, 300);
    return () => clearTimeout(handler);
  }, [settings, loaded]);

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
    await chrome.storage.local.set({ ndf_language: value });
  };

  const handleTestConnection = async () => {
    setStatus('loading');
    try {
      const result = await testConnection(settings);
      setStatus(result.connected ? 'success' : 'error');
    } catch {
      setStatus('error');
    }
  };

  const handleClose = () => {
    window.close();
  };

  if (!loaded) {
    return <div className="settings-page-container">{t('options.loading')}</div>;
  }

  return (
    <main className="settings-page-container" aria-labelledby="settings-page-title">
      <div className="header page-header">
        <div>
          <p className="eyebrow">News Deframer</p>
          <h1 id="settings-page-title">{t('options.settings_title', 'Settings')}</h1>
        </div>
        <StatusBadge
          status={status}
          enabled={settings.enabled}
          labels={{
            connected: t('options.status_connected'),
            error: t('options.status_error'),
            checking: t('options.status_checking'),
            disabled: t('options.status_disabled', 'Disabled'),
          }}
        />
      </div>
      <SettingsContent settings={settings} lang={lang} status={status} onSettingsChange={setSettings} onLanguageChange={handleLanguageChange} onTestConnection={handleTestConnection} />
      <div className="settings-actions">
        <button className="action-button action-button-enabled settings-close-button" onClick={handleClose} type="button">
          {t('options.apply', 'Apply')}
        </button>
      </div>
    </main>
  );
};
