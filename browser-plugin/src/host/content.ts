import Library from '../library/index';
import { Config } from '../library/api';

const DEFAULT_BACKEND_URL = 'http://localhost:8080';

function loadLibrary() {
  chrome.storage.local.get(['backendUrl', 'username', 'password'], (items) => {
    const { backendUrl, username, password } = items as Record<string, string | undefined>;
    const config: Config = {
        backendUrl: backendUrl ?? DEFAULT_BACKEND_URL,
        username,
        password
    };

    console.log("[NDF Host] Starting library...");
    Library.start(config);
  });
}

loadLibrary();
