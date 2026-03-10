# Mobile Port Plan

## Major Refactoring Tasks Before The New Architecture

Before introducing a real shared frontend layout, the current browser code needs to be simplified and renamed in a way that matches what the code actually does.

### 1. Replace the current popup-first settings model

- [done] popup is now a small control surface and settings live in a fullscreen browser page under the browser host

### 2. Consolidate the current browser code before any mobile work

- [done] browser code is consolidated under `frontend/browser-extension/`, shared browser/mobile primitives live under `frontend/shared/`, and the old debug/browser-plugin leftovers are removed

### 3. Hello World App

- [done] a super simple News Deframer React Native app now lives in `frontend/mobile/`, shows `News Deframer` and `Hello world`, and includes local-build make targets for install, Metro, Android, iOS, web, and APK builds

### 4. Reserve a future mobile host structure, but do not build it yet

We do not want to implement the future mobile app structure yet.

We only want to document a likely future shape so it can be used later for prompting and planning.

After the hello world app exists, a reasonable target could look like this:

```text
frontend/
  mobile/
    README.md
    Makefile
    app/
    src/
      index.tsx
      navigation/
      screens/
      host/
      services/
```

This is intentionally deferred.

## Goal

Build News Deframer as one shared frontend codebase with two hosts:

- a browser extension host
- a mobile app host for iOS and Android

The mobile app gets its own index, menu system, and navigation. The old page-cleaning model is browser-specific host behavior and must not define the shared app architecture.

The old `browser-plugin` project has now been renamed into the `frontend` tree.

## Direction

We are committing to the shared-code architecture only.

That means:

- no separate mobile source tree with duplicated features
- no WebView-first mobile architecture
- no DOM takeover logic in shared feature code
- one shared application layer used by extension, mobile, and future dashboard surfaces

## Core Architecture

Organize the system into three layers:

1. shared product code
2. platform adapters
3. platform hosts

### Shared product code

This layer owns:

- domain models and schemas
- API repository code and shared types
- feature state and query hooks
- shared screens and feature modules
- design system and theme tokens
- auth/session logic
- dashboard, domain selection, and future browser/mobile account-linking flows

This layer must not import:

- `chrome.*`
- `window.*`
- `document.*`
- React DOM APIs
- native mobile APIs

### Platform adapters

This layer implements interfaces for:

- storage
- cache
- HTTP transport
- auth/session persistence
- external navigation
- sharing
- page context
- runtime-specific behaviors

### Platform hosts

Hosts are thin shells.

The browser host owns browser-only behavior like:

- manifest and entrypoints
- content script mounting
- page context detection
- Shadow DOM mounting if still needed
- background proxy wiring
- tiny popup bootstrapping for start/stop and opening the fullscreen settings page
- fullscreen browser settings or dashboard page bootstrapping

The mobile host owns mobile-only behavior like:

- app bootstrap
- app index and menu system
- native navigation
- share sheet integration
- secure storage
- in-app browser or external browser opening

### UI split inside that structure

Some components that currently live under a host directory will need to split between shell wiring and reusable UI.

- keep browser-only wrappers, mount points, and API-bound controls in `browser-extension/src/host/`
- move cross-host visual pieces into `frontend/shared/components/`
- keep feature composition in `browser-extension/src/ndf/` until the mobile host is real
- treat browser-extension settings pages as host shells
- treat shared settings/options content as reusable feature UI
- keep NDF pages separate from host settings pages even when they reuse the same shared components

In practice that means:

- a popup launcher button that directly calls Chrome APIs stays in `src/host/`
- the button, card, form section, status pill, or layout it renders can live in `frontend/shared/components/` or `frontend/shared/settings/`
- a browser settings page becomes a host shell that renders shared settings content from `frontend/shared/settings/`
- a mobile settings screen becomes a mobile host shell that renders the same shared settings content with mobile navigation and storage adapters
- NDF article, portal, trends, and dashboard pages stay under `browser-extension/src/ndf/`

## What Gets Abstracted

