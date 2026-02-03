import { useTranslation } from 'react-i18next';

import { AnalyzedItem } from '../client';
import { MetaData } from './MetaData';
import { RatingBarOverlay } from './RatingBarOverlay';

const tileCss = `
  .tile-link { text-decoration: none; color: inherit; display: block; height: 100%; }
  .tile { background-color: var(--card-bg); border-radius: 8px; box-shadow: var(--card-shadow); overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; height: 100%; position: relative; display: flex; flex-direction: column; }
  .tile:hover { transform: translateY(-5px); box-shadow: var(--card-shadow-hover); }
  .image-container { position: relative; }
  .image-container img { width: 100%; height: auto; display: block; }
  .content { padding: 15px; flex: 1; }
  h3 { margin: 0 0 10px; font-size: 1.1em; color: var(--text-color); }
  p { font-size: 0.9em; color: var(--secondary-text); margin: 0; }
  .tile-footer {
    height: 30px;
    padding: 0 15px;
    display: flex;
    align-items: center;
    border-top: 1px solid var(--border-color, #eee);
  }
  .tile-footer .meta-data {
    margin: 0;
    font-size: 0.9em;
  }
`;

interface ArticleTileProps {
  item: AnalyzedItem;
}

export const ArticleTile = ({ item }: ArticleTileProps) => {
  const { t } = useTranslation();
  const title = item.title_corrected || item.title_original || t('article.no_title');
  const description = item.description_corrected || item.description_original || t('article.no_description');
  const imageUrl = item.media?.medium === 'image' ? item.media.url : '';

  return (
    <a href={item.url} className="tile-link">
        <style>{tileCss}</style>
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
