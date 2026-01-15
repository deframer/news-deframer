### Architecture Overview
1.  **Library (UniversalStorage):** A TypeScript class that unifies access to storage (Chrome Extension, Tampermonkey, or LocalStorage).
2.  **Singleton Pattern:** The library instantiates the class **internally** when loaded. This captures the Host's permissions immediately.
3.  **Host Script:** Loads the library, holds the `@grant` permissions, and bridges the library to the page (`unsafeWindow`).
4.  **Consumer Script:** Polls the page for the bridged library and uses the Host's storage instance to read/write data.

---

### 1. The Core Class (`src/UniversalStorage.ts`)
This class handles the logic. It uses the modern asynchronous `GM.*` API (Tampermonkey 5+) and falls back to a memory object if `localStorage` fails.

```typescript
// src/UniversalStorage.ts

declare global {
    // 1. Browser Extension API
    const chrome: any;

    // 2. Tampermonkey 5+ API (GM.*)
    const GM: {
        setValue: (key: string, value: any) => Promise<void>;
        getValue: <T>(key: string, defaultValue?: T) => Promise<T>;
        deleteValue: (key: string) => Promise<void>;
        listValues: () => Promise<string[]>;
        addValueChangeListener?: (
            key: string,
            callback: (name: string, oldVal: any, newVal: any, remote: boolean) => void
        ) => Promise<number>;
    };
}

export class UniversalStorage {

    /**
     * Internal global variable for memory fallback when LocalStorage is unavailable.
     */
    private static memoryStorage: { [key: string]: any } = {};

    /**
     * Set a value.
     * Priority: Chrome Storage -> Tampermonkey 5+ (GM.setValue) -> LocalStorage
     */
    async set(key: string, value: any): Promise<void> {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => chrome.storage.local.set({ [key]: value }, resolve));
        }

        // 2. Tampermonkey 5+
        if (typeof GM !== 'undefined' && GM.setValue) {
            return await GM.setValue(key, value);
        }

        // 3. Web Fallback
        try {
            localStorage.setItem(key, JSON.stringify(value));
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable, using memory fallback.", e);
            UniversalStorage.memoryStorage[key] = value;
        }
        return Promise.resolve();
    }

    /**
     * Get a value.
     */
    async get<T>(key: string, defaultValue: T): Promise<T> {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => {
                chrome.storage.local.get([key], (res: any) => {
                    resolve(res[key] !== undefined ? res[key] : defaultValue);
                });
            });
        }

        // 2. Tampermonkey 5+
        if (typeof GM !== 'undefined' && GM.getValue) {
            return await GM.getValue(key, defaultValue);
        }

        // 3. Web Fallback
        try {
            const item = localStorage.getItem(key);
            return Promise.resolve(item ? JSON.parse(item) : defaultValue);
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable, checking memory fallback.", e);
            const memVal = UniversalStorage.memoryStorage[key];
            return Promise.resolve(memVal !== undefined ? memVal : defaultValue);
        }
    }

    /**
     * Remove a value.
     */
    async remove(key: string): Promise<void> {
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.local) {
            return new Promise((resolve) => chrome.storage.local.remove(key, resolve));
        }
        if (typeof GM !== 'undefined' && GM.deleteValue) {
            return await GM.deleteValue(key);
        }
        try {
            localStorage.removeItem(key);
        } catch (e) {
            console.warn("UniversalStorage: LocalStorage unavailable (remove), clearing memory fallback.", e);
        }
        delete UniversalStorage.memoryStorage[key];
        return Promise.resolve();
    }

    /**
     * Sync: Subscribe to changes from other Tabs / Scripts.
     */
    subscribe(key: string, callback: (newValue: any, oldValue: any) => void): void {
        // 1. Browser Extension
        if (typeof chrome !== 'undefined' && chrome.storage && chrome.storage.onChanged) {
            chrome.storage.onChanged.addListener((changes: any, namespace: string) => {
                if (namespace === 'local' && changes[key]) {
                    callback(changes[key].newValue, changes[key].oldValue);
                }
            });
            return;
        }

        // 2. Tampermonkey 5+
        if (typeof GM !== 'undefined' && GM.addValueChangeListener) {
            GM.addValueChangeListener(key, (name: string, oldVal: any, newVal: any, remote: boolean) => {
                if (remote) callback(newVal, oldVal);
            });
            return;
        }

        // 3. Web Fallback
        window.addEventListener('storage', (event) => {
            if (event.key === key) {
                try {
                    const n = event.newValue ? JSON.parse(event.newValue) : null;
                    const o = event.oldValue ? JSON.parse(event.oldValue) : null;
                    callback(n, o);
                } catch { callback(event.newValue, event.oldValue); }
            }
        });
    }
}
```

---

### 2. The Library Entry (`src/index.ts`)
This is crucial. We instantiate the **Singleton** here. When the Host Script loads this file, this code runs in the Host's context (with permissions). We export the *instance*, not the class.

