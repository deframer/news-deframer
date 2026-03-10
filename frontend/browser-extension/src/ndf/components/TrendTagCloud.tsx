import { Dispatch, SetStateAction, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { buildWordCloudWords, getWordCloudColor, layoutWordCloud, MeasureWordOptions, WordCloudWord } from '../../shared/wordcloud';
import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, TrendMetric } from '../client';
import { TrendDetails } from './TrendDetails';

interface TrendTagCloudProps {
  domain: DomainEntry;
  days: number;
  searchEngineUrl: string;
  activeTab: 'lifecycle' | 'context' | 'articles';
  setActiveTab: Dispatch<SetStateAction<'lifecycle' | 'context' | 'articles'>>;
}

const BULLET_DELIMITER = '•';
const DEFAULT_WIDTH = 640;
const DEFAULT_HEIGHT = 400;

const approximateMeasureWord = ({ text, fontSize }: MeasureWordOptions) => ({
  width: Math.max(fontSize, text.length * fontSize * 0.58),
  height: fontSize * 1.2,
});

const createCanvasMeasureWord = () => {
  if (typeof document === 'undefined') {
    return approximateMeasureWord;
  }

  const canvas = document.createElement('canvas');
  let context: CanvasRenderingContext2D | null = null;

  try {
    context = canvas.getContext('2d');
  } catch {
    return approximateMeasureWord;
  }

  if (!context) {
    return approximateMeasureWord;
  }

  return ({ text, fontSize, fontFamily, fontWeight = 'normal' }: MeasureWordOptions) => {
    context.font = `${fontWeight} ${fontSize}px ${fontFamily}`;
    const metrics = context.measureText(text);

    return {
      width: Math.max(metrics.width, fontSize),
      height: fontSize * 1.2,
    };
  };
};

