import { useEffect, useState } from 'react';
import { createRoot } from 'react-dom/client';

interface Settings {
  debugMode: boolean;
  debugUrl: string;
}

const DEFAULT_DEBUG_URL = 'http://localhost:8080/library.bundle.js';

const Options = () => {
  const [settings, setSettings] = useState<Settings>({
    debugMode: false,
    debugUrl: DEFAULT_DEBUG_URL,
  });
  const [status, setStatus] = useState('');

  useEffect(() => {
    chrome.storage.local.get(['debugMode', 'debugUrl'], (items) => {
      setSettings({
        debugMode: items.debugMode ?? false,
        debugUrl: items.debugUrl ?? DEFAULT_DEBUG_URL,
      });
    });
  }, []);

  const saveSettings = () => {
    chrome.storage.local.set(settings, () => {
      setStatus('Saved!');
      setTimeout(() => setStatus(''), 2000);
    });
  };

  return (
    <div>
      <h2>News Deframer</h2>
      
      <div style={{ marginBottom: '15px' }}>
        <label style={{ display: 'block', marginBottom: '5px' }}>
          <input
            type="checkbox"
            checked={settings.debugMode}
            onChange={(e) => setSettings({ ...settings, debugMode: e.target.checked })}
          />
          Enable Debug Mode
        </label>
        <small style={{ color: '#666' }}>
          Loads library from remote URL instead of internal bundle.
        </small>
      </div>

      {settings.debugMode && (
        <div style={{ marginBottom: '15px' }}>
          <label style={{ display: 'block', marginBottom: '5px' }}>Debug URL:</label>
          <input
            type="text"
            value={settings.debugUrl}
            onChange={(e) => setSettings({ ...settings, debugUrl: e.target.value })}
            style={{ width: '100%' }}
          />
        </div>
      )}

      <button onClick={saveSettings}>Save Settings</button>
      {status && <span style={{ marginLeft: '10px', color: 'green' }}>{status}</span>}
    </div>
  );
};

const root = createRoot(document.getElementById('root')!);
root.render(<Options />);
