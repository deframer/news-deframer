import { installPreemptiveDomGuard } from '../ndf/domGuard';
import * as ndf from '../ndf/index';
import log from '../shared/logger';
import { getSettings } from '../shared/settings';

const readEnabledFlag = (): Promise<boolean> => {
  return new Promise((resolve) => {
    chrome.storage.local.get({ enabled: false }, (result) => {
      resolve(Boolean(result.enabled));
    });
  });
};

async function startNdf() {
  if (sessionStorage.getItem('__ndf-bypass')) {
    log.debug('Bypass flag detected. Skipping NDF for this load.');
    sessionStorage.removeItem('__ndf-bypass');
    return;
  }

  // Install guard immediately to block popups, even before we know if we are enabled.
  // If we are disabled, we will reload the page without the guard.
  const guard = installPreemptiveDomGuard({ allowedIds: ['ndf-root', 'ndf-global-styles'] });

  const enabled = await readEnabledFlag();
  if (!enabled) {
    log.debug('NDF is disabled by user settings. Reloading without guard.');
    guard?.release();
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
    return;
  }

  const settings = await getSettings();

  log.debug('Starting NDF...');
  try {
    await ndf.start(settings);
  } finally {
    guard?.release();
  }
}

// Early exit if the content is not HTML.
// This is to prevent the script from running on non-HTML content like XML, JSON, etc.
if (document.contentType === 'text/html') {
  startNdf();
}
