import { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';

interface Settings {
  backendUrl: string;
  username?: string;
  password?: string;
}

const DEFAULT_BACKEND_URL = 'http://localhost:8080';

const Options = () => {
  const [settings, setSettings] = useState<Settings>({
    backendUrl: DEFAULT_BACKEND_URL,
    username: '',
    password: '',
  });
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [loaded, setLoaded] = useState(false);

  // Load settings on mount
  useEffect(() => {
    chrome.storage.local.get(['backendUrl', 'username', 'password'], (items) => {
      const newSettings = {
        backendUrl: (items.backendUrl as string | undefined) ?? DEFAULT_BACKEND_URL,
        username: (items.username as string | undefined) ?? '',
        password: (items.password as string | undefined) ?? '',
      };
      setSettings(newSettings);
      setLoaded(true);
      testConnection(newSettings);
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
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000);

    try {
      const headers: HeadersInit = {};
      if (currentSettings.username && currentSettings.password) {
        headers['Authorization'] = 'Basic ' + btoa(`${currentSettings.username}:${currentSettings.password}`);
      }

      const url = currentSettings.backendUrl.replace(/\/$/, '') + '/api/domains';
      const response = await fetch(url, { headers, signal: controller.signal });

      if (response.ok) {
        setStatus('success');
      } else {
        setStatus('error');
      }
    } catch {
      setStatus('error');
    } finally {
      clearTimeout(timeoutId);
    }
  };

  const handleTestClick = () => {
    testConnection(settings);
  };

  if (!loaded) return <div>Loading...</div>;

  return (
    <div style={{ padding: '20px', maxWidth: '400px' }}>
      <h2>News Deframer Settings</h2>

      <div style={{ marginBottom: '20px', border: '1px solid #ddd', padding: '15px', borderRadius: '4px' }}>
        <h3 style={{ marginTop: 0 }}>Connection</h3>
        <div style={{ marginBottom: '10px' }}>
          <label style={{ display: 'block', marginBottom: '5px' }}>News Deframer URL:</label>
          <input
            type="text"
            value={settings.backendUrl}
            onChange={(e) => setSettings({ ...settings, backendUrl: e.target.value })}
            style={{ width: '100%', padding: '5px' }}
            placeholder="http://localhost:8080"
          />
        </div>

        <div style={{ marginBottom: '10px' }}>
          <label style={{ display: 'block', marginBottom: '5px' }}>Username (Optional):</label>
          <input
            type="text"
            value={settings.username}
            onChange={(e) => setSettings({ ...settings, username: e.target.value })}
            style={{ width: '100%', padding: '5px' }}
          />
        </div>

        <div style={{ marginBottom: '10px' }}>
          <label style={{ display: 'block', marginBottom: '5px' }}>Password (Optional):</label>
          <input
            type="password"
            value={settings.password}
            onChange={(e) => setSettings({ ...settings, password: e.target.value })}
            style={{ width: '100%', padding: '5px' }}
          />
        </div>

        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginTop: '15px' }}>
          <button
            onClick={handleTestClick}
            disabled={status === 'loading'}
            style={{ padding: '5px 10px' }}
          >
            {status === 'loading' ? 'Testing...' : 'Test Connection'}
          </button>

          {status !== 'idle' && (
            <div style={{
              display: 'flex',
              alignItems: 'center',
              color: status === 'success' ? 'green' : 'red',
              fontWeight: 'bold'
            }}>
              <span style={{
                display: 'inline-block',
                width: '10px',
                height: '10px',
                borderRadius: '50%',
                backgroundColor: status === 'success' ? 'green' : 'red',
                marginRight: '5px'
              }}></span>
              {status === 'success' ? 'Connected' : 'Connection Failed'}
            </div>
          )}
        </div>
      </div>

      <div style={{ marginTop: '20px', fontSize: '12px', color: '#999', textAlign: 'center' }}>
        Settings are saved automatically
      </div>
    </div>
  );
};

const root = createRoot(document.getElementById('root')!);
root.render(<Options />);
