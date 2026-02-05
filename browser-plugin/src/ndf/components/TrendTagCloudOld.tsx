import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { TrendItem } from './TabTrend';
import { TrendContext } from './TrendContext';

// https://visx.airbnb.tech/wordcloud

interface TrendTagCloudOldProps {
  items: TrendItem[];
  days: number;
  domain: string;
}

export const TrendTagCloudOld = ({ items, days, domain }: TrendTagCloudOldProps) => {
  const { t } = useTranslation();
  const [selectedTopic, setSelectedTopic] = useState<string | null>(null);
  if (items.length === 0) return null;

  // Use Outlier Ratio for sizing as it indicates "Trendiness" better than raw count
  const maxVal = Math.max(...items.map(i => i.outlierRatio));
  const minVal = Math.min(...items.map(i => i.outlierRatio));

  const getFontSize = (val: number) => {
    if (maxVal === minVal) return '1.2em';
    const minSize = 0.9;
    const maxSize = 2.2;
    // Linear interpolation
    const size = minSize + ((val - minVal) / (maxVal - minVal)) * (maxSize - minSize);
    return `${size.toFixed(2)}em`;
  };

  // Sort alphabetically for cloud display
  const cloudItems = [...items].sort((a, b) => a.word.localeCompare(b.word));

  return (
    <>
      <div className="tag-cloud">
        {cloudItems.map((item) => (
          <span
            key={item.word}
            className={`tag-item ${selectedTopic === item.word ? 'active' : ''}`}
            style={{ fontSize: getFontSize(item.outlierRatio) }}
            onClick={() => setSelectedTopic(selectedTopic === item.word ? null : item.word)}
          >
            {item.word}
            <div className="cloud-tooltip">
              {t('trends.rank', 'Rank')}: {item.rank}<br/>
              {t('trends.trend', 'Trend')}: {item.outlierRatio.toFixed(2)}x<br/>
              {t('trends.vol', 'Vol')}: {item.count}
            </div>
          </span>
        ))}

        {selectedTopic && (
          <div className="cloud-context-container">
            <TrendContext topic={selectedTopic} className="embedded-context" days={days} domain={domain} />
          </div>
        )}
      </div>
    </>
  );
};