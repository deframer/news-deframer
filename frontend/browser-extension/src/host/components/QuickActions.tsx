import { useTranslation } from 'react-i18next';

import { FEEDBACK_EMAIL } from '../../shared/contact';

interface QuickActionsProps {
  onOpenSettings: () => void;
}

export const QuickActions = ({ onOpenSettings }: QuickActionsProps) => {
  const { t } = useTranslation();

  const handleProjectClick = () => {
    window.setTimeout(() => {
      window.close();
    }, 0);
  };

  return (
    <div className="card quick-actions-card">
      <div className="quick-actions">
        <button className="action-button action-button-enabled" onClick={onOpenSettings} type="button">
          {t('options.open_settings', 'Open Settings')}
        </button>
        <a
          className="action-button action-button-enabled quick-link"
          href="https://deframer.github.io/"
          target="_blank"
          rel="noopener noreferrer"
          onClick={handleProjectClick}
        >
          {t('footer.github_link')}
        </a>
        <a
          className="action-button action-button-enabled quick-link"
          href={`mailto:${FEEDBACK_EMAIL}`}
          onClick={handleProjectClick}
        >
          {FEEDBACK_EMAIL}
        </a>
      </div>
    </div>
  );
};
