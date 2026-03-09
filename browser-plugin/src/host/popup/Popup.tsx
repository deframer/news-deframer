import '../../shared/i18n';

import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings, Settings } from '../../shared/settings';
import { getThemeCss, globalStyles, Theme } from '../../shared/theme';
import { QuickActions } from '../components/QuickActions';
import { HostStatus,StatusBadge } from '../components/StatusBadge';
import { ToggleSwitch } from '../components/ToggleSwitch';
import { testConnection } from '../lib/connection';

export const Popup = () => {
  const { t } = useTranslation();
  const [settings, setSettings] = useState<Settings | null>(null);
  const [status, setStatus] = useState<HostStatus>('idle');

  useEffect(() => {
    const styleId = 'ndf-theme-styles';
    let style = document.getElementById(styleId);
    if (!style) {
      style = document.createElement('style');
      style.id = styleId;
      document.head.appendChild(style);
    }
    style.textContent = getThemeCss((settings?.theme || 'system') as Theme) + globalStyles;
  }, [settings?.theme]);

  useEffect(() => {
    const load = async () => {
      const loaded = await getSettings();
      setSettings(loaded);
      if (!loaded.enabled) return;
      setStatus('loading');
      try {
        const result = await testConnection(loaded);
        setStatus(result.connected ? 'success' : 'error');
      } catch {
        setStatus('error');
      }
    };

    load();
  }, []);

  const handleEnableToggle = async (enabled: boolean) => {
    if (!settings) return;
    const nextSettings = { ...settings, enabled };
    setSettings(nextSettings);
    await chrome.storage.local.set(nextSettings);

    if (!enabled) {
      setStatus('idle');
      return;
    }

    setStatus('loading');
    try {
      const result = await testConnection(nextSettings);
      setStatus(result.connected ? 'success' : 'error');
    } catch {
      setStatus('error');
    }
  };

  const handleOpenSettings = () => {
    chrome.tabs.create({ url: chrome.runtime.getURL('options.html') });
    window.close();
  };

  if (!settings) {
    return <div className="popup-container">{t('options.loading')}</div>;
  }

  return (
    <main className="popup-container" aria-label={t('options.title', 'Options')}>
      <div className="header compact-header">
        <h2>News Deframer</h2>
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
      <div className="card">
        <ToggleSwitch id="enable-extension" label={t('options.label_enable_extension')} checked={settings.enabled} onChange={handleEnableToggle} />
      </div>
      <QuickActions onOpenSettings={handleOpenSettings} />
    </main>
  );
};
