import { CSSProperties, useEffect, useId, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import '../styles.css';


import log from '../../shared/logger';
import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, SentimentItem } from '../client';

interface TrendSentimentsProps {
  term: string;
  domain: DomainEntry;
  days: number;
  date?: string;
  className?: string;
}

interface SentimentMetric {
  key: string;
  min: number;
  max: number;
  baseline?: number;
  color: string;
}

const vadMetrics: SentimentMetric[] = [
  {
    key: 'valence',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--success-color)',
  },
  {
    key: 'arousal',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--warning-color)',
  },
  {
    key: 'dominance',
    min: 1,
    max: 9,
    baseline: 5,
    color: 'var(--accent-color)',
  },
];

const be5Metrics: SentimentMetric[] = [
  {
    key: 'joy',
    min: 1,
    max: 5,
    color: 'var(--warning-color)',
  },
  {
    key: 'anger',
    min: 1,
    max: 5,
    color: 'var(--danger-color)',
  },
  {
    key: 'sadness',
    min: 1,
    max: 5,
    color: 'var(--trend-steady)',
  },
  {
    key: 'fear',
    min: 1,
    max: 5,
    color: 'var(--primary-color)',
  },
  {
    key: 'disgust',
    min: 1,
    max: 5,
    color: 'var(--success-color)',
  },
];

const formatValue = (value: number) => value.toFixed(2);

const getPercent = (value: number, min: number, max: number) => ((value - min) / (max - min)) * 100;

const TOOLTIP_MAX_WIDTH = 280;
const TOOLTIP_GUTTER = 12;

