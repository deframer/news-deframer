import React from 'react';

const metaDataCss = `
  .meta-data {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 12px;
    font-size: 0.85em;
    color: var(--secondary-text);
    margin-top: 0.5em;
    margin-bottom: 2em;
    line-height: 1;
  }
  .meta-item {
    display: flex;
    align-items: center;
    gap: 6px;
    white-space: nowrap;
  }
  .meta-icon {
    width: 14px;
    height: 14px;
    stroke: currentColor;
    stroke-width: 2;
    fill: none;
  }
`;

interface MetaDataProps {
  pubDate?: string | Date;
  author?: string;
  category?: string;
}

const getRelativeTime = (dateStr: string | Date): string => {
  const date = new Date(dateStr);
  if (isNaN(date.getTime())) return '';

  const seconds = (new Date().getTime() - date.getTime()) / 1000;
  const rtf = new Intl.RelativeTimeFormat('en', { style: 'narrow' });

  if (seconds < 60) return rtf.format(-Math.floor(seconds), 'second');
  if (seconds < 3600) return rtf.format(-Math.floor(seconds / 60), 'minute');
  if (seconds < 86400) return rtf.format(-Math.floor(seconds / 3600), 'hour');
  if (seconds < 604800) return rtf.format(-Math.floor(seconds / 86400), 'day');
  if (seconds < 2592000) return rtf.format(-Math.floor(seconds / 604800), 'week');
  if (seconds < 31536000) return rtf.format(-Math.floor(seconds / 2592000), 'month');
  return rtf.format(-Math.floor(seconds / 31536000), 'year');
};

export const MetaData = ({ pubDate, author, category }: MetaDataProps) => {
  // If no metadata is available, don't render anything to save space
  if (!pubDate && !author && !category) return null;

  let timeAgo = '';
  if (pubDate) {
    timeAgo = getRelativeTime(pubDate);
  }

  return (
    <>
      <style>{metaDataCss}</style>
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

        {/* Placeholders for future author/category implementation */}
        {author && (
          <div className="meta-item">
             <span>{author}</span>
          </div>
        )}
      </div>
    </>
  );
};