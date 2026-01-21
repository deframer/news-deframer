import React from 'react';

const footerCss = `
  .page-footer-text {
    text-align: center;
    padding: 1.5em;
    color: #6c757d;
    font-size: 0.9em;
    line-height: 1.5;
  }
  .page-footer-text a {
    color: #495057;
    text-decoration: underline;
  }
  .page-footer-text a:hover {
    color: #212529;
  }
`;

export const Footer = () => {
  return (
    <>
      <style>{footerCss}</style>
      <footer className="page-footer-text">
        This content was replaced by the{' '}
        <a
          href="https://github.com/egandro/news-deframer"
          target="_blank"
          rel="noopener noreferrer"
        >
          News Deframer
        </a>{' '}
        browser plugin.
      </footer>
    </>
  );
};
