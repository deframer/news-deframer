const formatRelativeFallback = (locale: string, value: number, unit: 'minute' | 'hour' | 'day' | 'week' | 'month' | 'year') => {
  const isGerman = locale.toLowerCase().startsWith('de');
  const absoluteValue = Math.abs(value);

  if (isGerman) {
    const labels = {
      minute: 'Min.',
      hour: 'Std.',
      day: 'T.',
      week: 'Wo.',
      month: 'Mon.',
      year: 'J.',
    } as const;

    return `vor ${absoluteValue} ${labels[unit]}`;
  }

  const labels = {
    minute: 'm',
    hour: 'h',
    day: 'd',
    week: 'w',
    month: 'mo',
    year: 'y',
  } as const;

  return `${absoluteValue}${labels[unit]} ago`;
};

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

  const seconds = (Date.now() - date.getTime()) / 1000;
  const absoluteSeconds = Math.abs(seconds);
  const safeLocale = typeof locale === 'string' && locale.trim() ? locale : 'en';
  const rtf = RelativeTimeFormatCtor ? new RelativeTimeFormatCtor(safeLocale, { style: 'narrow' }) : null;

  if (absoluteSeconds < 60) return '';
  if (absoluteSeconds < 3600) {
    const value = -Math.round(seconds / 60);
    return rtf ? rtf.format(value, 'minute') : formatRelativeFallback(safeLocale, value, 'minute');
  }
  if (absoluteSeconds < 86400) {
    const value = -Math.round(seconds / 3600);
    return rtf ? rtf.format(value, 'hour') : formatRelativeFallback(safeLocale, value, 'hour');
  }
  if (absoluteSeconds < 604800) {
    const value = -Math.round(seconds / 86400);
    return rtf ? rtf.format(value, 'day') : formatRelativeFallback(safeLocale, value, 'day');
  }
  if (absoluteSeconds < 2592000) {
    const value = -Math.round(seconds / 604800);
    return rtf ? rtf.format(value, 'week') : formatRelativeFallback(safeLocale, value, 'week');
  }
  if (absoluteSeconds < 31536000) {
    const value = -Math.round(seconds / 2592000);
    return rtf ? rtf.format(value, 'month') : formatRelativeFallback(safeLocale, value, 'month');
  }

  const value = -Math.round(seconds / 31536000);
  return rtf ? rtf.format(value, 'year') : formatRelativeFallback(safeLocale, value, 'year');
};
