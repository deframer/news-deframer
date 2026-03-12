const parseDateInput = (dateStr: string | Date) => {
  if (dateStr instanceof Date) {
    return dateStr;
  }

  const trimmed = dateStr.trim();
  const direct = new Date(trimmed);
  if (!Number.isNaN(direct.getTime())) {
    return direct;
  }

  if (/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}(:\d{2})?$/.test(trimmed)) {
    const normalized = `${trimmed.replace(' ', 'T')}${trimmed.includes(':') && trimmed.split(':').length === 2 ? ':00' : ''}Z`;
    return new Date(normalized);
  }

  if (/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) {
    return new Date(`${trimmed}T00:00:00Z`);
  }

  return new Date(NaN);
};

export const getRelativeTime = (dateStr: string | Date, locale: string): string => {
  const date = parseDateInput(dateStr);
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
