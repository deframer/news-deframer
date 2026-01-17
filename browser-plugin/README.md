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

#### Debug Mode (Recommended for Library Dev)
This allows you to edit the library code and see changes on reload without rebuilding the extension.

1.  Start the development servers:
    ```bash
    npm run dev
    ```
    (This starts `webpack serve` for the library at `localhost:8090` and watches the host)

2.  Open the Extension Settings (click extension icon -> Options, or right-click -> Options).
3.  Check **Enable Debug Mode**.
4.  Ensure URL is `http://localhost:8090/library.bundle.js`.
5.  Save.
6.  Reload any target page. The extension will fetch the latest library code from localhost.

#### Release Mode
1.  Uncheck **Debug Mode** in settings.
2.  The extension will load the bundled version of the library.
3.  To update the bundled version, run `npm run build` and reload the extension in `chrome://extensions`.

## Project Structure

*   `src/host`: The Extension shell (Manifest, Loader, React Settings).
*   `src/library`: The "Guest" library containing the actual logic.
*   `webpack.host.js`: Builds the extension.
*   `webpack.lib.js`: Builds the library.
