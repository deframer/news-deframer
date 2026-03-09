import { useTranslation } from 'react-i18next';

export const Footer = () => {
  const { t } = useTranslation();
  return (
    <>
      <footer className="page-footer-text">
        {t('footer.text_pre', 'Visit')}
        <a
          href="https://deframer.github.io/"
          target="_blank"
          rel="noopener noreferrer"
        >
          News Deframer
        </a>{' '}
        {t('footer.text_post', 'for more info.')}
      </footer>
    </>
  );
};