```typescript
// src/index.ts
import { HelloWorld } from './helloworld'; // Optional
import { UniversalStorage } from './UniversalStorage';

console.log("📦 Library: Initializing singleton...");

// 1. Instantiate the Singleton HERE.
// This runs once when the library loads, capturing the Host's permissions.
const storageInstance = new UniversalStorage();

// 2. Export the Function as Default
// Webpack will assign the return value of this function to 'window.__lib_ndf'
export default function () {
    return {
        // Return the Class (if you ever need to make new Hellos)
        HelloWorld,

        // Return the PRE-MADE Instance (The Singleton)
        UniversalStorage: storageInstance,
    };
}
```

---

### 3. Build Configuration (`webpack.config.js`)
Ensures the library is attached to `window.__lib_ndf`.

```javascript
const path = require('path');

module.exports = {
    entry: './src/index.ts',
    mode: 'development', // 'production' for minified
    devtool: 'inline-source-map',
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.ts', '.js'],
    },
    output: {
        filename: 'news-deframer-lib.js',
        path: path.resolve(__dirname, 'dist'),
        library: {
            name: '__lib_ndf',  // The global variable name
            type: 'window',     // Attach to window
            export: 'default',  // Export the default function from index.ts
        },
        clean: true,
    },
};
```

---

### 4. The "Host" Script (Provider)
**Role:** Holds permissions (`GM.setValue`, etc.) and bridges the library to `unsafeWindow` so Consumers can see it.
**Important:** Update the `v=` parameter in `@require` whenever you push a new build to avoid CDN caching issues.

```javascript
// ==UserScript==
// @name        00_Lib_NDF_Host
// @namespace   egandro
// @match       *://*/*
// @run-at      document-start
// @grant       GM.setValue
// @grant       GM.getValue
// @grant       GM.deleteValue
// @grant       GM.listValues
// @grant       GM.addValueChangeListener
// @grant       unsafeWindow
// @require     https://cdn.jsdelivr.net/gh/egandro/news-deframer@browser-plugin-universal-storage/browser-plugin/news-deframer-lib.js?v=FIX_CACHE_1
// ==/UserScript==

(function () {
    'use strict';
    const LOG_PREFIX = "🟦 HOST:";

    // 1. Verify Library Loaded (in Sandbox)
    if (typeof window.__lib_ndf !== 'function') {
        console.error(`${LOG_PREFIX} ❌ Library load failed. Check @require.`);
        return;
    }

    // 2. Bridge to Page Scope
    if (typeof unsafeWindow !== 'undefined') {
        // Expose the library factory to the page
        unsafeWindow.__lib_ndf = window.__lib_ndf;
        console.log(`${LOG_PREFIX} ✅ Library bridged to unsafeWindow. Ready for consumers.`);
    } else {
        console.warn(`${LOG_PREFIX} ⚠️ unsafeWindow missing. Consumers cannot connect.`);
    }
})();
```

---

### 5. The "Consumer" Script (Client)
**Role:** Polls for the library and uses the shared storage.
**Important:** It must use `// @grant unsafeWindow` to see the bridged library.

```javascript
// ==UserScript==
// @name        01_Consumer_Script
// @namespace   egandro
// @match       *://*/*
// @run-at      document-start
// @grant       unsafeWindow
// ==/UserScript==

(function () {
    'use strict';
    const LOG_PREFIX = "🟩 CONSUMER:";

    // Stable Polling: Waits indefinitely for Host
    function waitForLibrary() {
        return new Promise((resolve) => {
            // Check immediately
            if (typeof unsafeWindow !== "undefined" && unsafeWindow.__lib_ndf) {
                return resolve(unsafeWindow.__lib_ndf());
            }

            console.log(`${LOG_PREFIX} ⏳ Waiting for Host...`);

            // Poll every 100ms
            const interval = setInterval(() => {
                if (typeof unsafeWindow !== "undefined" && unsafeWindow.__lib_ndf) {
                    clearInterval(interval);
                    resolve(unsafeWindow.__lib_ndf());
                }
            }, 100);
        });
    }

    // Main Logic
    waitForLibrary().then(async (lib) => {
        const storage = lib.UniversalStorage; // Pre-made Singleton Instance
        console.log(`${LOG_PREFIX} ✅ Connected.`);

        // 1. Subscribe to Sync
        storage.subscribe('shared_demo_key', (newVal, oldVal) => {
            console.log(`${LOG_PREFIX} ♻️ Update detected: ${oldVal} -> ${newVal}`);
        });

        // 2. Write Data (Using Host's permissions)
        const ts = new Date().toLocaleTimeString();
        await storage.set('shared_demo_key', `Updated by Consumer at ${ts}`);
        console.log(`${LOG_PREFIX} 📝 Wrote: ${ts}`);

        // 3. Read Data
        const val = await storage.get('shared_demo_key');
        console.log(`${LOG_PREFIX} 📖 Read: ${val}`);

    });
})();
```