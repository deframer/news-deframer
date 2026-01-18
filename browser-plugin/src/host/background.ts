import log from '../shared/logger';

log.info('Background script running');

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
