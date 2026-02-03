import log from 'loglevel';
import prefix from 'loglevel-plugin-prefix';

import { LOG_LEVEL } from './mode';

const NDF_PREFIX = '***** NDF';

// Attach the prefix plugin
prefix.reg(log);

// Configure the prefix
prefix.apply(log, {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  format(level, _name, _timestamp) {
    return `${NDF_PREFIX} [${level}]`;
  },
});

// Set the initial log level
log.setLevel(LOG_LEVEL as log.LogLevelDesc);

export default log;
