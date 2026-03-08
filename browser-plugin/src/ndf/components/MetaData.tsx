import { useTranslation } from 'react-i18next';

interface MetaDataProps {
  pubDate?: string | Date;
  author?: string;
  category?: string;
}

const getRelativeTime = (dateStr: string | Date, locale: string): string => {
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '';

  const seconds = (new Date().getTime() - date.getTime()) / 1000;
  const absoluteSeconds = Math.abs(seconds);
  const rtf = new Intl.RelativeTimeFormat(locale, { style: 'narrow' });

  if (absoluteSeconds < 60) return '';
  if (absoluteSeconds < 3600) return rtf.format(-Math.round(seconds / 60), 'minute');
  if (absoluteSeconds < 86400) return rtf.format(-Math.round(seconds / 3600), 'hour');
  if (absoluteSeconds < 604800) return rtf.format(-Math.round(seconds / 86400), 'day');
  if (absoluteSeconds < 2592000) return rtf.format(-Math.round(seconds / 604800), 'week');
  if (absoluteSeconds < 31536000) return rtf.format(-Math.round(seconds / 2592000), 'month');
  return rtf.format(-Math.round(seconds / 31536000), 'year');
};

export const MetaData = ({ pubDate, author, category }: MetaDataProps) => {
  const { i18n, t } = useTranslation();
  // If no metadata is available, don't render anything to save space
  if (!pubDate && !author && !category) return null;

  let timeAgo = '';
  if (pubDate) {
    timeAgo = getRelativeTime(pubDate, i18n.language) || t('metadata.just_now', 'a moment ago');
  }

  const showSeparator = Boolean(timeAgo && author);

  return (
    <>
      <div className="meta-data">
        {timeAgo && (
          <div className="meta-item">
            <svg className="meta-icon" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10" />
              <polyline points="12 6 12 12 16 14" />
            </svg>
            <span>{timeAgo}</span>
          </div>
        )}

        {showSeparator && <span className="meta-separator">|</span>}

        {author && (
          <div className="meta-item">
            <span>{author}</span>
          </div>
        )}
      </div>
    </>
  );
};
