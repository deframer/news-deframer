# News Deframer Browser Extension

Browser extension host package for News Deframer.

## Developer Setup

### 1. Installation
1. Clone repository.
2. `cd frontend/browser-extension`
3. `npm install`
4. `cd ../shared && make icons`
5. `cd ../browser-extension && make dist`

### 2. Load in Browser
1. Open Chrome/Brave -> Extensions (`chrome://extensions`).
2. Enable Developer Mode.
3. Click Load unpacked.
4. Select `frontend/browser-extension/dist/host`.

### 3. Development Workflow
1. Start the backend from the project root if needed.
2. Run:

```bash
cd frontend/browser-extension && npm run dev
```

3. Reload the extension in the browser after rebuilds.

## Project Structure

- `src/host/`: browser-only entrypoints, popup, settings page shell, content script, background script, DOM/browser wiring
- `src/ndf/`: News Deframer product pages and feature flow that still live in the browser package
- `src/shared/`: browser-local support code such as settings, theme, i18n, logger, and browser-only types
- `../shared/components/`: reusable NDF components shared across host and product surfaces
- `../shared/settings/`: reusable settings UI rendered inside the browser settings shell
- `../shared/assets/`: shared assets copied into the built extension

## Common Commands

```bash
npm run dev
npm run build
npm run lint
npm run test
```

Build output is written to `dist/host/`.

Refresh shared browser extension icons with `cd ../shared && make icons` when `../shared/assets/icons/icon.svg` changes.