The current codebase has several functions that are valid product needs, but their implementation is platform-specific.

Those should be abstracted behind interfaces rather than removed conceptually.

### Platform-specific host behaviors to abstract

- page context lookup
- current URL lookup
- open external URL
- open search URL
- settings persistence
- cache persistence
- network transport
- host mounting
- share-in flow
- share-out flow

### Browser-only host behaviors

These stay browser-specific and should not leak into shared features:

- DOM takeover
- Shadow DOM mounting
- content-script bootstrap
- tab or page scripting
- background-worker proxy transport

### Mobile-only host behaviors

These stay mobile-specific and should not leak into shared features:

- app shell navigation
- native share sheet
- native secure storage
- deep-link handling

### Later and optional

- QR-based linking can be considered later as an optional browser-account pairing flow
- do not design the first shared architecture around QR scanning
- do not implement QR-specific abstractions now

## Current Code Mapping

### Already host-specific

- `browser-extension/src/host/manifest.json`
- `browser-extension/src/host/content.ts`
- `browser-extension/src/host/background.ts`
- `browser-extension/src/host/settings/SettingsPage.tsx`
- browser mounting flow in `browser-extension/src/host/` and `browser-extension/src/ndf/index.tsx`

The current `browser-extension/src/host/settings/SettingsPage.tsx` should not stay as the long-term settings experience. It should be replaced by:

- a browser host page that mounts shared settings content
- a mobile host screen that mounts the same shared settings content inside mobile navigation
- a very small extension popup that only starts or stops the browser behavior and opens the fullscreen settings page

That page should remain separate from `browser-extension/src/ndf/pages/*`.

- `host/settings/` is for browser-extension controls, host settings, and browser-specific wrappers
- `ndf/pages/` is for News Deframer product pages like article, portal, trends, and dashboard
- both can render reusable shared components where that makes sense

### Good candidates for shared extraction

- endpoint models and shared API interfaces from `browser-extension/src/ndf/ndf-api.ts`
- page and feature logic from `browser-extension/src/ndf/pages/*`
- reusable components from `browser-extension/src/ndf/components/*` into `frontend/shared/components/`
- shared settings form pieces into `frontend/shared/settings/`
- logger config into `frontend/shared/logger.ts`
- logging level constants into `frontend/shared/loglevel.ts` and `frontend/shared/loglevel-dev.ts`
- theme tokens into `frontend/shared/theme.ts`
- translations needed by shared components into `frontend/shared/i18n.ts`
- settings types and defaults into `frontend/shared/settings.ts`
- shared API types and interfaces into `frontend/shared/ndf-api-interfaces.ts`
- keep browser-specific client implementation out of `frontend/shared/` until it is actually transport-agnostic
- pure utilities like rating and URL classification

The important distinction is:

- `browser-extension/src/host/` is browser-only
- `browser-extension/src/ndf/` is the current product feature layer
- `frontend/shared/` is shared across browser and mobile

### Must be rewritten behind interfaces

- `browser-extension/src/ndf/ndf-api.ts`
- browser-host settings and storage wiring
- `browser-extension/src/domain-cache.ts`
- direct `window` usage in pages and components
- direct `chrome` usage in shared-looking logic

## Repository Structure

The previous `apps`, `packages`, `features`, `ports`, `providers`, and `services` layout was too abstract for the current codebase.

The structure we should actually move toward first is:

```text
frontend/
  .gitignore
  README.md
  mobile-port.md

  shared/
    components/
    settings/
    assets/
    logger.ts
    loglevel.ts
    loglevel-dev.ts
    theme.ts
    i18n.ts
    settings.ts
    ndf-api-interfaces.ts
    utils/

  browser-extension/
    promo_assets/
    package.json
    package-lock.json
    webpack.config.js
    tsconfig.json
    tsconfig.build.json
    jest.config.js
    jest.setup.js
    eslint.config.mjs
    Makefile
    src/
      host/
        manifest.json
        background.ts
        content.ts
        domGuard.ts
        i18n.ts
        popup/
        settings/
        settings-store.ts
        styles.css
      domain-cache.ts
      ndf/

  mobile/
    README.md
```

