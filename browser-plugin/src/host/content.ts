export {};

// Define the interface for our library
interface Config {
  backendUrl: string;
  username?: string;
  password?: string;
}

interface NewsDeframerLib {
  start: (config: Config) => void;
}

declare global {
  interface Window {
    NewsDeframer?: NewsDeframerLib;
  }
}

const DEFAULT_DEBUG_URL = 'http://localhost:8080/library.bundle.js';
const DEFAULT_BACKEND_URL = 'http://localhost:8080';

async function injectScript(url: string, isRemote: boolean, config: Config) {
    if (isRemote) {
        // Fetch via background script to avoid CORS/Mixed-Content issues
        console.log(`[NDF Host] Requesting background fetch for ${url}...`);
        chrome.runtime.sendMessage({ type: 'FETCH_DEBUG_LIB', url }, (response) => {
            if (chrome.runtime.lastError) {
                console.error("[NDF Host] Messaging error:", chrome.runtime.lastError);
                return;
            }
            
            if (response && response.success) {
                console.log("[NDF Host] Received library from background. Executing...");
                try {
                    (0, eval)(response.code);
                    if (window.NewsDeframer) {
                        console.log("[NDF Host] Library loaded. Starting...");
                        window.NewsDeframer.start(config);
                    } else {
                        console.error("[NDF Host] Library loaded but window.NewsDeframer is missing.");
                    }
                } catch (e) {
                    console.error("[NDF Host] Eval failed:", e);
                }
            } else {
                 console.error("[NDF Host] Background fetch failed:", response?.error || "Unknown error");
            }
        });
    } else {
        // Local release mode - inject script tag
        try {
            console.log(`[NDF Host] Injecting local library from ${url}...`);
            const response = await fetch(url);
            if (!response.ok) throw new Error(`HTTP ${response.status}`);
            const code = await response.text();
            
            console.log("[NDF Host] Executing local library...");
            (0, eval)(code);

            if (window.NewsDeframer) {
                console.log("[NDF Host] Library loaded. Starting...");
                window.NewsDeframer.start(config);
            }
        } catch (err) {
            console.error(`[NDF Host] Failed to load local library:`, err);
        }
    }
}

function loadLibrary() {
  chrome.storage.local.get(['debugMode', 'debugUrl', 'backendUrl', 'username', 'password'], (items) => {
    const debugMode = items.debugMode ?? false;
    const config: Config = {
        backendUrl: items.backendUrl ?? DEFAULT_BACKEND_URL,
        username: items.username,
        password: items.password
    };
    
    if (debugMode) {
      const debugUrl = items.debugUrl ?? DEFAULT_DEBUG_URL;
      injectScript(debugUrl, true, config);
    } else {
      const localUrl = chrome.runtime.getURL('assets/library.bundle.js');
      injectScript(localUrl, false, config);
    }
  });
}

loadLibrary();
