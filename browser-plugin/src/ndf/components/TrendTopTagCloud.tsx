import { useState } from 'react';
import { TrendItem } from './TrendTop';
import { TrendContext } from './TrendContext';

interface TrendTopTagCloudProps {
  items: TrendItem[];
}

const cloudCss = `
  .tag-cloud {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    align-content: center;
    gap: 12px;
    padding: 20px 10px;
    min-height: 200px;
  }
  .tag-item {
    transition: transform 0.2s, color 0.2s;
    cursor: pointer;
    line-height: 1;
    color: var(--text-color);
  }
  .tag-item:hover {
    transform: scale(1.1);
    color: var(--primary-color, #0056b3);
  }
  .tag-item.active {
    color: var(--primary-color, #0056b3);
    font-weight: bold;
    background-color: var(--bg-color-secondary, #f0f0f0);
    padding: 4px 10px;
    border-radius: 12px;
    transform: scale(1.1);
    box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  }
  .cloud-context-container {
    margin-top: 10px;
    animation: fadeIn 0.3s ease;
  }
`;

export const TrendTopTagCloud = ({ items }: TrendTopTagCloudProps) => {
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
      <style>{cloudCss}</style>
      <div className="tag-cloud">
        {cloudItems.map((item) => (
          <span
            key={item.word}
            className={`tag-item ${selectedTopic === item.word ? 'active' : ''}`}
            style={{ fontSize: getFontSize(item.outlierRatio) }}
            title={`Rank: ${item.rank}, Trend: ${item.outlierRatio.toFixed(2)}x, Vol: ${item.count}`}
            onClick={() => setSelectedTopic(selectedTopic === item.word ? null : item.word)}
          >
            {item.word}
          </span>
        ))}

        {selectedTopic && (
          <div className="cloud-context-container">
            <TrendContext topic={selectedTopic} className="embedded-context" />
          </div>
        )}
      </div>
    </>
  );
};