import { useTranslation } from 'react-i18next';

import { formatRatingPercent, getRatingColors } from '../utils/rating';

interface RatingBarProps {
  value: number | undefined;
  label?: string;
  id?: string;
  reason?: string;
}

export const RatingBar = ({ value, label, id, reason }: RatingBarProps) => {
  const { t } = useTranslation();
  const percentage = formatRatingPercent(value);
  const colors = getRatingColors(percentage);
  const textShadow =
    colors.text === '#ffffff' ? '0 0 3px rgba(0,0,0,0.7)' : '';

  return (
    <>
      {label && <div className="metric-label">{label}</div>}
      <div className="metric-details">
        <div
          className="bar-container"
          role="meter"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={percentage}
          aria-label={t('rating.aria_label', '{{label}}: {{percentage}}%', { label: label || t('rating.label', 'Rating'), percentage })}
          aria-describedby={id ? `${id}-reason` : undefined}
        >
          <div
            className="bar"
            style={{ width: `${percentage}%`, backgroundColor: colors.bg }}
          />
          {value !== undefined && (
            <div className="bar-overlay" style={{ color: colors.text, textShadow }}>
              <span>{percentage}%</span>
            </div>
          )}
        </div>
        {reason && (
          <p id={id ? `${id}-reason` : undefined} className="reason">
            {reason}
          </p>
        )}
      </div>
    </>
  );
};
