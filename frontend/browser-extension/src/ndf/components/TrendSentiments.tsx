import { useEffect, useId, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import log from '../../shared/logger';
import type { AnalysisOutput, EmotionVector } from '../../shared/sentiments';
import { sentimentsToLabels } from '../../shared/sentiments';
import { getSettings } from '../../shared/settings';
import { DomainEntry, NewsDeframerClient, SentimentItem, SentimentScores } from '../client';
import { TrendSentimentsInterpretation } from './TrendSentimentsInterpretation';
import { TrendVADBE5 } from './TrendVADBE5';

interface TrendSentimentsProps {
  term: string;
  domain: DomainEntry;
  days: number;
  date?: string;
  className?: string;
}

const TOOLTIP_MAX_WIDTH = 280;
const TOOLTIP_GUTTER = 12;

export const TrendSentiments = ({ term, domain, days, date, className }: TrendSentimentsProps) => {
  const { t, i18n } = useTranslation();
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

  const metricValues = useMemo((): SentimentScores => {
    return sentiments?.[sentimentType] || { valence: 5, arousal: 5, dominance: 5, joy: 1, anger: 1, sadness: 1, fear: 1, disgust: 1 };
  }, [sentiments, sentimentType]);

  const interpretation = useMemo((): AnalysisOutput | null => {
    const ev: EmotionVector = {
      valence: metricValues.valence,
      arousal: metricValues.arousal,
      dominance: metricValues.dominance,
      joy: metricValues.joy,
      anger: metricValues.anger,
      sadness: metricValues.sadness,
      fear: metricValues.fear,
      disgust: metricValues.disgust,
    };
    const lang = i18n.language || 'en';
    return sentimentsToLabels(ev, lang) as AnalysisOutput;
  }, [metricValues, i18n.language]);

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

  return (
    <div className={`sentiment-panel ${className || ''}`}>
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

      {!hasData && hasLoaded && (
        <div className="sentiment-panel-state sentiment-panel-empty">
          {t('trends.sentiments_no_data_for_type', 'No {{type}} sentiment data available.', {
            type: t(`trends.sentiments_${sentimentType}`, sentimentType),
          })}
        </div>
      )}

      {hasData && (
        <section className="sentiment-columns-layout">
          <div className="sentiment-column-left">
            <TrendVADBE5
              metricValues={metricValues}
              tooltipId={tooltipId}
              onShowTooltip={showTooltip}
              onHideTooltip={hideTooltip}
            />
          </div>

          <div className="sentiment-vertical-divider" />

          <div className="sentiment-column-right">
            {interpretation && (
              <TrendSentimentsInterpretation interpretation={interpretation} />
            )}
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
