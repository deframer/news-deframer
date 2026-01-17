console.log("[NDF Background] Background script running");

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
    if (request.type === 'PROXY_REQ') {
        const { url, headers } = request;
        console.log(`[NDF Background] Proxying request to ${url}`);

        fetch(url, { headers })
            .then(response => response.json().then(data => ({ ok: response.ok, status: response.status, data })))
            .then(result => sendResponse(result))
            .catch(err => {
                console.error("[NDF Background] Proxy request failed:", err);
                sendResponse({ ok: false, error: err.toString() });
            });
        return true;
    }
});
