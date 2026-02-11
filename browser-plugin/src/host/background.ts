import log from '../shared/logger';
import { ProxyResponse } from '../shared/types';

log.info('Background script running');

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
  if (request.type === 'SET_BYPASS') {
    sendResponse({ ok: true });
    return true;
  }

  if (request.type === 'PROXY_REQ') {
    const { url, headers, timeout = 5000 } = request;
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
