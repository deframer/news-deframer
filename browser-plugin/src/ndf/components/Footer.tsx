import React from 'react';
import { useTranslation } from 'react-i18next';

const footerCss = `
  .page-footer-text {
    text-align: center;
    padding: 1.5em;
    color: var(--secondary-text);
    font-size: 0.9em;
    line-height: 1.5;
  }
  .page-footer-text a {
    color: var(--secondary-text);
    text-decoration: underline;
  }
  .page-footer-text a:hover {
    color: var(--text-color);
  }
`;

export const Footer = () => {
  const { t } = useTranslation();
  return (
    <>
      <style>{footerCss}</style>
      <footer className="page-footer-text">
        {t('footer.text_pre')}
        <a
          href="https://deframer.github.io/"
          target="_blank"
          rel="noopener noreferrer"
        >
          News Deframer
        </a>{' '}
        {t('footer.text_post')}
      </footer>
    </>
  );
};
