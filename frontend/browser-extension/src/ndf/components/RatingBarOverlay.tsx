import { useTranslation } from 'react-i18next';

import { formatRatingPercent, getRatingColors } from '../utils/rating';

interface RatingBarOverlayProps {
  value: number | undefined;
  reason?: string;
  className?: string;
}

export const RatingBarOverlay = ({
  value,
  reason,
  className,
}: RatingBarOverlayProps) => {
  const { t } = useTranslation();
  const percentage = formatRatingPercent(value);
  const colors = getRatingColors(percentage);
  const textShadow =
    colors.text === '#ffffff' ? '0 0 3px rgba(0,0,0,0.7)' : '';

  return (
    <div className={className}>
      <div
        className="bar-container"
        role="meter"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={percentage}
        aria-label={t('rating.overlay_aria_label', 'Rating: {{percentage}}%. {{reason}}', { percentage, reason: reason || t('rating.no_reason', 'No reason provided') })}
      >
        <div
          className="bar"
          style={{ width: `${percentage}%`, backgroundColor: colors.bg }}
        />
        <div
          className="bar-overlay"
          style={{ color: colors.text, textShadow }}
        >
          <span>{percentage}%</span>
        </div>
        <div className="tooltip-text">{reason || t('rating.no_reason', 'No reason provided')}</div>
      </div>
    </div>
  );
};
