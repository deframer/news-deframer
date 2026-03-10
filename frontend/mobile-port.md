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

- [done] replace the placeholder article screen with a full mobile article view that uses the article domain in the top header, keeps global back/header behavior, adds a clickable top content block with visible open action to launch the original URL in the system browser, shows metadata and ratings with expandable details/original sections, and uses sticky bottom `Original title` / `Details` toggles with close-on-active and switch-on-other behavior

## 8. A11Y Font Size

- add an accessibility font-size check pass across all mobile screens
- evaluate dynamic type/text scaling for key text styles (titles, body, metadata, buttons, ratings)
- add a font-size setting in mobile settings so users can increase/decrease text size consistently

## 9. Mobile Trends

- [done] rename the main tab to `Trends` and implement the Trends sub-layout in portal with a time-selector row (`24h`, `7d`, `30d`, `90d`, `365d`, default `7d`), subtabs (`Tag Cloud`, `Compare`, `Search`) as mobile components, and placeholder panels only for this first iteration (no footer, `Articles` unchanged)

## 10. Mobile Trend Tag Cloud + Trend Details

- [done] port browser-like trend flow to mobile: interactive `TrendTagCloud` with selected-term drill-down into `TrendDetails` tabs (`Lifecycle`, `Context`, `Articles`)
- [done] add mobile trend components and details-tab behavior (`TrendTagCloudPanel`, `TrendDetailsPanel`, `TrendLifecyclePanel`, `TrendContextPanel`, `TrendArticleListPanel`) including `Articles` panel placement

## 11. Trend Details + Lifecycle

- make `TrendDetails` panel content use full inner real estate of the container: keep existing top spacing, but remove extra inner left/right/bottom spacing from lifecycle/context/articles panels
- implement `TrendLifecycle` in mobile based on `browser-extension/src/ndf/components/TrendLifecycleChart.tsx` using mobile components and backend lifecycle data
- add lifecycle chart display modes for mobile: default wide mode with minimum bar width and horizontal scroll, plus compact fit mode toggled by a chart-header icon button (top-left)
- keep lifecycle bars selectable in wide mode; compact mode can be non-clickable when needed for readability
- show selected lifecycle details below the chart in a compact separated info area (date, frequency, velocity), similar to bottom metadata/info treatment in `mobile/src/components/ArticleTile.tsx`
- display `TrendArticleListPanel` below the lifecycle chart when a date/bar is selected
- when date range/domain/term changes, keep selection only if selected date is still present; otherwise clear selection safely
