import { ParentSize } from '@visx/responsive';
import { scaleLog } from '@visx/scale';
import { Text } from '@visx/text';
import { Wordcloud } from '@visx/wordcloud';
import { useMemo } from 'react';

import { TrendItem } from './TabTrend';

interface TrendTagCloudProps {
  items: TrendItem[];
}

export const TrendTagCloud = ({ items }: TrendTagCloudProps) => {
  if (items.length === 0) return null;

  const words = useMemo(() => {
    return items.map((i) => ({
      text: i.word,
      value: i.outlierRatio,
    }));
  }, [items]);

  const fontScale = useMemo(() => {
    const values = words.map((w) => w.value);
    const minVal = Math.min(...values);
    const maxVal = Math.max(...values);

    // Avoid log(0) or invalid domains
    const safeMin = Math.max(minVal, 0.1);
    const safeMax = Math.max(maxVal, safeMin * 1.1);

    return scaleLog({
      domain: [safeMin, safeMax],
      range: [14, 50],
    });
  }, [words]);

  // Use CSS variables for colors to match theme
  const colors = ['var(--text-color)', 'var(--primary-color)', 'var(--secondary-text)'];

  return (
    <div className="tag-cloud" style={{ width: '100%', height: '400px' }}>
      <ParentSize>
        {({ width, height }) => (
          <Wordcloud
            words={words}
            width={width}
            height={height}
            fontSize={(datum) => fontScale(datum.value)}
            font={'sans-serif'}
            padding={4}
            spiral={'archimedean'}
            rotate={0}
            random={() => 0.5}
          >
            {(cloudWords) =>
              cloudWords.map((w, i) => (
                <g key={w.text} transform={`translate(${w.x}, ${w.y}) rotate(${w.rotate})`}>
                  <Text
                    className="word-cloud-text"
                    fill={colors[i % colors.length]}
                    textAnchor={'middle'}
                    fontSize={w.size}
                    fontFamily={w.font}
                  >
                    {w.text}
                  </Text>
                </g>
              ))
            }
          </Wordcloud>
        )}
      </ParentSize>
    </div>
  );
};