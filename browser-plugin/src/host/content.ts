import Library from '../library/index';
import { Config } from '../library/api';

const DEFAULT_BACKEND_URL = 'http://localhost:8080';

function loadLibrary() {
  chrome.storage.local.get(['backendUrl', 'username', 'password'], (items) => {
    const config: Config = {
        backendUrl: items.backendUrl ?? DEFAULT_BACKEND_URL,
        username: items.username,
        password: items.password
    };

    console.log("[NDF Host] Starting library...");
    Library.start(config);
  });
}

loadLibrary();
