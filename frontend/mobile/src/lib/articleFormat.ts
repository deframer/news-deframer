import { AppPalette } from '../theme';

export const stripHtml = (value?: string): string => {
  if (!value) {
    return '';
  }

  return value.replace(/<[^>]+>/g, ' ').replace(/\s+/g, ' ').trim();
};

export const toPercent = (value?: number): number => Math.round((value || 0) * 100);

export const getRatingColors = (percentage: number, palette: AppPalette) => {
  if (percentage < 11) {
    return { backgroundColor: palette.success, textColor: palette.text };
  }

  if (percentage < 34) {
    return { backgroundColor: palette.success, textColor: '#ffffff' };
  }

  if (percentage < 67) {
    return { backgroundColor: palette.warning, textColor: '#000000' };
  }

  return { backgroundColor: palette.danger, textColor: '#ffffff' };
};
