# Mobile Port Plan

## 1. Hello World App

- [done] a super simple News Deframer React Native app now lives in `frontend/mobile/`, shows `News Deframer` and `Hello world`, and includes local-build make targets for install, Metro, Android, iOS, web, and APK builds

## 2. Revert To The Browser-Extension Rename Point

- revert to the point where `browser-plugin` was renamed to `frontend/browser-extension/`
- do not continue with the `frontend/shared/` architecture for mobile
- React Native cannot use components that render pure HTML or DOM elements such as `div`, `input`, `select`, `button`, and `a`
- browser extension UI and mobile UI should be separate
- mobile should be built independently from the browser extension UI
