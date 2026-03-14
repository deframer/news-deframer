import { CSSProperties, useEffect, useId, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import log from '../../shared/logger';
import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, SentimentItem, SentimentScores } from '../client';

interface TrendSentimentsProps {
  term: string;
  domain: DomainEntry;
  days: number;
  date?: string;
  className?: string;
}

interface SentimentMetric {
  key: string;
  label: string;
  description: string;
  min: number;
  max: number;
  baseline?: number;
  color: string;
}

const vadMetrics: SentimentMetric[] = [
  {
    key: 'valence',
    label: 'Valence',
    description: 'How pleasant or unpleasant the language is.',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--success-color)',
  },
  {
    key: 'arousal',
    label: 'Arousal',
    description: 'How calm or activated the language feels.',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--warning-color)',
  },
  {
    key: 'dominance',
    label: 'Dominance',
    description: 'How much control, strength, or power the language conveys.',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--accent-color)',
  },
];

const be5Metrics: SentimentMetric[] = [
  {
    key: 'joy',
    label: 'Joy',
    description: 'The strongest positive uplift signaled by the language.',
    min: 1,
    max: 5,
    color: 'var(--warning-color)',
  },
  {
    key: 'anger',
    label: 'Anger',
    description: 'The strongest hostility or outrage signaled by the language.',
    min: 1,
    max: 5,
    color: 'var(--danger-color)',
  },
  {
    key: 'sadness',
    label: 'Sadness',
    description: 'The strongest sorrow or loss signaled by the language.',
    min: 1,
    max: 5,
    color: 'var(--trend-steady)',
  },
  {
    key: 'fear',
    label: 'Fear',
    description: 'The strongest anxiety or threat signaled by the language.',
    min: 1,
    max: 5,
    color: 'var(--primary-color)',
  },
  {
    key: 'disgust',
    label: 'Disgust',
    description: 'The strongest revulsion or aversion signaled by the language.',
    min: 1,
    max: 5,
    color: 'var(--success-color)',
  },
];

const fallbackSentiments: Required<SentimentScores> = {
  valence: 5.14,
  arousal: 3.91,
  dominance: 5.3,
  joy: 3.48,
  anger: 2.49,
  sadness: 2.4,
  fear: 2.74,
  disgust: 2.59,
};

const formatValue = (value: number) => value.toFixed(2);

const getPercent = (value: number, min: number, max: number) => ((value - min) / (max - min)) * 100;

const getMetricValue = (value: number | null | undefined, fallback: number) => value ?? fallback;

const TOOLTIP_MAX_WIDTH = 280;
const TOOLTIP_GUTTER = 12;

