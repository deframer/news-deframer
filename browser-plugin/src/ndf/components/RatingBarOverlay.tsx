import { useTranslation } from 'react-i18next';

import { formatRatingPercent, getRatingColors } from '../utils/rating';

const overlayCss = `
  .bar-container {
    background-color: var(--rating-bg);
    height: 30px;
    width: 100%;
    position: relative; /* For overlay and tooltip */
  }
  .bar {
    height: 100%;
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
    padding-left: 10px;
    font-weight: bold;
    font-size: 0.9em;
    pointer-events: none;
  }
  .bar-container .tooltip-text {
    visibility: hidden;
    opacity: 0;
    width: 250px;
    background-color: var(--tooltip-bg);
    color: var(--tooltip-text);
    text-align: left;
    border-radius: 6px;
    padding: 10px;
    position: absolute;
    z-index: 1;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    transition: opacity 0.2s;
    font-size: 0.9em;
    line-height: 1.4;
    pointer-events: none;
  }
  .bar-container:hover .tooltip-text {
    visibility: visible;
    opacity: 1;
  }
`;

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
      <style>{overlayCss}</style>
      <div
        className="bar-container"
        role="meter"
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={percentage}
        aria-label={t('rating.overlay_aria_label', { percentage, reason: reason || t('rating.no_reason') })}
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
        <div className="tooltip-text">{reason || t('rating.no_reason')}</div>
      </div>
    </div>
  );
};
