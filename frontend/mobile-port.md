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

- [done] move shared colors and icons into `frontend/shared/` and add the shared icon-generation make target

## 5. Mobile App Structure And Navigation

- [done] add the mobile app shell with `screens`, `components`, and `services`, including cross-platform storage, a mobile API client, dashboard-first navigation with settings fallback, browser-equivalent settings/session layout, shared theme/i18n behavior, burger-menu navigation, and an About screen with back navigation

## 6. Mobile Portal Screen

- [done] replace the old session flow with a mobile portal screen that opens from the dashboard, keeps the top header/back navigation, loads portal articles from `/mobile/api`, renders native article tiles with loading/error/empty states, and includes sticky `Articles` / `Trend Mining` tabs with a placeholder trend-mining view

## 7. Mobile Article Screen

- replace the placeholder `mobile/src/screens/ArticleScreen.tsx` with a real mobile article view based on `browser-extension/src/ndf/pages/ArticlePage.tsx`
- use the selected article domain as the top app header title instead of `Article`
- keep the existing global mobile header/back behavior and do not add in-page back or hide buttons
- remove the browser-only `view original` button and remove the bottom replacement note
- make the image, corrected title, and corrected description a single clickable top block that opens the original URL in the platform default/system browser
- keep metadata under the upper article block
- show overall rating first, then support expandable detail metrics
- add two bottom toggle buttons: `Original title` and `Details`
- make `Details` toggle the full analysis section (detail metrics + original section)
- make `Original title` toggle only the original title/original text section
- pressing an already-active toggle closes everything again