export const TrendSentiments = ({ term, domain, days, date, className }: TrendSentimentsProps) => {
  const { t } = useTranslation();
  const tooltipId = useId();
  const [sentiments, setSentiments] = useState<SentimentItem | null>(null);
  const [sentimentType, setSentimentType] = useState<'sentiments' | 'sentiments_deframed'>('sentiments');
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

  const metricValues = useMemo(() => {
    const selectedSentiments = sentiments?.[sentimentType];
    return {
      valence: selectedSentiments?.valence ?? null,
      arousal: selectedSentiments?.arousal ?? null,
      dominance: selectedSentiments?.dominance ?? null,
      joy: selectedSentiments?.joy ?? null,
      anger: selectedSentiments?.anger ?? null,
      sadness: selectedSentiments?.sadness ?? null,
      fear: selectedSentiments?.fear ?? null,
      disgust: selectedSentiments?.disgust ?? null,
    };
  }, [sentiments, sentimentType]);

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

  const selectedSentimentData = sentiments?.[sentimentType];
  const hasData = selectedSentimentData && Object.keys(selectedSentimentData).length > 0;

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
      {!hasData && hasLoaded && (
        <div className="sentiment-panel-state sentiment-panel-empty">
          {t('trends.sentiments_no_data_for_type', 'No {{type}} sentiment data available.', {
            type: t(`trends.sentiments_${sentimentType}`, sentimentType),
          })}
        </div>
      )}

      {hasData && (
        <section className="sentiment-section">
          <div className="sentiment-toggle-pill" role="group" aria-label={t('trends.sentiments_toggle', 'Sentiment Data Source')}>
            <button
              type="button"
              className={`sentiment-toggle-btn ${sentimentType === 'sentiments' ? 'active' : ''}`}
              aria-pressed={sentimentType === 'sentiments'}
              onClick={() => setSentimentType('sentiments')}
            >
              {t('trends.sentiments_original', 'Original')}
            </button>
            <button
              type="button"
              className={`sentiment-toggle-btn ${sentimentType === 'sentiments_deframed' ? 'active' : ''}`}
              aria-pressed={sentimentType === 'sentiments_deframed'}
              onClick={() => setSentimentType('sentiments_deframed')}
              disabled={!sentiments?.sentiments_deframed}
            >
              {t('trends.sentiments_deframed', 'Deframed')}
            </button>
          </div>

          <div className="sentiment-section-header">
            <strong className="sentiment-header-label">{t('trends.sentiments_vad', 'VAD')}</strong>
            <span className="sentiment-header-scale">{t('trends.sentiments_vad_scale', 'Scale 1-9, neutral at 5')}</span>
          </div>
          <div className="sentiment-rows">
            {vadMetrics.map((metric) => {
              const value = metricValues[metric.key as keyof typeof metricValues];

              if (value === null) {
                return (
                  <div key={metric.key} className="sentiment-metric">
                    <div className="sentiment-metric-row">
                      <span className="sentiment-label">
                        {t(`trends.sentiments_metrics.${metric.key}.label`)}
                      </span>
                      <span className="sentiment-value">{t('trends.sentiments_na', 'N/A')}</span>
                      <div className="sentiment-track-wrap" />
                      {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`))}
                    </div>
                  </div>
                );
              }

              const baselinePercent = getPercent(metric.baseline || metric.min, metric.min, metric.max);
              const valuePercent = getPercent(value, metric.min, metric.max);
              const fillStyle: CSSProperties = {
                background: metric.color,
                left: `${Math.min(baselinePercent, valuePercent)}%`,
                width: `${Math.abs(valuePercent - baselinePercent)}%`,
              };

              return (
                <div key={metric.key} className="sentiment-metric">
                  <div className="sentiment-metric-row">
                    <span className="sentiment-label">
                      {t(`trends.sentiments_metrics.${metric.key}.label`)}
                    </span>
                    <span className="sentiment-value">{formatValue(value)}</span>
                    <div className="sentiment-track-wrap">
                      <div className="sentiment-vad-track">
                        <div className="sentiment-vad-fill" style={fillStyle} />
                        <div className="sentiment-vad-baseline" style={{ left: `${baselinePercent}%` }} />
                      </div>
                    </div>
                    {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`))}
                  </div>
                </div>
              );
            })}
          </div>

          <hr className="sentiment-divider" />

          <div className="sentiment-section-header">
            <strong className="sentiment-header-label">{t('trends.sentiments_be5', 'BE5')}</strong>
            <span className="sentiment-header-scale">{t('trends.sentiments_be5_scale', 'Scale 1-5, absent to max')}</span>
          </div>
          <div className="sentiment-rows">
            {be5Metrics.map((metric) => {
              const value = metricValues[metric.key as keyof typeof metricValues];

              if (value === null) {
                return (
                  <div key={metric.key} className="sentiment-metric">
                    <div className="sentiment-metric-row">
                      <span className="sentiment-label">
                        {t(`trends.sentiments_metrics.${metric.key}.label`)}
                      </span>
                      <span className="sentiment-value">{t('trends.sentiments_na', 'N/A')}</span>
                      <div className="sentiment-track-wrap" />
                      {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`))}
                    </div>
                  </div>
                );
              }

              return (
                <div key={metric.key} className="sentiment-metric">
                  <div className="sentiment-metric-row">
                    <span className="sentiment-label">
                      {t(`trends.sentiments_metrics.${metric.key}.label`)}
                    </span>
                    <span className="sentiment-value">{formatValue(value)}</span>
                    <div className="sentiment-track-wrap">
                      <div className="sentiment-be5-grid">
                        {[1, 2, 3, 4, 5].map((step) => {
                          const fillWidth = Math.max(0, Math.min(1, value - (step - 1))) * 100;
                          return (
                            <div key={`${metric.key}-${step}`} className="sentiment-be5-box">
                              <div
                                className="sentiment-be5-fill"
                                style={{ width: `${fillWidth}%`, background: metric.color }}
                              />
                            </div>
                          );
                        })}
                      </div>
                    </div>
                    {renderInfoButton(t(`trends.sentiments_metrics.${metric.key}.description`))}
                  </div>
                </div>
              );
            })}
          </div>
        </section>
      )}

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
        {t(
          'trends.sentiments_disclaimer',
          'These scores describe the emotional meaning of the selected article set. They are lexical estimates based on word usage, not direct measurements of what authors, readers, or the public feel.'
        )}
      </div>
    </div>
  );
};
