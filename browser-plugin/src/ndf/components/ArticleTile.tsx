import { useTranslation } from 'react-i18next';

import { AnalyzedItem } from '../client';
import { stripHtml } from '../utils/html-utils';
import { MetaData } from './MetaData';
import { RatingBarOverlay } from './RatingBarOverlay';

interface ArticleTileProps {
  item: AnalyzedItem;
}

export const ArticleTile = ({ item }: ArticleTileProps) => {
  const { t } = useTranslation();
  const title = item.title_corrected || stripHtml(item.title_original) || t('article.no_title', 'No Title');
  const description = item.description_corrected || stripHtml(item.description_original) || t('article.no_description', 'No Description');
  const imageUrl = item.media?.medium === 'image' ? item.media.url : '';

  return (
    <a href={item.url} className="tile-link">
      <div className="tile">
        <RatingBarOverlay value={item.rating} reason={item.overall_reason} />
        {imageUrl ? (
          <div className="image-container">
            <img src={imageUrl} alt={item.media?.description || ''} />
          </div>
        ) : null}
        <div className="content">
          <h3>{title}</h3>
          <p>{description}</p>
        </div>
        <div className="tile-footer">
          <MetaData pubDate={(item as AnalyzedItem & { pubDate?: string | Date }).pubDate} />
        </div>
      </div>
    </a>
  );
};
