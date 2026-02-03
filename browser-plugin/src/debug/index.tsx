import '../shared/i18n';

import { createRoot } from 'react-dom/client';

import { AnalyzedItem } from '../ndf/client';
import { ArticlePage } from '../ndf/pages/ArticlePage';
import { PortalPage } from '../ndf/pages/PortalPage';
import log from '../shared/logger';
import { getThemeCss, globalStyles, Theme } from '../shared/theme';

//log.setLevel('trace');
log.setLevel('debug');

// Import data (using require to ensure it loads without strict TS json config)
// eslint-disable-next-line @typescript-eslint/no-require-imports
const portalData = require('./portal.json');

// Inject Theme CSS
const style = document.createElement('style');
style.textContent = getThemeCss('light' as Theme) + globalStyles;
document.head.appendChild(style);

const container = document.getElementById('root');
if (container) {
  const root = createRoot(container);
  const params = new URLSearchParams(window.location.search);
  const type = params.get('type');

  // Mock window.location for the components if needed,
  // though they run in a real browser now so it's less critical.

  if (type === 'article') {
    const rawItems = Array.isArray(portalData) ? portalData : (portalData as { items: AnalyzedItem[] }).items;
    const i = params.get('i');
    let item = rawItems[0] as AnalyzedItem;
    if (i) {
      const idx = parseInt(i, 10);
      if (rawItems && rawItems[idx]) {
        item = rawItems[idx] as AnalyzedItem;
      }
    }
    root.render(<ArticlePage item={item} />);
  } else {
    // Handle array or object wrapper
    const rawItems = Array.isArray(portalData) ? portalData : (portalData as { items: AnalyzedItem[] }).items;

    // Update URLs to be dynamic and point to the article debug view
    const baseUrl = window.location.origin + window.location.pathname;
    const items = (rawItems as AnalyzedItem[]).map((item, index) => ({
      ...item,
      url: `${baseUrl}?type=article&i=${index}`,
    }));

    root.render(<PortalPage items={items} />);
  }
}