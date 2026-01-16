# News Deframer - Chrome Extension

This is the Browser Extension wrapper for the News Deframer library.

It uses the exact same `news-deframer-lib.js` as the Tampermonkey scripts, but runs it in the native Extension context, enabling `chrome.storage.local`.

## 🛠️ How to Build

Run this command from the project root:

```bash
make ext
```

## 📦 Remote Development (Windows / VS Code SSH)

If you are developing remotely and want to install without unzipping:

1. **Pack:**
   Run this command on the server:
   ```bash
   make pack
   ```
   This creates `extension.zip`.

2. **Download & Extract:**
   - Download `extension.zip` to Windows.
   - **Right-click -> Extract All**.

3. **Install:**
   - Open `chrome://extensions`.
   - Toggle **Developer mode**.
   - Click **Load unpacked**.
   - Select the extracted `extension` folder.

## 🔍 How to Debug & Explore Storage

Since this extension uses `chrome.storage.local`, you cannot see the data in the "Application" tab (which is for LocalStorage).

**To see the database:**

1. Open any webpage.
2. Open DevTools (**F12**).
3. Go to the **Console** tab.
4. **Crucial:** In the top toolbar of the console, click the dropdown that says **"Top"** (or "Main") and select **"News Deframer Extension"**.
5. Run this command to see all data:
   ```js
   chrome.storage.local.get(null, console.log)
   ```
6. Run this to clear data:
   ```js
   chrome.storage.local.clear()
   ```
