# News Deframer - Architecture Plan

## 1. Architecture: The "Hybrid Loader" Pattern

We will split the codebase into two distinct parts managed within this repository.

### A. Host Extension (The Shell)
*   **Role:** Installed in the browser. Handles permissions, storage, and code injection.
*   **Logic:** Minimal "Loader" logic.
    *   Checks `storage` for configuration.
    *   **Debug Mode:** Fetches JS from a user-provided URL (e.g., `localhost`) and executes it using `eval` (Manifest V2).
    *   **Release Mode:** Loads the bundled local version of the library included in the extension package.
*   **UI:** A React-based Popup/Settings page to configure the "Debug URL" and toggle modes.
*   **Tech:** TypeScript, React, Chrome Extension API (Manifest V2).

### B. Guest Library (The Logic)
*   **Role:** The actual DOM modification logic.
*   **Format:** Compiles to a single `library.bundle.js` (UMD/IIFE format).
*   **Entry:** Exposes a standard global entry point (e.g., `window.NewsDeframer.start()`).
*   **Tech:** TypeScript, Webpack.

## 2. Technology Stack
*   **Platform:** Chrome Extension (Manifest V2)
    *   *Reason:* Manifest V2 is required to allow `unsafe-eval` for loading remote code during development.
*   **Language:** TypeScript
*   **UI Framework:** React (for extension popup/settings)
*   **Bundler:** Webpack
    *   `webpack.host.js`: Builds the extension shell.
    *   `webpack.lib.js`: Builds the guest library.

## 3. Build & Execution Modes

| Mode | Description | Loading Mechanism |
| :--- | :--- | :--- |
| **Debug / Dev** | fast iteration, hot-reload-like exp. | Extension fetches `library.bundle.js` from `http://localhost:xxxx`. |
| **Release** | Final distributable | `library.bundle.js` is copied into the extension folder at build time and loaded from disk. |

## 4. Development Phases

### Phase 1: Clean & Init
*   Archive/Remove legacy files.
*   Initialize `package.json`.
*   Install dependencies (React, TypeScript, Webpack).
*   Set up project structure (`src/host`, `src/library`).

### Phase 2: The Guest Library (Core)
*   Set up `webpack.lib.js`.
*   Create a simple "Hello World" DOM modifier.
*   Ensure it exports a clean `init()` function.

### Phase 3: The Host Extension (Shell)
*   Set up `webpack.host.js`.
*   Create `manifest.json` (V2) with permissions: `<all_urls>`, `storage`, `unsafe-eval`.
*   **Settings UI:** React app to input Library URL and toggle Debug Mode.
*   **Content Script:** The "Loader". Reads settings -> injects appropriate script.

### Phase 4: Integration & Polish
*   Implement `npm run build:release` (Copy library to host assets).
*   Implement `npm run dev` (Start dev server for library + watch host).
*   Verify safety and error handling.
