import { useTranslation } from 'react-i18next';

import { FEEDBACK_EMAIL } from '../../shared/contact';

export const SettingsAbout = () => {
  const { t } = useTranslation();

  return (
    <div className="content-grid">
      <div className="settings-column">
        <div className="card">
          <h3 className="section-title">{t('options.section_project')}</h3>
          <a
            href="https://deframer.github.io/"
            target="_blank"
            rel="noopener noreferrer"
            className="project-link"
          >
            {t('footer.github_link')}
          </a>
        </div>
      </div>
      <div className="settings-column">
        <div className="card">
          <h3 className="section-title">{t('options.section_feedback')}</h3>
          <a href={`mailto:${FEEDBACK_EMAIL}`} className="project-link">
            {FEEDBACK_EMAIL}
          </a>
        </div>
      </div>
    </div>
  );
};
