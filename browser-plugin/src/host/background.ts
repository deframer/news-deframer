import log from '../shared/logger';
import { DEFAULT_BACKEND_URL } from '../shared/settings';

log.info('Background script running');

const updateAuthRules = (backendUrl: string) => {
  try {
    if (typeof backendUrl !== 'string' || !backendUrl.startsWith('http')) {
      log.error('Invalid backendUrl for auth rules:', backendUrl);
      return;
    }
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
  if (typeof backendUrl === 'string') {
    updateAuthRules(backendUrl);
  }
});

// Listen for settings changes
chrome.storage.onChanged.addListener((changes, namespace) => {
  if (namespace === 'local' && changes.backendUrl) {
    const newUrl = changes.backendUrl.newValue;
    if (typeof newUrl === 'string') {
      updateAuthRules(newUrl);
    }
  }
});

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.type === 'PROXY_REQ') {
    const { url, headers, timeout = 5000 } = request;
    log.info(`Proxying request to ${url} with timeout ${timeout}ms`);

    let responseSent = false;
    const onceSendResponse = (response: any) => {
      if (!responseSent) {
        responseSent = true;
        sendResponse(response);
      }
    };

    const controller = new AbortController();
    const timeoutId = setTimeout(() => {
      controller.abort();
      onceSendResponse({ ok: false, error: 'Connection timed out' });
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
