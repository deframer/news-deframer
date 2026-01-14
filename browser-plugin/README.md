# WARNING

This is developer code. Do not use in your regular browser.

It makes you browser insecure!

## Setup

- Create a new Profile in Chrome/Brave/etc
- Install the Tampermonkey <https://chromewebstore.google.com/detail/tampermonkey/dhdgffkkebhmkfjojejmpbldmpobfkfo?hl=en>
- Tampermonkey documentation: <https://www.tampermonkey.net/>
- Open the extensions setting in your webbrowser
- Click "details" on Tampermonkey
- Enable user scripts. **WARNING** understand what this means.

## Tampermonkey Scripts

- Click on the Tampermonkey Extension Button
- Add new User Script

You are good to go.

## Hello world

File: Save or Ctrl-S

```js
// ==UserScript==
// @name        Hello-World-On-Every-Page
// @match       *://*/*
// @run-at      document-start
// ==/UserScript==

(function () {
    'use strict';

    console.log("Tampermonkey Hello World");
})();
```

- open the console
- open any page
- you see "Tampermonkey Hello World" in the console

