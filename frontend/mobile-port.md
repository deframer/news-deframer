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

- [done] finalize `TrendDetails` layout to use full inner panel real estate (top spacing kept, extra side/bottom inset removed)
- [done] implement mobile `TrendLifecycle` parity with browser behavior: backend lifecycle bars with wide/compact modes, selection-safe refresh handling, and `TrendArticleListPanel` shown below selected bars

## 12. Mobile Trend Context

- [done] implement mobile `TrendContext` equivalent to `browser-extension/src/ndf/components/TrendContext.tsx` with `getContextByDomain` data fetch, loading/empty states, and context chips including frequency

## 13. Mobile Trend Search

- [done] implement mobile `TrendSearch` layout aligned with browser behavior and current mobile Trends spacing: search input + search icon button in the top row
- [done] reuse the same lower details stack as Tag Cloud by rendering `TrendDetailsPanel` (`Lifecycle`, `Context`, `Articles`) for the searched term below the search row, with shared time-range/domain/language/settings wiring

## 14. Mobile Trend Compare Panel

- [done] replace the dummy `TrendComparePanel` with a compare panel wired to `getDomainComparison(domainA, domainB, language, daysInPast)`
- [done] keep the `Compare with` dropdown at the top, using the same compare-domain filtering/naming as the browser version
- [done] split compare results into 3 readable sections: current domain, selected compare domain, and `Shared Trends`
- [done] show only the topic plus article-open buttons in each row; do not show badges, scores, or `A` / `B`
- [done] use the default color for the current-domain action button and a fixed compare color for the selected-domain action button
- [done] use the section headers and action buttons as the main visual indicator so the compare UI stays compact on phone screens

## 15. Compare Panel UX Article List

- [done] keep `mobile/src/screens/ArticleScreen.tsx` untouched and handle compare article UX inside the compare panel only
- [done] tap a row in section A/B or a domain indicator in `Shared Trends` to select `{ term, domain }`, hide the compare list and `Compare with` dropdown, show a selected-term pill, and render the article list below it
- [done] tapping the selected trend pill closes the article list and restores the compare list plus the `Compare with` dropdown
- [done] keep or restore the compare-list scroll position when the compare list is shown again, and clear stale selection when current domain, compare domain, time range, or compare data changes invalidate it

## 16. Mobile Back Button / System Back Handling

- [done] on every mobile screen that shows a visible back button in the header, make Android hardware back do the same thing
- [done] keep the existing iOS header-back behavior aligned with that same screen-level behavior
- [done] do not add extra nested back-state handling beyond what the current header back already does
- [done] use built-in React Native back handling only where needed for Android hardware back
- [done] keep the implementation minimal and screen-based, not a global back-state system