export const TrendSentiments = ({ term, domain, days, date, className }: TrendSentimentsProps) => {
  const { t } = useTranslation();
  const tooltipId = useId();
  const [sentiments, setSentiments] = useState<SentimentItem | null>(null);
  const [loading, setLoading] = useState(true);
  const [hasLoaded, setHasLoaded] = useState(false);
  const [tooltipContent, setTooltipContent] = useState('');
  const [isTooltipOpen, setIsTooltipOpen] = useState(false);
  const [tooltipPos, setTooltipPos] = useState({ x: 0, y: 0 });

  useEffect(() => {
    let mounted = true;

    const fetchSentiments = async () => {
      setLoading(true);

      try {
        const settings = await getSettings();
        const client = new NewsDeframerClient(settings);
        const result = await client.getSentimentsByTrend(domain.domain, term, date, days);
        log.debug('Sentiment API result:', result);

        if (mounted) {
          setSentiments(result);
          setHasLoaded(true);
        }
      } catch (error) {
        console.error('Failed to fetch sentiment data', error);
        if (mounted) {
          setSentiments(null);
          setHasLoaded(true);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    fetchSentiments();

    return () => {
      mounted = false;
    };
  }, [date, days, domain.domain, term]);

  const metricValues = useMemo(() => ({
    valence: getMetricValue(sentiments?.sentiments?.valence, fallbackSentiments.valence),
    arousal: getMetricValue(sentiments?.sentiments?.arousal, fallbackSentiments.arousal),
    dominance: getMetricValue(sentiments?.sentiments?.dominance, fallbackSentiments.dominance),
    joy: getMetricValue(sentiments?.sentiments?.joy, fallbackSentiments.joy),
    anger: getMetricValue(sentiments?.sentiments?.anger, fallbackSentiments.anger),
    sadness: getMetricValue(sentiments?.sentiments?.sadness, fallbackSentiments.sadness),
    fear: getMetricValue(sentiments?.sentiments?.fear, fallbackSentiments.fear),
    disgust: getMetricValue(sentiments?.sentiments?.disgust, fallbackSentiments.disgust),
  }), [sentiments]);

  const showTooltip = (target: HTMLElement, description: string) => {
    const rect = target.getBoundingClientRect();
    const preferredX = rect.left + rect.width / 2;
    const tooltipWidth = Math.min(TOOLTIP_MAX_WIDTH, window.innerWidth * 0.6, window.innerWidth - TOOLTIP_GUTTER * 2);
    const minX = tooltipWidth / 2 + TOOLTIP_GUTTER;
    const maxX = window.innerWidth - tooltipWidth / 2 - TOOLTIP_GUTTER;
    const clampedX = Math.max(minX, Math.min(preferredX, maxX));
    setTooltipPos({ x: clampedX, y: rect.top - 10 });
    setTooltipContent(description);
    setIsTooltipOpen(true);
  };

  const hideTooltip = () => {
    setIsTooltipOpen(false);
  };

  if (loading) {
    return (
      <div className={`sentiment-panel ${className || ''}`}>
        <div className="sentiment-panel-state">
          <div className="spinner-small" />
        </div>
      </div>
    );
  }

  if (hasLoaded && (!sentiments || !sentiments.sentiments)) {
    return (
      <div className={`sentiment-panel ${className || ''}`}>
          <div className="sentiment-panel-state sentiment-panel-empty">
          {t('trends.sentiments_no_data', 'No sentiment data available for this topic.')}
        </div>
      </div>
    );
  }

  const renderInfoButton = (description: string) => (
    <div className="article-header-info-trigger-wrap">
      <button
        type="button"
        className="article-header-info-trigger"
        aria-label={description}
        aria-describedby={tooltipId}
        onMouseEnter={(e) => showTooltip(e.currentTarget, description)}
        onMouseLeave={hideTooltip}
        onFocus={(e) => showTooltip(e.currentTarget, description)}
        onBlur={hideTooltip}
      >
        i
      </button>
    </div>
  );

  return (
    <div className={`sentiment-panel ${className || ''}`}>
      <section className="sentiment-section">
        <div className="sentiment-section-header">
          <strong>{t('trends.sentiments_vad', 'VAD')}</strong>
          <span>{t('trends.sentiments_vad_scale', 'Scale 1-9, neutral at 5')}</span>
        </div>
        <div className="sentiment-rows">
          {vadMetrics.map((metric) => {
            const baselinePercent = getPercent(metric.baseline || metric.min, metric.min, metric.max);
            const value = metricValues[metric.key as keyof typeof metricValues];
            const valuePercent = getPercent(value, metric.min, metric.max);
            const fillStyle: CSSProperties = {
              background: metric.color,
              left: `${Math.min(baselinePercent, valuePercent)}%`,
              width: `${Math.abs(valuePercent - baselinePercent)}%`,
            };

            return (
              <div key={metric.key} className="sentiment-metric">
                <div className="sentiment-metric-row">
                  <span className="sentiment-label">{t(`trends.sentiments_metrics.${metric.key}.label`, metric.label)}</span>
                  <span className="sentiment-value">{formatValue(value)}</span>
                  <div className="sentiment-track-wrap">
                    <div className="sentiment-vad-track">
                      <div className="sentiment-vad-fill" style={fillStyle} />
                      <div className="sentiment-vad-baseline" style={{ left: `${baselinePercent}%` }} />
                    </div>
                  </div>
                  {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`, metric.description))}
                </div>
              </div>
            );
          })}
        </div>
      </section>

      <section className="sentiment-section">
        <div className="sentiment-section-header">
          <strong>{t('trends.sentiments_be5', 'BE5')}</strong>
          <span>{t('trends.sentiments_be5_scale', 'Scale 1-5, absent to max')}</span>
        </div>
        <div className="sentiment-rows">
          {be5Metrics.map((metric) => {
            const value = metricValues[metric.key as keyof typeof metricValues];

            return (
              <div key={metric.key} className="sentiment-metric">
                <div className="sentiment-metric-row">
                  <span className="sentiment-label">{t(`trends.sentiments_metrics.${metric.key}.label`, metric.label)}</span>
                  <span className="sentiment-value">{formatValue(value)}</span>
                  <div className="sentiment-track-wrap">
                    <div className="sentiment-be5-grid">
                      {[1, 2, 3, 4, 5].map((step) => {
                        const fillWidth = Math.max(0, Math.min(1, value - (step - 1))) * 100;
                        return (
                          <div key={`${metric.key}-${step}`} className="sentiment-be5-box">
                            <div className="sentiment-be5-fill" style={{ width: `${fillWidth}%`, background: metric.color }} />
                          </div>
                        );
                      })}
                    </div>
                  </div>
                  {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`, metric.description))}
                </div>
              </div>
            );
          })}
        </div>
      </section>
      {isTooltipOpen && (
        <div
          id={tooltipId}
          role="tooltip"
          className="article-header-floating-tooltip"
          style={{ left: tooltipPos.x, top: tooltipPos.y }}
        >
          {tooltipContent}
        </div>
      )}
      <div className="sentiment-footnote">
        {t('trends.sentiments_disclaimer', 'These scores describe the emotional meaning of the selected article set. They are lexical estimates based on word usage, not direct measurements of what authors, readers, or the public feel.')}
      </div>
    </div>
  );
};
