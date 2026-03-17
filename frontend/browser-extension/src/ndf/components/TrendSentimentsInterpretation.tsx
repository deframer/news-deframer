import { useTranslation } from 'react-i18next';

import { AnalysisOutput, getSentimentUIText } from '../../shared/sentiments';

interface TrendSentimentsInterpretationProps {
  interpretation: AnalysisOutput;
}

export const TrendSentimentsInterpretation = ({ interpretation }: TrendSentimentsInterpretationProps) => {
  const { t, i18n } = useTranslation();
  const lang = i18n.language || 'en';

  const interpretationMetrics = [
    { label: getSentimentUIText('sentiments_interpretation', lang), value: interpretation.interpretation },
    { label: getSentimentUIText('sentiments_core_state', lang), value: interpretation.core_state },
    { label: getSentimentUIText('sentiments_emotions', lang), value: `${interpretation.primary_emotion} / ${interpretation.secondary_emotion}` },
    { label: getSentimentUIText('sentiments_tension', lang), value: interpretation.tension_label },
    { label: getSentimentUIText('sentiments_control', lang), value: interpretation.control_label },
    { label: getSentimentUIText('sentiments_mood', lang), value: interpretation.mood_label },
    { label: getSentimentUIText('sentiments_clarity', lang), value: interpretation.clarity_label },
  ];

  return (
    <div className="sentiment-interpretation-container">
      <div className="sentiment-interpretation-header-wrap">
        <strong className="sentiment-interpretation-header-label">{t('trends.sentiments_interpretation_header', 'Interpretation')}</strong>
      </div>
      <div className="sentiment-interpretation-list">
        {interpretationMetrics.map((metric, index) => (
          <div key={index} className="sentiment-interpretation-item">
            <span className="sentiment-interpretation-label">{metric.label}</span>
            <span className="sentiment-interpretation-value">{metric.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
};
