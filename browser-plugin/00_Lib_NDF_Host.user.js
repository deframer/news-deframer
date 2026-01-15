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
// @require     https://raw.githubusercontent.com/egandro/news-deframer/refs/heads/main/browser-plugin/news-deframer-lib.js?v=1
// ==/UserScript==

(function () {
    // https://cdn.jsdelivr.net/gh/egandro/news-deframer@main/browser-plugin/news-deframer-lib.js
    'use strict';
    const LOG_PREFIX = "🟦 HOST:";

    if (typeof window.__lib_ndf !== 'function') {
        console.error(`${LOG_PREFIX} ❌ Library load failed. Check @require.`);
        return;
    }

    if (typeof unsafeWindow !== 'undefined') {
        // Bridge the internal Sandbox library to the public Page Scope
        unsafeWindow.__lib_ndf = window.__lib_ndf;
        console.log(`${LOG_PREFIX} ✅ Library bridged to unsafeWindow. Ready for consumers.`);
    } else {
        console.warn(`${LOG_PREFIX} ⚠️ unsafeWindow missing. Consumers cannot connect.`);
    }
})();
