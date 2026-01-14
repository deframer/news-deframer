// ==UserScript==
// @name        Lib-Test-Framework
// @match       *://*/*
// @run-at      document-start
// @require     https://raw.githubusercontent.com/egandro/news-deframer/refs/heads/main/browser-plugin/news-deframer-lib.js
// ==/UserScript==

(function () {
    'use strict';

    if (typeof window._lib_test === "function") {
        const lib = window._lib_test();

        if (lib && lib.HelloWorld) {
            const hello = new lib.HelloWorld();
            console.log(hello.message("Lib-Test-Framework"));
        }
    } else {
        console.warn("_lib_test is NOT installed yet");
    }
})();
