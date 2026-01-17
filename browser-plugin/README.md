# WARNING

This is developer code. Do not use in your regular browser.

It makes you browser insecure!

## Setup

- Create a new Profile in Chrome/Brave/etc
- Install the Tampermonkey <https://chromewebstore.google.com/detail/tampermonkey/dhdgffkkebhmkfjojejmpbldmpobfkfo?hl=en>
- Install the Extensions Reloader <https://chromewebstore.google.com/detail/extensions-reloader/fimgfedafeadlieiabdeeaodndnlbhid>
- Tampermonkey documentation: <https://www.tampermonkey.net/>
- Open the extensions setting in your webbrowser
- Click "settings" in Tampermonkey
- Enable user scripts. **WARNING** understand what this means.

## Tampermonkey Scripts

- Click the Tampermonkey Dashboard Button
- Click "Settings", "General" - enable Advanced Mode (to have access to the storage)
- Click "+" to add a User Script
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

## Why?

- Yes sure! It feels a bit like console hacking 🙃
- During the development cycle we want to edit a JS file on a server that serves the file via https.
- We can use npm + other tools to create a slim and nice library.
- Messing with the Tampermonkey editor doesn't have proper tooling support.
- Again this is a developer hack.