export const TrendTagCloud = ({ domain, days, searchEngineUrl, activeTab, setActiveTab }: TrendTagCloudProps) => {
  const { t } = useTranslation();
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [selectedTerm, setSelectedTerm] = useState<string | null>(null);
  const [tooltip, setTooltip] = useState<WordCloudWord<TrendMetric & { rank: number }> | null>(null);
  const [items, setItems] = useState<TrendMetric[]>([]);
  const [loading, setLoading] = useState(false);
  const [dimensions, setDimensions] = useState({ width: DEFAULT_WIDTH, height: DEFAULT_HEIGHT });

  const currentHeight = selectedTerm ? 250 : 400;

  useEffect(() => {
    const node = containerRef.current;

    if (!node) {
      return;
    }

    const updateDimensions = () => {
      const nextWidth = Math.round(node.clientWidth) || DEFAULT_WIDTH;
      setDimensions({ width: nextWidth, height: currentHeight });
    };

    updateDimensions();

    if (typeof ResizeObserver === 'undefined') {
      return;
    }

    const observer = new ResizeObserver(() => {
      updateDimensions();
    });

    observer.observe(node);

    return () => {
      observer.disconnect();
    };
  }, [currentHeight]);

  const handleSelect = (term: string) => {
    setSelectedTerm((current) => (current === term ? null : term));
    setTooltip(null);
  };

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

  const words = useMemo(
    () =>
      buildWordCloudWords(
        items,
        (item) => item.trend_topic,
        (item) => item.outlier_ratio,
        (item, index) => ({ ...item, rank: index + 1 }),
      ),
    [items],
  );

  const measureWord = useMemo(() => createCanvasMeasureWord(), []);

  const positionedWords = useMemo(
    () =>
      layoutWordCloud({
        words,
        width: dimensions.width,
        height: dimensions.height,
        measureWord,
        getColor: (text) => getWordCloudColor(text),
      }),
    [dimensions.height, dimensions.width, measureWord, words],
  );

  const cloudBounds = useMemo(() => {
    if (positionedWords.length === 0) {
      return null;
    }

    return positionedWords.reduce(
      (acc, word) => ({
        left: Math.min(acc.left, word.x - word.width / 2),
        top: Math.min(acc.top, word.y - word.height / 2),
        right: Math.max(acc.right, word.x + word.width / 2),
        bottom: Math.max(acc.bottom, word.y + word.height / 2),
      }),
      { left: Infinity, top: Infinity, right: -Infinity, bottom: -Infinity },
    );
  }, [positionedWords]);

  const tooltipPosition = useMemo(() => {
    if (!tooltip || !containerRef.current || !cloudBounds) {
      return null;
    }

    const rect = containerRef.current.getBoundingClientRect();
    const cloudWidth = cloudBounds.right - cloudBounds.left;
    const cloudHeight = cloudBounds.bottom - cloudBounds.top;
    const cloudLeft = (rect.width - cloudWidth) / 2;
    const cloudTop = (rect.height - cloudHeight) / 2;

    return {
      left: rect.left + cloudLeft + (tooltip.x - cloudBounds.left),
      top: Math.max(rect.top + cloudTop + (tooltip.y - cloudBounds.top) - tooltip.height / 2 - 12, 12),
    };
  }, [cloudBounds, tooltip]);

  const isSelectedTermVisible = useMemo(() => {
    if (!selectedTerm) return false;
    return positionedWords.some((word) => word.text === selectedTerm);
  }, [positionedWords, selectedTerm]);

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
        ref={containerRef}
        className="tag-cloud"
        style={{ width: '100%', height: `${currentHeight}px`, position: 'relative', display: 'flex', alignItems: 'center', justifyContent: 'center' }}
      >
        {isSelectedTermVisible && selectedTerm && (
          <div style={{ position: 'absolute', top: '10px', right: '10px', zIndex: 10 }}>
            <button
              onClick={handleSearch}
              className="open-icon-btn"
              type="button"
              title={`${selectedTerm} ${BULLET_DELIMITER} ${domain.domain}`}
              style={{
                background: 'var(--bg-color)',
                border: '1px solid var(--border-color)',
                borderRadius: '4px',
                cursor: 'pointer',
                padding: '5px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: 'var(--text-color)',
                boxShadow: 'var(--shadow-medium)'
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

        {cloudBounds && (
          <div
            style={{
              position: 'relative',
              width: `${cloudBounds.right - cloudBounds.left}px`,
              height: `${cloudBounds.bottom - cloudBounds.top}px`,
            }}
          >
            {positionedWords.map((word) => {
          const isSelected = selectedTerm === word.text;
          const tooltipLabel = `${t('trends.rank', 'Rank')}: ${word.rank}. ${t('trends.trend', 'Trend')}: ${word.data.outlier_ratio.toFixed(2)}x. ${t('trends.vol', 'Vol')}: ${word.data.frequency}.`;

          return (
            <button
              key={word.text}
              type="button"
              className="word-cloud-item"
              aria-pressed={isSelected}
              aria-label={`${word.text}. ${tooltipLabel}`}
              onClick={(event) => {
                handleSelect(word.text);
                event.currentTarget.focus();
              }}
              onFocus={() => setTooltip(word)}
              onBlur={() => setTooltip((current) => (current?.text === word.text ? null : current))}
              onMouseEnter={() => setTooltip(word)}
              onMouseLeave={() => setTooltip((current) => (current?.text === word.text ? null : current))}
              onKeyDown={(event) => {
                if (event.key === 'Enter' || event.key === ' ') {
                  event.preventDefault();
                  handleSelect(word.text);
                }
              }}
              style={{
                position: 'absolute',
                left: `${word.x - cloudBounds.left}px`,
                top: `${word.y - cloudBounds.top}px`,
                transform: 'translate(-50%, -50%)',
                fontSize: `${word.fontSize}px`,
                color: isSelected ? 'var(--accent-color)' : word.color,
                fontWeight: isSelected ? 'bold' : 'normal',
                textDecoration: isSelected ? 'underline' : 'none',
                width: `${Math.ceil(word.width)}px`,
                height: `${Math.ceil(word.height)}px`,
              }}
            >
              {word.text}
            </button>
          );
            })}
          </div>
        )}

        {tooltip && (
          <div
            className="cloud-tooltip"
            style={{
              visibility: 'visible',
              opacity: 1,
              position: 'fixed',
              left: tooltipPosition?.left,
              top: tooltipPosition?.top,
              transform: 'translate(-50%, -100%)',
              marginTop: '-5px',
              pointerEvents: 'none',
              bottom: 'auto',
              zIndex: 2000,
            }}
          >
            {t('trends.rank', 'Rank')}: {tooltip.rank}<br/>
            {t('trends.trend', 'Trend')}: {tooltip.data.outlier_ratio.toFixed(2)}x<br/>
            {t('trends.vol', 'Vol')}: {tooltip.data.frequency}
          </div>
        )}
      </div>
      {isSelectedTermVisible && selectedTerm && (
        <TrendDetails term={selectedTerm} domain={domain} days={days} activeTab={activeTab} setActiveTab={setActiveTab} />
      )}
    </>
  );
};
