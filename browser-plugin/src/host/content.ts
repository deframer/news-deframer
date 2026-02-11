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

  const enabled = await readEnabledFlag();
  if (!enabled) {
    log.debug('NDF is disabled by user settings.');
    return;
  }

  // Install guard only when enabled so hostile scripts never execute while we fetch settings.
  const guard = installPreemptiveDomGuard({ allowedIds: ['ndf-root', 'ndf-global-styles'] });

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
