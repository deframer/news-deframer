import ndf from '../ndf/index';
import log from '../shared/logger';
import { getSettings } from '../shared/settings';

async function startNdf() {
  const settings = await getSettings();
  log.info('Starting NDF...');
  ndf.start(settings);
}

startNdf();
