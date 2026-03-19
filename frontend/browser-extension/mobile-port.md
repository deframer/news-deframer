# Mobile Port Plan

## Major Refactoring Tasks Before The New Architecture

Before introducing a real shared frontend layout, the current browser code needs to be simplified and renamed in a way that matches what the code actually does.

### 1. Replace the current popup-first settings model

- [done] popup is now a small control surface and settings live in a fullscreen browser page under `src/host/pages/`

### 2. Consolidate the current browser code before any mobile work

#### a. Rename `browser-plugin` to `frontend` first

The rename should happen before the internal split.

That gives us one stable top-level name for all later work and avoids planning the rest of the restructure around a directory name we already know we want to replace.

After that rename, the browser-specific build tooling should live with the browser extension host instead of pretending to be shared root tooling.

The immediate target is:

```text
frontend/
  .gitignore
  README.md
  mobile-port.md

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
      manifest.json
      background.ts
      content.ts
      domGuard.ts
      popup/
      options/
      components/
      assets/
      styles.css

  product/
    src/
      index.tsx
      client.ts
      components/
      pages/
      utils/
      styles.css
      styles.ts

  shared/
    src/
      settings.ts
      theme.ts
      i18n.ts
      logger.ts
      domain-cache.ts
      types.ts

  mobile/
    README.md
```

Notes:

- keep `.gitignore` at `frontend/.gitignore`
- do not keep browser-specific `package.json`, webpack, Jest, or TypeScript files at `frontend/` root
- keep the browser extension as a self-contained package under `frontend/browser-extension/`
- move `promo_assets/` into `frontend/browser-extension/promo_assets/` because those assets are browser-extension specific
- mobile gets only a placeholder `README.md` for now

#### b. Separate product screens from browser page manipulation

The content script and page takeover logic should become a thin browser-only mount layer.

That means:

- browser page manipulation stays in the extension host
- article, portal, trends, settings, and dashboard screens move into shared product code
- `src/ndf/index.tsx` should stop being both product bootstrap and browser takeover bootstrap
- `domGuard` stays browser-specific and should move with the extension host, not with shared product code

#### c. Remove the standalone debug browser app

The old standalone debug browser app has been removed.

### 3. Reserve a future mobile host structure, but do not build it yet

We do not want to implement the mobile app structure now.

We only want to document a likely future shape so it can be used later for prompting and planning.

For now, the actual repository change is only:

- create `frontend/mobile/README.md`

Later, if and when mobile work starts, a reasonable target could look like this:

```text
frontend/
  mobile/
    README.md
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

The current `browser-plugin` project should become `frontend` over time.

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
- API repositories and contracts
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

- `src/host/manifest.json`
- `src/host/content.ts`
- `src/host/background.ts`
- `src/host/ui/options.tsx`
- browser mounting flow in `src/ndf/index.tsx`

The current `src/host/ui/options.tsx` should not stay as the long-term settings experience. It should be replaced by:

- a shared fullscreen settings screen that runs in browser and mobile
- a very small extension popup that only starts or stops the browser behavior and opens the fullscreen settings page

### Good candidates for extraction into shared product code

- endpoint models from `src/ndf/client.ts`
- page and feature logic from `src/ndf/pages/*`
- reusable components from `src/ndf/components/*`
- theme tokens from `src/shared/theme.ts`
- translations from `src/shared/i18n.ts`
- pure utilities like rating and URL classification

### Must be rewritten behind interfaces

- `src/ndf/client.ts`
- `src/shared/settings.ts`
- `src/shared/domain-cache.ts`
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
      manifest.json
      background.ts
      content.ts
      domGuard.ts
      popup/
      options/
      components/
      assets/
      styles.css

  product/
    src/
      index.tsx
      client.ts
      components/
      pages/
      utils/
      styles.css
      styles.ts

  shared/
    src/
      settings.ts
      theme.ts
      i18n.ts
      logger.ts
      domain-cache.ts
      types.ts

  mobile/
    README.md
```

This is a simpler transition structure:

- `frontend/` is the renamed top-level directory
- `browser-extension/` contains browser-only build tooling and runtime code
- `browser-extension/` also owns browser-specific promo assets
- `product/` contains the shared product UI and product logic
- `shared/` contains cross-cutting shared support code
- `mobile/` is only a placeholder for now

Later, once the browser refactor is done and mobile work is real, we can decide whether `product/` and `shared/` should become a more formal package layout.

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
- settings screens
- dashboard screens

Browser-only wrappers can still mount those shared screens inside popup, options, dashboard, or content-script containers.

For settings specifically, the preferred direction is:

- move settings UI into a shared fullscreen feature screen
- use the extension popup only as a tiny launcher and start/stop control
- stop treating the popup as the main product UI

## Interfaces

The shared layer should depend on a small set of named adapters.

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

## Refactoring `client.ts`

`src/ndf/client.ts` should become a pure shared repository layer.

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

- rename `browser-plugin` to `frontend`
- keep the popup small and keep settings in a fullscreen browser page
- simplify or remove reload-centric start or stop behavior
- separate browser takeover code from product screens
- remove the standalone debug browser app

### Phase 2

- move browser-specific build files into `frontend/browser-extension/`
- move browser-only runtime code into `frontend/browser-extension/src/`
- move `promo_assets/` into `frontend/browser-extension/promo_assets/`
- move product UI and product logic into `frontend/product/src/`
- move cross-cutting shared code into `frontend/shared/src/`
- create `frontend/mobile/README.md` as a placeholder only

### Phase 3

- replace `src/ndf/client.ts` with a transport-agnostic shared API layer
- move settings and cache behind interfaces
- remove direct platform access from feature code
- define the host-facing interfaces once the boundaries are clearer in the new layout

### Phase 4

- build shared article, portal, trends, settings, and dashboard modules on top of the cleaned structure
- keep browser-specific page takeover as a thin host-only integration
- plan the real mobile app shell and directory structure later, after the browser-side split is stable

## Final Recommendation

Treat News Deframer as one shared product with multiple hosts.

Build:

- a shared source layer for features, state, contracts, and UI
- a thin browser extension host
- a thin mobile host with its own index and menu system
- interface-driven platform services for storage, transport, navigation, sharing, page context, and host rendering

This is the cleanest path for iOS and Android while preserving a single evolving product codebase.
