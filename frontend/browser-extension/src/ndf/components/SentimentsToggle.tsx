import { useTranslation } from 'react-i18next';

interface SentimentsToggleProps {
  sentimentType: 'sentiments' | 'sentiments_deframed';
  onTypeChange: (type: 'sentiments' | 'sentiments_deframed') => void;
  hasDeframed: boolean;
}

export const SentimentsToggle = ({ sentimentType, onTypeChange, hasDeframed }: SentimentsToggleProps) => {
  const { t } = useTranslation();
  return (
    <div className="sentiment-toggle-pill" role="group" aria-label={t('trends.sentiments_toggle', 'Sentiment Data Source')}>
      <button
        type="button"
        className={`sentiment-toggle-btn ${sentimentType === 'sentiments' ? 'active' : ''}`}
        aria-pressed={sentimentType === 'sentiments'}
        onClick={() => onTypeChange('sentiments')}
      >
        {t('trends.sentiments_original', 'Original')}
      </button>
      <button
        type="button"
        className={`sentiment-toggle-btn ${sentimentType === 'sentiments_deframed' ? 'active' : ''}`}
        aria-pressed={sentimentType === 'sentiments_deframed'}
        onClick={() => onTypeChange('sentiments_deframed')}
        disabled={!hasDeframed}
      >
        {t('trends.sentiments_deframed', 'Deframed')}
      </button>
    </div>
  );
};
