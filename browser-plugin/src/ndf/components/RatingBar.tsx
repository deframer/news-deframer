import { useTranslation } from 'react-i18next';

import { formatRatingPercent, getRatingColors } from '../utils/rating';

const ratingBarCss = `
  .bar-container {
    background-color: var(--rating-bg);
    border-radius: 5px;
    height: 30px;
    width: 100%;
    position: relative;
    overflow: hidden;
  }
  .bar {
    height: 100%;
    border-radius: 5px;
    transition: width 0.3s ease;
  }
  .bar-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    padding-left: 8px;
    font-weight: bold;
    font-size: 0.95em;
    pointer-events: none;
  }
`;

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
      <style>{ratingBarCss}</style>
      {label && <div className="metric-label">{label}</div>}
      <div className="metric-details">
        <div
          className="bar-container"
          role="meter"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={percentage}
          aria-label={t('rating.aria_label', { label: label || t('rating.label'), percentage })}
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
