import { ColorSchemeName } from 'react-native';

import { getPalette } from '../theme';
import { Settings } from './settingsService';

export const themeService = {
  getPalette(settings: Settings, colorScheme: ColorSchemeName) {
    return getPalette(settings.theme, colorScheme);
  },
};
