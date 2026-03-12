import { useTranslation } from 'react-i18next';

import { formatTime } from '../../../../shared/formatTime';

interface MetaDataProps {
  pubDate?: string | Date;
  author?: string;
  category?: string;
}

export const MetaData = ({ pubDate, author, category }: MetaDataProps) => {
  const { i18n } = useTranslation();
  // If no metadata is available, don't render anything to save space
  if (!pubDate && !author && !category) return null;

  let timeAgo = '';
  if (pubDate) {
    timeAgo = formatTime(pubDate, i18n.language);
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
