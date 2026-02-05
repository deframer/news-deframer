export type Theme = 'light' | 'dark' | 'system';

export const cssVariables = {
  light: `
    --bg-color: #f0f2f5;
    --text-color: #333333;
    --card-bg: #ffffff;
    --card-shadow: 0 4px 6px rgba(0,0,0,0.1);
    --card-shadow-hover: 0 8px 12px rgba(0,0,0,0.15);
    --header-bg: #ffffff;
    --border-color: #eeeeee;
    --secondary-text: #666666;
    --btn-bg: #ffffff;
    --btn-text: #333333;
    --btn-border: #dddddd;
    --btn-hover-bg: #f8f9fa;
    --accent-color: #0056b3;
    --accent-hover: #004494;
    --accent-text: #ffffff;
    --rating-bg: #e9ecef;
    --success-color: #198754;
    --warning-color: #ffc107;
    --danger-color: #b02a37;
    --tooltip-bg: rgba(0,0,0,0.85);
    --tooltip-text: #ffffff;
    --bg-color-secondary: #f7f8fa;
    --hover-bg: rgba(0,0,0,0.05);
    --badge-bg: #e9ecef;
    --primary-color: #0056b3;
    --primary-color-dark: #004494;
  `,
  dark: `
    --bg-color: #18191a;
    --text-color: #e4e6eb;
    --card-bg: #242526;
    --card-shadow: 0 4px 6px rgba(0,0,0,0.3);
    --card-shadow-hover: 0 8px 12px rgba(0,0,0,0.4);
    --header-bg: #242526;
    --border-color: #3e4042;
    --secondary-text: #b0b3b8;
    --btn-bg: #3a3b3c;
    --btn-text: #e4e6eb;
    --btn-border: #3e4042;
    --btn-hover-bg: #4e4f50;
    --accent-color: #2d88ff;
    --accent-hover: #4599ff;
     --accent-text: #000000;
     --rating-bg: #3a3b3c;
     --success-color: #4caf50;
     --warning-color: #ffca28;
     --danger-color: #f44336;
     --tooltip-bg: rgba(255,255,255,0.9);
     --tooltip-text: #000000;
     --bg-color-secondary: #18191a;
     --hover-bg: rgba(255,255,255,0.05);
     --badge-bg: #3a3b3c;
     --primary-color: #2d88ff;
     --primary-color-dark: #4599ff;
   `
};

export const globalStyles = `
  *, *::before, *::after {
    box-sizing: border-box;
  }
  body {
    margin: 0;
    padding: 0;
    background-color: var(--bg-color);
    color: var(--text-color);
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    min-height: 100vh;
    width: 100%;
  }
`;

export const getThemeCss = (theme: Theme): string => {
  if (theme === 'system') {
    return `
      :host, :root {
        ${cssVariables.light}
      }
      @media (prefers-color-scheme: dark) {
        :host, :root {
          ${cssVariables.dark}
        }
      }
    `;
  }
  return `
    :host, :root {
      ${cssVariables[theme]}
    }
  `;
};
