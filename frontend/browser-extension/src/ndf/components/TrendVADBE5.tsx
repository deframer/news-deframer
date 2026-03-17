import { CSSProperties } from 'react';
import { useTranslation } from 'react-i18next';

import { SentimentScores } from '../client';

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

interface TrendVADBE5Props {
  metricValues: SentimentScores;
  tooltipId: string;
  onShowTooltip: (target: HTMLElement, description: string) => void;
  onHideTooltip: () => void;
}

export const TrendVADBE5 = ({ metricValues, tooltipId, onShowTooltip, onHideTooltip }: TrendVADBE5Props) => {
  const { t } = useTranslation();

  const renderInfoButton = (description: string) => (
    <div className="article-header-info-trigger-wrap">
      <button
        type="button"
        className="article-header-info-trigger"
        aria-label={description}
        aria-describedby={tooltipId}
        onMouseEnter={(e) => onShowTooltip(e.currentTarget, description)}
        onMouseLeave={onHideTooltip}
        onFocus={(e) => onShowTooltip(e.currentTarget, description)}
        onBlur={onHideTooltip}
      >
        i
      </button>
    </div>
  );

  return (
    <>
      <div className="sentiment-section-header">
        <strong className="sentiment-header-label">{t('trends.sentiments_vad', 'VAD')}</strong>
        <span className="sentiment-header-scale">{t('trends.sentiments_vad_scale', 'Scale 1-9, neutral at 5')}</span>
      </div>
      <div className="sentiment-rows">
        {vadMetrics.map((metric) => {
          const value = metricValues[metric.key as keyof typeof metricValues];

          if (value === null || value === undefined) {
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

      <div className="sentiment-section-header">
        <strong className="sentiment-header-label">{t('trends.sentiments_be5', 'BE5')}</strong>
        <span className="sentiment-header-scale">{t('trends.sentiments_be5_scale', 'Scale 1-5, absent to max')}</span>
      </div>
      <div className="sentiment-rows">
        {be5Metrics.map((metric) => {
          const value = metricValues[metric.key as keyof typeof metricValues];

          if (value === null || value === undefined) {
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
    </>
  );
};
