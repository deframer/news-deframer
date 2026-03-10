# Mobile Port Plan

## 1. Hello World App

- [done] a super simple News Deframer React Native app now lives in `frontend/mobile/`, shows `News Deframer` and `Hello world`, and includes local-build make targets for install, Metro, Android, iOS, web, and APK builds

## 2. Mobile/Extension Separation

- [done] mobile and browser extension remain separate UI implementations
- [done] React Native does not reuse browser extension components that render pure HTML or DOM elements such as `div`, `input`, `select`, `button`, and `a`
- [done] mobile is built independently from the browser extension UI

## 3. Shared Wordcloud

- [done] create our own shared wordcloud in `frontend/shared/wordcloud/`
- [done] keep the layout logic shared and render it separately for browser and mobile

## 4. Shared Colors

- create `frontend/shared/` for shared colors
- move the browser extension theme colors there
- move the icons there too
- add a `frontend/shared/Makefile` to generate the shared icons

## 5. Mobile App Structure And Navigation

- build the mobile app with `screens`, `components`, and `services`
- add a cross-platform storage service for iOS, Android, and web via `@react-native-async-storage/async-storage`
- add a mobile API client service equivalent to `browser-extension/src/ndf/client.ts`
- start on a Dashboard screen that lists domains; this screen does not exist in the browser extension yet
- if required configuration is missing, open Settings first instead of Dashboard
- keep the session page layout equivalent to the browser version
- keep the settings layout equivalent to the browser settings cards and form sections
- support the same `light`, `dark`, and `system` theme modes as the browser version
- support i18n on mobile too and reuse the same texts and translation keys as the browser version
- add a Material-style burger menu with entries for Dashboard, Settings, and About
- remove mobile-unneeded browser controls from Settings: no enable/disable and no Apply button
- in About, show `News Deframer` and the open-project link; the About layout can differ from the current browser options page
- add a `<<` back button on Settings and About to return to Dashboard
