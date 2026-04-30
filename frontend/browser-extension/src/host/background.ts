import log from '../shared/logger';
import { CONNECTION_TIMEOUT_MS, getSettings } from '../shared/settings';
import { ProxyResponse } from '../shared/types';

log.info('Background script running');

const DEFAULT_ICON_PATH = {
  16: 'assets/browser-extension/icon16.png',
  32: 'assets/browser-extension/icon32.png',
  48: 'assets/browser-extension/icon48.png',
  128: 'assets/browser-extension/icon128.png',
} as const;

const ENABLED_ICON_PATH = {
  16: 'assets/browser-extension/icon-enabled16.png',
  32: 'assets/browser-extension/icon-enabled32.png',
  48: 'assets/browser-extension/icon-enabled48.png',
  128: 'assets/browser-extension/icon-enabled128.png',
} as const;

const syncActionIcon = async () => {
  const settings = await getSettings();
  await chrome.action.setIcon({ path: settings.enabled ? ENABLED_ICON_PATH : DEFAULT_ICON_PATH });
};

void syncActionIcon();

chrome.runtime.onInstalled.addListener(() => {
  void syncActionIcon();
});

chrome.runtime.onStartup.addListener(() => {
  void syncActionIcon();
});

chrome.storage.onChanged.addListener((changes, areaName) => {
  if (areaName !== 'local') return;
  if (changes.enabled) {
    void syncActionIcon();
  }
});

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.type === 'SET_BYPASS') {
    sendResponse({ ok: true });
    return true;
  }

  if (request.type === 'PROXY_REQ') {
    const { url, headers, timeout = CONNECTION_TIMEOUT_MS } = request;
    log.info(`Proxying request to ${url} with timeout ${timeout}ms`);

    let responseSent = false;
    const onceSendResponse = (response: ProxyResponse) => {
      if (!responseSent) {
        responseSent = true;
        sendResponse(response);
      }
    };

    const controller = new AbortController();
    const timeoutId = setTimeout(() => {
      controller.abort();
      onceSendResponse({ ok: false, error: `Timeout after ${timeout}ms` });
    }, timeout);

    fetch(url, { headers, signal: controller.signal })
      .then(async (response) => {
        clearTimeout(timeoutId);
        if (response.ok) {
          const data = await response.json();
          onceSendResponse({ ok: true, status: response.status, data });
        } else {
          const errorText = await response.text();
          onceSendResponse({
            ok: false,
            status: response.status,
            error: errorText,
          });
        }
      })
      .catch((err) => {
        clearTimeout(timeoutId);
        if (err.name !== 'AbortError') {
          log.error('Proxy request failed:', err);
          onceSendResponse({ ok: false, error: err.toString() });
        }
      });

    return true; // Indicates an asynchronous response
  }
});
