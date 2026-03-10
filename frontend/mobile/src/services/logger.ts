import { logger as createLogger } from 'react-native-logs';

const isDevRuntime = typeof __DEV__ === 'boolean' ? __DEV__ : true;

const nativeLogger = createLogger.createLogger({
  severity: isDevRuntime ? 'debug' : 'info',
  transportOptions: {
    colors: false,
  },
  printLevel: true,
  printDate: false,
  enabled: true,
});

const withPrefix = (message: string) => `[news-deframer-mobile] ${message}`;

export const logger = {
  info: (message: string, meta?: unknown) => {
    if (meta === undefined) {
      nativeLogger.info(withPrefix(message));
      return;
    }

    nativeLogger.info(withPrefix(message), meta);
  },
  warn: (message: string, meta?: unknown) => {
    if (meta === undefined) {
      nativeLogger.warn(withPrefix(message));
      return;
    }

    nativeLogger.warn(withPrefix(message), meta);
  },
  error: (message: string, meta?: unknown) => {
    if (meta === undefined) {
      nativeLogger.error(withPrefix(message));
      return;
    }

    nativeLogger.error(withPrefix(message), meta);
  },
};
