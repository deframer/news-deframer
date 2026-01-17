// ==UserScript==
// @name        02_Readonly_Tester
// @namespace   egandro
// @match       *://www.iana.org/*
// @run-at      document-start
// @grant       unsafeWindow
// ==/UserScript==

(function () {
    'use strict';
    const LOG_PREFIX = "🟧 TESTER (IANA):";

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
        const storage = lib.UniversalStorage;
        console.log(`${LOG_PREFIX} ✅ Connected.`);

        // 1. Read Initial Value
        const val = await storage.get('shared_test_key', 'No Data Yet');
        console.log(`${LOG_PREFIX} 📖 Current Database Value: "${val}"`);

        // 2. Listen for changes coming from example.com
        storage.subscribe('shared_test_key', (newVal) => {
            console.log(`${LOG_PREFIX} 🔔 Notification: Database changed to "${newVal}"`);
        });
    });
})();
