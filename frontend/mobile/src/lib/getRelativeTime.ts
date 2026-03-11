export const getRelativeTime = (dateStr: string | Date, locale: string): string => {
  const date = new Date(dateStr);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  const RelativeTimeFormatCtor =
    typeof Intl !== 'undefined' && typeof Intl.RelativeTimeFormat === 'function'
      ? Intl.RelativeTimeFormat
      : null;
  if (!RelativeTimeFormatCtor) {
    return '';
  }

  const seconds = (Date.now() - date.getTime()) / 1000;
  const absoluteSeconds = Math.abs(seconds);
  const safeLocale = typeof locale === 'string' && locale.trim() ? locale : 'en';
  const rtf = new RelativeTimeFormatCtor(safeLocale, { style: 'narrow' });

  if (absoluteSeconds < 60) return '';
  if (absoluteSeconds < 3600) return rtf.format(-Math.round(seconds / 60), 'minute');
  if (absoluteSeconds < 86400) return rtf.format(-Math.round(seconds / 3600), 'hour');
  if (absoluteSeconds < 604800) return rtf.format(-Math.round(seconds / 86400), 'day');
  if (absoluteSeconds < 2592000) return rtf.format(-Math.round(seconds / 604800), 'week');
  if (absoluteSeconds < 31536000) return rtf.format(-Math.round(seconds / 2592000), 'month');

  return rtf.format(-Math.round(seconds / 31536000), 'year');
};
