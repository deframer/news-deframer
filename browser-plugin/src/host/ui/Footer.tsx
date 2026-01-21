import React from 'react';

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
  return (
    <>
      <style>{footerCss}</style>
      <footer className="page-footer-text">
        <a
          href="https://github.com/egandro/news-deframer"
          target="_blank"
          rel="noopener noreferrer"
        >
          News Deframer on GitHub
        </a>
      </footer>
    </>
  );
};
