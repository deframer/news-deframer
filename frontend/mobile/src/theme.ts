import { ColorSchemeName } from 'react-native';

export type ThemeMode = 'light' | 'dark' | 'system';

export interface AppPalette {
  background: string;
  text: string;
  card: string;
  cardShadow: string;
  border: string;
  secondaryText: string;
  buttonBackground: string;
  buttonText: string;
  buttonBorder: string;
  buttonHover: string;
  accent: string;
  accentHover: string;
  accentText: string;
  ratingBackground: string;
  success: string;
  warning: string;
  danger: string;
  secondaryBackground: string;
  badgeBackground: string;
}

const lightPalette: AppPalette = {
  background: '#f0f2f5',
  text: '#333333',
  card: '#ffffff',
  cardShadow: 'rgba(0,0,0,0.1)',
  border: '#eeeeee',
  secondaryText: '#666666',
  buttonBackground: '#ffffff',
  buttonText: '#333333',
  buttonBorder: '#dddddd',
  buttonHover: '#f8f9fa',
  accent: '#0056b3',
  accentHover: '#004494',
  accentText: '#ffffff',
  ratingBackground: '#e9ecef',
  success: '#198754',
  warning: '#ffc107',
  danger: '#b02a37',
  secondaryBackground: '#f7f8fa',
  badgeBackground: '#e9ecef',
};

const darkPalette: AppPalette = {
  background: '#18191a',
  text: '#e4e6eb',
  card: '#242526',
  cardShadow: 'rgba(0,0,0,0.35)',
  border: '#3e4042',
  secondaryText: '#b0b3b8',
  buttonBackground: '#3a3b3c',
  buttonText: '#e4e6eb',
  buttonBorder: '#3e4042',
  buttonHover: '#4e4f50',
  accent: '#2d88ff',
  accentHover: '#4599ff',
  accentText: '#000000',
  ratingBackground: '#3a3b3c',
  success: '#4caf50',
  warning: '#ffca28',
  danger: '#f44336',
  secondaryBackground: '#18191a',
  badgeBackground: '#3a3b3c',
};

export const resolveThemeMode = (theme: ThemeMode, colorScheme: ColorSchemeName): Exclude<ThemeMode, 'system'> => {
  if (theme === 'system') {
    return colorScheme === 'dark' ? 'dark' : 'light';
  }

  return theme;
};

export const getPalette = (theme: ThemeMode, colorScheme: ColorSchemeName): AppPalette => {
  const resolved = resolveThemeMode(theme, colorScheme);
  return resolved === 'dark' ? darkPalette : lightPalette;
};
