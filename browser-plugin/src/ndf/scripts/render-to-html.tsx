import fs from 'fs';
import path from 'path';
import React from 'react';
import ReactDOMServer from 'react-dom/server';

// 1. Mock Browser Environment
// This must be done before importing components that might access window at module level
const mockWindow = {
  location: {
    hostname: 'example.com',
    href: 'http://example.com/article',
    search: '',
    pathname: '/',
    reload: () => console.log('window.location.reload called'),
  },
  scrollTo: () => {},
  sessionStorage: {
    getItem: () => null,
    setItem: () => {},
  },
};

(global as any).window = mockWindow;
(global as any).document = {
  createElement: () => ({}),
  cookie: '',
};

if (!('navigator' in global)) {
  (global as any).navigator = {
    userAgent: 'node',
  };
}

// 2. Helper to load components
// We use require() here to ensure the global mocks are applied before the modules are evaluated
const loadComponents = () => {
  const { ArticlePage } = require('../pages/ArticlePage');
  const { PortalPage } = require('../pages/PortalPage');
  const themeModule = require('../../shared/theme');
  const { getThemeCss, globalStyles, Theme } = themeModule.default || themeModule;
  return { ArticlePage, PortalPage, getThemeCss, globalStyles, Theme };
};

const render = () => {
  const args = process.argv.slice(2);
  if (args.length < 1) {
    console.error('Usage: npx ts-node src/ndf/scripts/render-to-html.tsx <path-to-json> [type: article|portal]');
    process.exit(1);
  }

  const jsonPath = args[0];
  const type = args[1] || 'article'; // default to article

  // Read and parse JSON
  const absolutePath = path.isAbsolute(jsonPath) ? jsonPath : path.join(process.cwd(), jsonPath);
  if (!fs.existsSync(absolutePath)) {
    console.error(`File not found: ${absolutePath}`);
    process.exit(1);
  }

  const jsonData = JSON.parse(fs.readFileSync(absolutePath, 'utf-8'));
  const { ArticlePage, PortalPage, getThemeCss, globalStyles, Theme } = loadComponents();

  let componentHtml = '';

  try {
    if (type === 'portal') {
      // Detect if JSON is array or object wrapping array
      const items = Array.isArray(jsonData) ? jsonData : jsonData.items;
      componentHtml = ReactDOMServer.renderToStaticMarkup(<PortalPage items={items} />);
    } else {
      componentHtml = ReactDOMServer.renderToStaticMarkup(<ArticlePage item={jsonData} />);
    }
  } catch (error) {
    console.error('Error rendering component:', error);
    process.exit(1);
  }

  // Generate CSS from theme (using LIGHT theme by default for debug)
  const themeValue = (Theme && Theme.LIGHT) || 'light';
  const css = getThemeCss(themeValue) + globalStyles;

  // 3. Wrap in HTML with CSS Variables for correct styling
  const fullHtml = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Debug Render - ${type}</title>
  <style>
    ${css}
  </style>
</head>
<body>
  ${componentHtml}
</body>
</html>`;

  console.log(fullHtml);
};

render();