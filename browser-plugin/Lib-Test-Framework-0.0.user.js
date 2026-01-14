// ==UserScript==
// @name        Lib-Test-Framework
// @match       *://*/*
// @run-at      document-start
// @require     https://github.com/egandro/news-deframer/browser-plugin/lib-test.js
// ==/UserScript==

(function () {
    'use strict';

    if (typeof window._lib_test === "function") {
        const lib = window._lib_test();

        if (lib && typeof lib.helloWorld === "function") {
            lib.helloWorld();
        }
    } else {
        console.warn("_lib_test is NOT installed yet");
    }
})();
