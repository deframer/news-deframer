// ==UserScript==
// @name        01_Consumer_Script
// @namespace   egandro
// @match       *://example.com/*
// @run-at      document-start
// @grant       unsafeWindow
// ==/UserScript==

(function () {
    'use strict';
    const LOG_PREFIX = "🟩 CONSUMER (Example.com):";

    // Stable Polling: Waits indefinitely for Host
    function waitForLibrary() {
        return new Promise((resolve) => {
            if (typeof unsafeWindow !== "undefined" && unsafeWindow.__lib_ndf) {
                return resolve(unsafeWindow.__lib_ndf());
            }
            console.log(`${LOG_PREFIX} ⏳ Waiting for Host...`);
            const interval = setInterval(() => {
                if (typeof unsafeWindow !== "undefined" && unsafeWindow.__lib_ndf) {
                    clearInterval(interval);
                    resolve(unsafeWindow.__lib_ndf());
                }
            }, 100);
        });
    }

    waitForLibrary().then(async (lib) => {
        const storage = lib.UniversalStorage; // Pre-made Singleton Instance
        console.log(`${LOG_PREFIX} ✅ Connected.`);

        // 1. Subscribe to Sync
        storage.subscribe('shared_test_key', (newVal, oldVal) => {
            console.log(`${LOG_PREFIX} ♻️ Update detected: ${oldVal} -> ${newVal}`);
        });

        // 2. Write Data (Using Host's permissions)
        const ts = new Date().toISOString();
        await storage.set('shared_test_key', `Written by Example.com at ${ts}`);
        console.log(`${LOG_PREFIX} 📝 Wrote: ${ts}`);

    });
})();
