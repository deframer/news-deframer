import { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';

import { invalidateDomainCache } from '../../shared/domain-cache';
import {
  DEFAULT_BACKEND_URL,
  getSettings,
  Settings,
} from '../../shared/settings';

type Status = 'idle' | 'loading' | 'success' | 'error';

const Options = () => {
  const [settings, setSettings] = useState<Settings>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
    enabled: true,
  });
  const [status, setStatus] = useState<Status>('idle');
  const [loaded, setLoaded] = useState(false);

  // Load settings on mount
  useEffect(() => {
    getSettings().then((loadedSettings) => {
      setSettings(loadedSettings);
      setLoaded(true);
      if (loadedSettings.enabled) {
        testConnection(loadedSettings);
      }
    });
  }, []);

  // Autosave when settings change (skip initial load)
  useEffect(() => {
    if (!loaded) return;

    const handler = setTimeout(() => {
      chrome.storage.local.set(settings);
    }, 500);

    return () => clearTimeout(handler);
  }, [settings, loaded]);

  const testConnection = async (currentSettings: Settings) => {
    setStatus('loading');
    try {
      const headers: HeadersInit = {};
      if (currentSettings.username && currentSettings.password) {
        headers['Authorization'] =
          'Basic ' +
          btoa(`${currentSettings.username}:${currentSettings.password}`);
      }

      const url =
        currentSettings.backendUrl.replace(/\/$/, '') + '/api/domains';

      const response = await new Promise<any>((resolve, reject) => {
        chrome.runtime.sendMessage(
          { type: 'PROXY_REQ', url, headers, timeout: 5000 },
          (res) => {
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
        return true;
      } else {
        setStatus('error');
        return false;
      }
    } catch (e) {
      setStatus('error');
      return false;
    }
  };

  const handleTestClick = () => {
    if (status === 'loading') {
      return; // Ignore click if test is already in progress
    }
    testConnection(settings);
  };

  if (!loaded) return <div style={{ padding: '20px' }}>Loading...</div>;

  const isConnected = status === 'success';
  const isError = status === 'error';
  const isLoading = status === 'loading';

  return (
    <div
      style={{
        padding: '24px',
        maxWidth: '400px',
        fontFamily:
          '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
        color: '#333',
      }}
    >
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '24px',
          borderBottom: '1px solid #eee',
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
                ? '#10B981'
                : isError
                ? '#EF4444'
                : '#6B7280',
              backgroundColor: isConnected
                ? '#D1FAE5'
                : isError
                ? '#FEE2E2'
                : '#F3F4F6',
              padding: '4px 8px',
              borderRadius: '12px',
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

      <div
        style={{
          backgroundColor: '#f9fafb',
          border: '1px solid #e5e7eb',
          borderRadius: '8px',
          padding: '16px',
          marginBottom: '24px',
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
            color: '#6b7280',
          }}
        >
          Connection Settings
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
              padding: '8px 12px',
              borderRadius: '6px',
              border: '1px solid #d1d5db',
              boxSizing: 'border-box',
              fontSize: '14px',
              backgroundColor: settings.enabled ? '#fff' : '#f3f4f6',
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
            <span style={{ fontWeight: 400, color: '#9ca3af' }}>
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
              padding: '8px 12px',
              borderRadius: '6px',
              border: '1px solid #d1d5db',
              boxSizing: 'border-box',
              fontSize: '14px',
              backgroundColor: settings.enabled ? '#fff' : '#f3f4f6',
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
            <span style={{ fontWeight: 400, color: '#9ca3af' }}>
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
              padding: '8px 12px',
              borderRadius: '6px',
              border: '1px solid #d1d5db',
              boxSizing: 'border-box',
              fontSize: '14px',
              backgroundColor: settings.enabled ? '#fff' : '#f3f4f6',
            }}
          />
        </div>

        <button
          onClick={handleTestClick}
          disabled={!settings.enabled || isLoading}
          style={{
            width: '100%',
            padding: '10px',
            backgroundColor: settings.enabled ? '#fff' : '#f3f4f6',
            border: '1px solid #d1d5db',
            borderRadius: '6px',
            cursor: !settings.enabled || isLoading ? 'not-allowed' : 'pointer',
            fontSize: '14px',
            fontWeight: 500,
            color: settings.enabled ? '#374151' : '#9ca3af',
            transition: 'background-color 0.2s, color 0.2s',
          }}
        >
          {isLoading ? 'Testing...' : 'Test Connection'}
        </button>
      </div>

      {/* Main Enable Toggle */}
      <div style={{ marginBottom: '24px' }}>
        <label
          style={{
            display: 'flex',
            alignItems: 'center',
            cursor: 'pointer',
            fontSize: '16px',
            fontWeight: 500,
          }}
        >
          <input
            type="checkbox"
            checked={settings.enabled}
            onChange={async (e) => {
              const enabled = e.target.checked;
              const newSettings = { ...settings, enabled };
              setSettings(newSettings);

              // Persist immediately so the reloaded tab sees the new state
              await chrome.storage.local.set(newSettings);

              if (enabled) {
                const success = await testConnection(newSettings);
                if (success) {
                  // On successful enable, reload tab and close popup
                  const [tab] = await chrome.tabs.query({
                    active: true,
                    currentWindow: true,
                  });
                  if (tab?.id && tab.url?.match(/^https?:\/\//)) {
                    chrome.tabs.reload(tab.id);
                  }
                  window.close();
                }
                // If not successful, do nothing (window stays open, no reload)
              } else {
                // On disable, invalidate cache, reload tab, and close popup
                setStatus('idle');
                await invalidateDomainCache();
                const [tab] = await chrome.tabs.query({
                  active: true,
                  currentWindow: true,
                });
                if (tab?.id && tab.url?.match(/^https?:\/\//)) {
                  chrome.tabs.reload(tab.id);
                }
                window.close();
              }
            }}
            style={{ marginRight: '10px', width: '18px', height: '18px' }}
          />
          Enable Extension
        </label>
      </div>
    </div>
  );
};

const root = createRoot(document.getElementById('root')!);
root.render(<Options />);
