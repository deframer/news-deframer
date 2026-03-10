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

- [done] add the mobile app shell with `screens`, `components`, and `services`, including cross-platform storage, a mobile API client, dashboard-first navigation with settings fallback, browser-equivalent settings/session layout, shared theme/i18n behavior, burger-menu navigation, and an About screen with back navigation

## 6. Mobile Portal Screen

- remove `mobile/src/screens/SessionScreen.tsx` and replace the old session flow with a new `PortalScreen.tsx`
- rename the mobile screen flow from `session` to `portal` in app state, navigation, titles, and related copy
- update the dashboard so selecting a domain opens the new portal screen for that domain
- make the portal screen follow the settings/about navigation pattern: show a `<<` back icon, hide the burger menu, and return to the dashboard on close
- show the selected domain only once at the very top as the portal hero/header title
- use a two-tab layout labeled `Articles` and `Trend Mining`, without showing the domain in the tabs
- make `Articles` the default tab and fetch the portal article list from the backend via `getSite(domain.domain)`
- render the fetched results in a mobile-native article tile layout equivalent to `browser-extension/src/ndf/components/ArticleTile.tsx`
- include loading, empty, and error states for portal article loading
- add a `Trend Mining` placeholder view for now, to reserve the second tab without implementing trend features yet
- keep the mobile portal layout simple and native: top title hero, tabs below it, article list underneath
