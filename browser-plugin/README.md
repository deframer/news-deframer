# News Deframer

Browser extension to deframe news sites.

## Developer Setup

### 1. Installation
1.  Clone repository.
2.  `npm install`
3.  `make dist` (Builds everything and creates `extension.zip`)

### 2. Load in Browser
1.  Open Chrome/Brave -> Extensions (`chrome://extensions`).
2.  Enable **Developer Mode**.
3.  Click **Load unpacked**.
4.  Select the `dist/host` directory (created by the build).
    *   Alternatively, you can drag and drop `extension.zip` if supported, or use it for distribution.

### 3. Recommended Tools
*   **Extensions Reloader:** [Chrome Web Store](https://chromewebstore.google.com/detail/extensions-reloader/fimgfedafeadlieiabdeeaodndnlbhid)
    *   Useful for quickly reloading the extension without going to the extensions page.

### 4. Development Workflow

1.  **Start the Backend:**
    Ensure the Go backend is running (e.g., run `make start` in the project root).

2.  **Start the Extension Builder:**
    ```bash
    npm run dev
    ```
    (This runs Webpack in watch mode to automatically rebuild the extension on file changes. **It does not reload Chrome automatically.**)
3.  Reload the extension in `chrome://extensions` or use the Extensions Reloader tool.

## Project Structure

*   `src/host`: The Extension shell (Manifest, Loader, React Settings).
*   `src/library`: The "Guest" library containing the actual logic.
*   `webpack.config.js`: Builds the extension.
