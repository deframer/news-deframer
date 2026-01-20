import log from '../shared/logger';
import { DEFAULT_BACKEND_URL } from '../shared/settings';

log.info('Background script running');

const updateAuthRules = (backendUrl: string) => {
  try {
    const url = new URL(backendUrl);
    // Construct a pattern that matches the origin (protocol + host + port)
    // We match any path under this origin
    const originPattern = `${url.origin}/*`;

    log.info(`Updating auth rules for ${originPattern}`);

    chrome.declarativeNetRequest.updateDynamicRules({
      removeRuleIds: [1],
      addRules: [
        {
          id: 1,
          priority: 1,
          action: {
            type: chrome.declarativeNetRequest.RuleActionType.MODIFY_HEADERS,
            responseHeaders: [
              {
                header: 'www-authenticate',
                operation: chrome.declarativeNetRequest.HeaderOperation.REMOVE,
              },
            ],
          },
          condition: {
            urlFilter: originPattern,
            resourceTypes: [
              chrome.declarativeNetRequest.ResourceType.MAIN_FRAME,
              chrome.declarativeNetRequest.ResourceType.SUB_FRAME,
              chrome.declarativeNetRequest.ResourceType.STYLESHEET,
              chrome.declarativeNetRequest.ResourceType.SCRIPT,
              chrome.declarativeNetRequest.ResourceType.IMAGE,
              chrome.declarativeNetRequest.ResourceType.FONT,
              chrome.declarativeNetRequest.ResourceType.OBJECT,
              chrome.declarativeNetRequest.ResourceType.XMLHTTPREQUEST,
              chrome.declarativeNetRequest.ResourceType.PING,
              chrome.declarativeNetRequest.ResourceType.CSP_REPORT,
              chrome.declarativeNetRequest.ResourceType.MEDIA,
              chrome.declarativeNetRequest.ResourceType.WEBSOCKET,
              chrome.declarativeNetRequest.ResourceType.OTHER,
            ],
          },
        },
      ],
    });
  } catch (e) {
    log.error('Failed to update auth rules', e);
  }
};

// Initialize rules on startup
chrome.storage.local.get(['backendUrl'], (result) => {
  const backendUrl = result.backendUrl || DEFAULT_BACKEND_URL;
  updateAuthRules(backendUrl);
});

// Listen for settings changes
chrome.storage.onChanged.addListener((changes, namespace) => {
  if (namespace === 'local' && changes.backendUrl) {
    updateAuthRules(changes.backendUrl.newValue);
  }
});

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.type === 'PROXY_REQ') {
    const { url, headers } = request;
    log.info(`Proxying request to ${url}`);

    fetch(url, { headers })
      .then(async (response) => {
        if (response.ok) {
          const data = await response.json();
          return { ok: response.ok, status: response.status, data };
        } else {
          const errorText = await response.text();
          return { ok: response.ok, status: response.status, error: errorText };
        }
      })
      .then((result) => sendResponse(result))
      .catch((err) => {
        log.error('Proxy request failed:', err);
        sendResponse({ ok: false, error: err.toString() });
      });
    return true;
  }
});