This is a simpler transition structure:

- `frontend/` is the renamed top-level directory
- `frontend/shared/` is the shared code layer for browser and mobile
- `frontend/shared/components/` holds reusable UI, including pieces that may currently be rendered from browser host shells
- `frontend/shared/settings/` holds reusable settings and options UI shared by browser and mobile hosts
- `frontend/shared/assets/` holds assets that should be shared across hosts instead of being buried in one host
- `frontend/shared/ndf-api-interfaces.ts` holds shared News Deframer API data types and the interface consumed by shared UI
- `frontend/shared/i18n.ts` holds only the translation resources needed by shared components while each host owns its own runtime initialization and host-specific copy
- `frontend/shared/settings.ts` holds shared settings types/defaults while each host owns persistence adapters
- `frontend/shared/theme.ts`, `frontend/shared/logger.ts`, `frontend/shared/loglevel.ts`, and `frontend/shared/loglevel-dev.ts` hold cross-host presentation and logging primitives
- `browser-extension/` contains browser-only build tooling and runtime code
- `browser-extension/` also owns browser-specific promo assets
- `browser-extension/src/host/` contains browser-only entrypoints, mounting, popup, settings storage, i18n bootstrap, and DOM takeover logic
- `browser-extension/src/ndf/` keeps the current product UI and product logic for now
- `browser-extension/src/domain-cache.ts` stays browser-specific
- `mobile/` is only a placeholder for now

`frontend/shared/` does not need its own `src/` unless it later becomes a real package with its own tooling.

## UI Strategy

If we want to share real components between browser and mobile, we should stop treating the current React DOM plus CSS-string approach as the shared UI base.

Recommended direction:

- build shared UI on React Native primitives
- render it on mobile with React Native
- render it in browser surfaces with a shared web renderer

This gives us real shared components for:

- article screens
- portal screens
- trends screens
- settings and options content
- dashboard screens

Browser-only wrappers can still mount those shared screens inside popup, options, dashboard, or content-script containers.

For settings specifically, the preferred direction is:

- move settings UI into a shared fullscreen feature screen
- use the extension popup only as a tiny launcher and start/stop control
- stop treating the popup as the main product UI

## Interfaces

The shared layer should depend on a small set of named adapters.

For now, `frontend/shared/ndf-api-interfaces.ts` should contain the shared News Deframer DTOs and the `NewsDeframerApi` interface. `frontend/shared/settings.ts`, `frontend/shared/theme.ts`, `frontend/shared/i18n.ts`, and `frontend/shared/logger.ts` can hold cross-host primitives, while lower-level platform adapters stay browser-specific until the mobile host is real.

- `HttpClient`: sends requests to the backend without the shared code needing to know whether the host uses browser proxying, direct fetch, or native networking.
- `SettingsStore`: loads and saves user settings such as backend URL, theme, enabled state, and language.
- `CacheStore`: stores short-lived derived data such as the domain list so the UI does not refetch everything on every screen open. This is an optimization layer, not a source of truth.
- `AuthSessionProvider`: provides auth headers or tokens to the API layer and centralizes login/session handling.
- `ExternalNavigator`: opens external URLs or search URLs in the way the current host expects.
- `PageContextProvider`: provides browser-page context such as current URL or current domain when the extension host needs it. Mobile will use this little or not at all.
- `ShareBridge`: handles incoming shared links and outgoing share actions.
- `HostRenderer`: mounts the shared application into the current host container. In the browser this may mean popup, options, dashboard, or content-script mounting; on mobile this means the app shell and navigation root.

### Platform mappings

Browser extension:

- `SettingsStore` -> `chrome.storage.local`
- `CacheStore` -> `chrome.storage.local`
- `HttpClient` -> background-proxied fetch or direct fetch where allowed
- `ExternalNavigator` -> tabs API or `window.open`
- `PageContextProvider` -> active tab and content-page context
- `HostRenderer` -> popup root, options root, dashboard root, or content-script mount

