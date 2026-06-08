import { useTranslation } from 'react-i18next';

import { NEWS_DEFRAMER_URL, REFERENCE_TAB_TARGET } from '../../shared/links';

export const Footer = () => {
  const { t } = useTranslation();
  return (
    <>
      <footer className="page-footer-text">
        {t('footer.text_pre', 'Visit')}
        <a
          href={NEWS_DEFRAMER_URL}
          target={REFERENCE_TAB_TARGET}
        >
          News Deframer
        </a>{' '}
        {t('footer.text_post', 'for more info.')}
      </footer>
    </>
  );
};
