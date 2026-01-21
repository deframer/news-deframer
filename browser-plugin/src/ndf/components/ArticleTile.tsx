import React from 'react';

import { AnalyzedItem } from '../client';
import { RatingBarOverlay } from './RatingBarOverlay';

const tileCss = `
  .tile-link { text-decoration: none; color: inherit; display: block; height: 100%; }
  .tile { background-color: #fff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; height: 100%; position: relative; }
  .tile:hover { transform: translateY(-5px); box-shadow: 0 8px 12px rgba(0,0,0,0.15); }
  .image-container { position: relative; }
  .image-container img { width: 100%; height: auto; display: block; }
  .content { padding: 15px; }
  h3 { margin: 0 0 10px; font-size: 1.1em; }
  p { font-size: 0.9em; color: #666; margin: 0; }
  
  .image-container .image-overlay-bar {
    /* This class is no longer used */
  }
  .image-overlay-bar .reason {
      display: none; /* Hide reason text on overlays */
  }
`;

interface ArticleTileProps {
  item: AnalyzedItem;
}

export const ArticleTile = ({ item }: ArticleTileProps) => {
  const title = item.title_corrected || item.title_original || 'No title';
  const description = item.description_corrected || item.description_original || 'No description';
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
      </div>
    </a>
  );
};