Mobile:

- `SettingsStore` -> AsyncStorage or MMKV
- `CacheStore` -> AsyncStorage, MMKV, or SQLite
- `HttpClient` -> native fetch client
- `ExternalNavigator` -> Linking or in-app browser
- `ShareBridge` -> share sheet integration
- `HostRenderer` -> app shell, navigation tree, and screen container

## Refactoring `ndf-api.ts`

`browser-extension/src/ndf/ndf-api.ts` should implement a shared repository interface that shared UI can depend on.

It should no longer know about:

- `chrome.runtime.sendMessage`
- extension proxy messages
- browser-specific transport rules

Instead it should receive:

- base URL
- `HttpClient`
- `AuthSessionProvider`

Suggested shape:

```ts
export class NewsDeframerApi {
  constructor(
    private baseUrl: string,
    private http: HttpClient,
    private auth: AuthSessionProvider,
  ) {}

  async getDomains(): Promise<DomainEntry[]> {
    const headers = await this.auth.getAuthHeaders();
    const response = await this.http.request<DomainEntry[]>({
      method: 'GET',
      url: `${this.baseUrl}/api/domains`,
      headers,
    });
    return response.data;
  }
}
```

The same shared API client can then be used by:

- extension popup
- extension dashboard
- content-script-mounted shared screens
- mobile app
- future web dashboard

## Shared Feature Boundaries

Shared features should own:

- article screen
- portal screen
- trends screen
- trend details
- settings screen
- dashboard screen
- auth and pairing flows

Shared features should not directly do:

- `getSettings()` from Chrome storage
- `window.open(...)`
- `window.location...`
- `sessionStorage...`
- content-script mounting
- DOM cleanup

Those actions must come from injected services.

## Future-Proofing

This structure supports future work without splitting the codebase:

- browser dashboard for domain selection
- mobile share-sheet article import
- synced preferences and account state
- optional later account-linking flows if we decide to add them
- additional app surfaces using the same feature packages

## Migration Plan

### Phase 1

- [done] rename `browser-plugin` to `frontend`
- keep the popup small and keep settings in a fullscreen browser page
- remove the standalone debug browser app
- move all browser-specific tooling into `frontend/browser-extension/`
- keep the current runtime code together under `frontend/browser-extension/src/`

### Phase 2

- move `promo_assets/` into `frontend/browser-extension/promo_assets/`
- restore and keep `src/host` as the browser-only code location
- keep `src/ndf` as the current product code location for now
- create `frontend/shared/` for shared browser/mobile UI, `logger.ts`, `loglevel.ts`, `loglevel-dev.ts`, `theme.ts`, `i18n.ts`, `settings.ts`, `ndf-api-interfaces.ts`, and pure utilities
- keep `domain-cache.ts` browser-extension specific
- create `frontend/mobile/README.md` as a placeholder only

### Phase 3

- replace browser-only dependencies inside `browser-extension/src/ndf/ndf-api.ts` with transport-agnostic adapters over time
- move settings and cache behind interfaces
- remove direct platform access from feature code
- move pure reusable pieces from browser-specific code into `frontend/shared/`
- split mixed host components so browser wiring stays in `src/host/` while reusable UI moves into `frontend/shared/components/`
- define the host-facing interfaces once the boundaries are clearer in the simplified layout

### Phase 4

- build shared article, portal, trends, settings, and dashboard modules on top of the cleaned structure
- keep browser-specific page takeover as a thin host-only integration
- plan the real mobile app shell and directory structure later, after the browser-side split is stable

## Final Recommendation

Treat News Deframer as one shared product with multiple hosts.

Build:

- a shared source layer for features, state, interfaces, API code, and UI
- a thin browser extension host
- a thin mobile host with its own index and menu system
- interface-driven platform services for storage, transport, navigation, sharing, page context, and host rendering

This is the cleanest path for iOS and Android while preserving a single evolving product codebase.
