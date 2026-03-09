import log from 'loglevel';
import prefix from 'loglevel-plugin-prefix';

import { LOG_LEVEL } from './loglevel';

const NDF_PREFIX = '***** NDF';

prefix.reg(log);
prefix.apply(log, {
  format(level) {
    return `${NDF_PREFIX} [${level}]`;
  },
});

log.setLevel(LOG_LEVEL as log.LogLevelDesc);

export default log;
