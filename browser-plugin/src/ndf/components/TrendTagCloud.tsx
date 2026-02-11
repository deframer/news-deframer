import { ParentSize } from '@visx/responsive';
import { scaleLog } from '@visx/scale';
import { Text } from '@visx/text';
import { Wordcloud } from '@visx/wordcloud';
import { memo, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, TrendMetric } from '../client';
import { TrendDetails } from './TrendDetails';


// https://visx.airbnb.tech/wordcloud

interface TrendTagCloudProps {
  domain: DomainEntry;
  days: number;
  searchEngineUrl: string;
}

const TrendWordCloud = memo(({ width, height, words, selectedTerm, onSelect, onHover }: {
  width: number;
  height: number;
  words: { text: string; value: number; original: TrendMetric & { rank: number } }[];
  selectedTerm: string | null;
  onSelect: (term: string | null) => void;
  onHover: (data: { x: number; y: number; item: TrendMetric & { rank: number } } | null) => void;
}) => {
  const fontScale = useMemo(() => {
    if (words.length === 0) {
      return scaleLog({ domain: [1, 10], range: [14, 50] });
    }
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

  // Memoize accessors to prevent re-layout on every render
  const getFontSize = useCallback((datum: { value: number }) => fontScale(datum.value), [fontScale]);
  const fixedRandom = useCallback(() => 0.5, []);

  // Use CSS variables for colors to match theme
  const colors = ['var(--text-color)', 'var(--primary-color)', 'var(--secondary-text)'];

  const getColor = (text: string) => {
    let hash = 0;
    for (let i = 0; i < text.length; i++) {
      hash = text.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  };

  return (
    <Wordcloud
      words={words}
      width={width}
      height={height}
      fontSize={getFontSize}
      font={'sans-serif'}
      padding={4}
      spiral={'archimedean'}
      rotate={0}
      random={fixedRandom}
    >
      {(cloudWords) =>
        cloudWords.map((w) => (
          <g key={w.text} transform={`translate(${w.x}, ${w.y}) rotate(${w.rotate})`}>
            <Text
              className="word-cloud-text"
              fill={selectedTerm === w.text ? 'var(--accent-color)' : getColor(w.text || '')}
              textAnchor={'middle'}
              fontSize={w.size}
              fontWeight={selectedTerm === w.text ? 'bold' : 'normal'}
              textDecoration={selectedTerm === w.text ? 'underline' : 'none'}
              fontFamily={w.font}
              onClick={() => onSelect(w.text || null)}
              onMouseEnter={() => onHover({ x: w.x || 0, y: w.y || 0, item: (w as unknown as { original: TrendMetric & { rank: number } }).original })}
              onMouseLeave={() => onHover(null)}
            >
              {w.text}
            </Text>
          </g>
        ))
      }
    </Wordcloud>
  );
});
TrendWordCloud.displayName = 'TrendWordCloud';

export const TrendTagCloud = ({ domain, days, searchEngineUrl }: TrendTagCloudProps) => {
  const { t } = useTranslation();
  const [selectedTerm, setSelectedTerm] = useState<string | null>(null);
  const [tooltip, setTooltip] = useState<{ x: number; y: number; item: TrendMetric & { rank: number } } | null>(null);
  const [items, setItems] = useState<TrendMetric[]>([]);
  const [loading, setLoading] = useState(false);

  const handleSearch = () => {
    if (!selectedTerm) return;
    const query = encodeURIComponent(`${selectedTerm} site:${domain.domain}`);
    const baseUrl = searchEngineUrl.replace(/\/$/, '');
    window.open(`${baseUrl}/search?q=${query}`, '_blank');
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const data = await client.getTopTrendByDomain(domain.domain, domain.language, days);
        setItems(data);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [domain, days]);

  const words = useMemo(() => {
    return items.map((i, index) => ({
      text: i.trend_topic,
      value: i.outlier_ratio,
      original: { ...i, rank: index + 1 },
    })).sort((a, b) => b.value - a.value);
  }, [items]);

  // Memoize the render function for ParentSize to prevent it from triggering updates
  // when TrendTagCloud re-renders (e.g. tab switch) but dimensions haven't changed.
  const renderCloud = useCallback(({ width, height }: { width: number; height: number }) => (
    <>
      <TrendWordCloud
        width={width}
        height={height}
        words={words}
        selectedTerm={selectedTerm}
        onSelect={setSelectedTerm}
        onHover={setTooltip}
      />
      {tooltip && (
        <div
          className="cloud-tooltip"
          style={{
            visibility: 'visible',
            opacity: 1,
            position: 'absolute',
            left: width / 2 + tooltip.x,
            top: height / 2 + tooltip.y,
            transform: 'translate(-50%, -100%)',
            marginTop: '-5px',
            pointerEvents: 'none',
            bottom: 'auto',
          }}
        >
          {t('trends.rank', 'Rank')}: {tooltip.item.rank}<br/>
          {t('trends.trend', 'Trend')}: {tooltip.item.outlier_ratio.toFixed(2)}x<br/>
          {t('trends.vol', 'Vol')}: {tooltip.item.frequency}
        </div>
      )}
    </>
  ), [words, selectedTerm, tooltip, t]);

  if (loading) {
    return <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <div className="spinner-small" />
    </div>;
  }

  if (items.length === 0) {
    return (
      <div style={{ height: '400px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--secondary-text)' }}>
        {t('trends.no_data', 'No trending topics found for this period.')}
      </div>
    );
  }

  return (
    <>
      <div
        className="tag-cloud"
        style={{ width: '100%', height: selectedTerm ? '250px' : '400px', position: 'relative' }}
      >
        {selectedTerm && (
          <div style={{ position: 'absolute', top: '10px', right: '10px', zIndex: 10 }}>
            <button
              onClick={handleSearch}
              className="open-icon-btn"
              type="button"
              title={`${selectedTerm} • ${domain.domain}`}
              style={{
                background: 'var(--bg-color, #fff)',
                border: '1px solid var(--border-color, #ccc)',
                borderRadius: '4px',
                cursor: 'pointer',
                padding: '5px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: 'var(--text-color)',
                boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
              }}
            >
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"></path>
                <polyline points="15 3 21 3 21 9"></polyline>
                <line x1="10" y1="14" x2="21" y2="3"></line>
              </svg>
            </button>
          </div>
        )}
        <ParentSize>
          {renderCloud}
        </ParentSize>
      </div>
      {selectedTerm && (
        <TrendDetails term={selectedTerm} domain={domain} days={days} />
      )}
    </>
  );
};