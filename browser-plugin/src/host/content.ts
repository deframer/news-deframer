import ndf from '../ndf/index';
import log from '../shared/logger';
import { getSettings } from '../shared/settings';

async function startNdf() {
  const settings = await getSettings();

  if (settings.enabled === false) {
    log.info('NDF is disabled by user settings.');
    return;
  }

  log.info('Starting NDF...');
  ndf.start(settings);
}

// Early exit if the content is not HTML.
// This is to prevent the script from running on non-HTML content like XML, JSON, etc.
if (document.contentType === 'text/html') {
  startNdf();
}
