export const createFooterHtml = (): string => {
  return `
    <footer class="page-footer-text">
      This content was replaced by the <a href="https://github.com/egandro/news-deframer" target="_blank" rel="noopener noreferrer">News Deframer</a> browser plugin.
    </footer>
  `;
};

export const getFooterCss = (): string => {
  return `
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
};