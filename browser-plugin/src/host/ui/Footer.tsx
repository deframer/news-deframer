import { useTranslation } from 'react-i18next';

export const Footer = () => {
  const { t } = useTranslation();

  return (
    <footer className="page-footer-text">
      <a
        href="https://deframer.github.io/"
        target="_blank"
        rel="noopener noreferrer"
      >
        {t('footer.github_link')}
      </a>
    </footer>
  );
};
