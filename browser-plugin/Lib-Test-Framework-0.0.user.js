// ==UserScript==
// @name        Lib-Test-Framework
// @match       *://*/*
// @run-at      document-start
// @require     https://cdn.jsdelivr.net/gh/egandro/news-deframer@main/browser-plugin/news-deframer-lib.js
// ==/UserScript==

(function () {
    'use strict';

    if (typeof window.__lib_ndf !== "function") {
        console.warn("__lib_ndf is NOT installed yet");
        return;
    }

    const lib = window.__lib_ndf();

    if (!lib || !lib.HelloWorld) {
        console.warn("__lib_ndf does not have a class lib.HelloWorld");
        return;
    }

    const hello = new lib.HelloWorld();
    console.log(hello.message("Lib-Test-Framework"));
})();
