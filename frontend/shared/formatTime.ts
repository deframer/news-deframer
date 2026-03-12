import { formatDistanceToNowStrict, parseISO } from 'date-fns';
import { de, enUS } from 'date-fns/locale';

const resolveDateFnsLocale = (locale: string) => (locale.toLowerCase().startsWith('de') ? de : enUS);

export const parseTimeInput = (dateStr: string | Date) => {
  if (dateStr instanceof Date) {
    return dateStr;
  }

  const parsed = parseISO(dateStr);
  if (!Number.isNaN(parsed.getTime())) {
    return parsed;
  }

  const direct = new Date(dateStr);
  return Number.isNaN(direct.getTime()) ? new Date(NaN) : direct;
};

export const formatShortDate = (dateStr: string | Date, locale: string): string => {
  const date = parseTimeInput(dateStr);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  return date.toLocaleDateString(locale, { month: 'short', day: 'numeric' });
};

export const formatTime = (dateStr: string | Date, locale: string): string => {
  const date = parseTimeInput(dateStr);
  if (Number.isNaN(date.getTime())) {
    return '';
  }

  return formatDistanceToNowStrict(date, {
    addSuffix: true,
    locale: resolveDateFnsLocale(locale),
  });
};
