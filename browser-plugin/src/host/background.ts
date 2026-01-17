console.log("[NDF Background] Background script running");

chrome.runtime.onMessage.addListener((request, _sender, sendResponse) => {
    if (request.type === 'FETCH_DEBUG_LIB') {
        const url = request.url;
        console.log(`[NDF Background] Fetching ${url}...`);

        fetch(url)
            .then(response => {
                if (!response.ok) throw new Error(`HTTP ${response.status}`);
                return response.text();
            })
            .then(code => {
                sendResponse({ success: true, code });
            })
            .catch(err => {
                console.error("[NDF Background] Fetch failed:", err);
                sendResponse({ success: false, error: err.toString() });
            });

        return true; // Keep channel open for async response
    }

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
