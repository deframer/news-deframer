export const ndfStyles = `
/* ArticlePage.tsx */
:host {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  display: block;
  background-color: var(--bg-color);
  color: var(--text-color);
  min-height: 100vh;
}

/* All styles from createArticleHtml are moved here */
.page-header { background-color: var(--header-bg); color: var(--text-color); padding: 8px 1.5em; text-align: left; border-bottom: 1px solid var(--border-color); z-index: 1001; display: flex; align-items: center; justify-content: space-between; }
.btn-back { display: inline-flex; align-items: center; gap: 6px; padding: 8px 12px; border: 1px solid var(--accent-color); border-radius: 8px; background-color: var(--accent-color); color: var(--accent-text); font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
.btn-back:hover { background-color: var(--accent-hover); }
.btn-hide { padding: 8px 12px; border: 1px solid var(--btn-border); border-radius: 8px; background-color: var(--btn-bg); color: var(--btn-text); font-size: 0.9em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; }
.btn-hide:hover { background-color: var(--btn-hover-bg); }
.article-container { max-width: 800px; margin: 0; background-color: var(--card-bg); border-radius: 0; box-shadow: none; padding-bottom: 20px; color: var(--text-color); }
.image-container img { width: 100%; height: auto; border-radius: 0; display: block; }
.main-content { padding: 1.5em; }
h1 { margin: 0 0 10px; font-size: 2em; color: var(--text-color); }
.description { font-size: 1.1em; color: var(--secondary-text); margin-bottom: 0.5em; }
.analysis-section { border-top: 2px solid var(--border-color); padding-top: 1.5em; margin-top: 1.5em; }
.metric-item { display: block; margin-bottom: 1.5em; }
.metric-label { font-weight: bold; font-size: 1.1em; margin-bottom: 5px; display: block; color: var(--text-color); }
.metric-details { display: block; width: 100%; }
.reason { margin: 5px 0 0; font-size: 0.95em; color: var(--secondary-text); line-height: 1.4; }
.original-content { margin-top: 2em; padding-top: 1em; border-top: 1px solid var(--border-color); }
.action-buttons { display: flex; justify-content: space-around; gap: 10px; padding: 1.5em; background-color: var(--card-bg); border-top: 1px solid var(--border-color); }
.btn { flex: 1; display: flex; align-items: center; justify-content: center; gap: 6px; padding: 10px; border: 1px solid var(--btn-border); border-radius: 8px; background-color: var(--btn-bg); color: var(--btn-text); font-size: 0.95em; font-weight: 500; cursor: pointer; transition: background-color 0.2s; text-decoration: none; }
.btn:hover { background-color: var(--btn-hover-bg); }
.btn-primary { background-color: var(--accent-color); color: var(--accent-text); border-color: var(--accent-color); }
.btn-primary:hover { background-color: var(--accent-hover); }

@media (min-width: 800px) {
  .page-header { border-radius: 12px 12px 0 0; }
  .article-container { margin: 2em auto; border-radius: 12px; box-shadow: var(--card-shadow); }
}

@media (max-width: 799px) {
  .article-container { padding-top: 55px; padding-bottom: 80px; }
  .page-header { position: fixed; top: 0; left: 0; width: 100%; box-sizing: border-box; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
  .action-buttons { position: fixed; bottom: 0; left: 0; width: 100%; z-index: 1000; box-shadow: 0 -2px 10px rgba(0,0,0,0.1); margin-top: 0; box-sizing: border-box; }
}

/* PortalPage.tsx */
:host {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  display: block;
  background-color: var(--bg-color);
  color: var(--text-color);
  min-height: 100vh;
}
.portal-container { padding: 2em; }
.header {
  background-color: var(--header-bg);
  color: var(--text-color);
  padding: 1em;
  text-align: center;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  margin-bottom: 2em;
}
h1 { margin: 0; font-size: 1.5em; }
.title-mobile { display: none; }
.btn {
  padding: 10px 16px;
  border: 1px solid var(--btn-border);
  border-radius: 8px;
  background-color: var(--btn-bg);
  color: var(--btn-text);
  font-size: 0.95em;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  text-decoration: none;
  flex-shrink: 0;
}
.btn:hover { background-color: var(--btn-hover-bg); }

#btn-hide {
  position: absolute;
  right: 1.5em;
  top: 50%;
  transform: translateY(-50%);
}

.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }

.tabs { display: flex; justify-content: center; margin-bottom: 20px; border-bottom: 1px solid var(--border-color, #ccc); }
.tab-btn { padding: 12px 24px; background: none; border: none; border-bottom: 3px solid transparent; cursor: pointer; font-size: 1.1em; color: var(--text-color); opacity: 0.6; transition: all 0.2s; }
.tab-btn:hover { opacity: 1; background-color: rgba(0,0,0,0.05); }
.tab-btn.active { border-bottom-color: var(--btn-bg); opacity: 1; font-weight: 600; }

@media (max-width: 799px) {
  .portal-container { padding: 0 1em 1em; }
  .header { position: sticky; top: 0; z-index: 1000; margin-bottom: 0; justify-content: space-between; padding: 10px 1em; height: 60px; margin-left: -1em; margin-right: -1em; width: calc(100% + 2em); }
  .tabs { position: sticky; top: 60px; z-index: 999; background-color: var(--bg-color); margin-bottom: 0; margin-left: -1em; margin-right: -1em; width: calc(100% + 2em); }
  .title-desktop { display: none; }
  .title-mobile { display: inline; }
  #btn-hide { right: 1em; }
  .footer-container.trends { display: none; }
  .grid { margin-top: 20px; }
}

/* Footer.tsx */
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

/* RatingBarOverlay.tsx */
.bar-container {
  background-color: var(--rating-bg);
  height: 30px;
  width: 100%;
  position: relative; /* For overlay and tooltip */
}
.bar {
  height: 100%;
  transition: width 0.3s ease;
}
.bar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  padding-left: 10px;
  font-weight: bold;
  font-size: 0.9em;
  pointer-events: none;
}
.bar-container .tooltip-text {
  visibility: hidden;
  opacity: 0;
  width: 250px;
  background-color: var(--tooltip-bg);
  color: var(--tooltip-text);
  text-align: left;
  border-radius: 6px;
  padding: 10px;
  position: absolute;
  z-index: 1;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  transition: opacity 0.2s;
  font-size: 0.9em;
  line-height: 1.4;
  pointer-events: none;
}
.bar-container:hover .tooltip-text {
  visibility: visible;
  opacity: 1;
}

/* ArticleTile.tsx */
.tile-link { text-decoration: none; color: inherit; display: block; height: 100%; }
.tile { background-color: var(--card-bg); border-radius: 8px; box-shadow: var(--card-shadow); overflow: hidden; transition: transform 0.2s, box-shadow 0.2s; height: 100%; position: relative; display: flex; flex-direction: column; }
.tile:hover { transform: translateY(-5px); box-shadow: var(--card-shadow-hover); }
.image-container { position: relative; }
.image-container img { width: 100%; height: auto; display: block; }
.content { padding: 15px; flex: 1; }
h3 { margin: 0 0 10px; font-size: 1.1em; color: var(--text-color); }
p { font-size: 0.9em; color: var(--secondary-text); margin: 0; }
.tile-footer {
  height: 30px;
  padding: 0 15px;
  display: flex;
  align-items: center;
  border-top: 1px solid var(--border-color, #eee);
}
.tile-footer .meta-data {
  margin: 0;
  font-size: 0.9em;
}

/* Spinner.tsx */
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* TrendContext.tsx */
.trend-context {
  padding: 12px 16px;
  background-color: var(--bg-color-secondary, #f9f9f9);
  border-radius: 8px;
  margin-top: 8px;
  border-left: 4px solid var(--primary-color, #0056b3);
  animation: fadeIn 0.3s ease-in-out;
}
.context-header {
  font-size: 0.85em;
  text-transform: uppercase;
  color: var(--secondary-text);
  margin-bottom: 8px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}
.context-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.context-chip {
  background: var(--card-bg, #fff);
  border: 1px solid var(--border-color, #e0e0e0);
  padding: 4px 10px;
  border-radius: 16px;
  font-size: 0.9em;
  color: var(--text-color);
  display: flex;
  align-items: center;
  gap: 6px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}
.context-chip {
  position: relative;
}
.context-chip:hover {
  z-index: 10;
}
.chip-tooltip {
  visibility: hidden;
  opacity: 0;
  background-color: var(--tooltip-bg, rgba(0,0,0,0.8));
  color: var(--tooltip-text, #fff);
  text-align: center;
  border-radius: 6px;
  padding: 6px 10px;
  position: absolute;
  z-index: 20;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 6px;
  font-size: 0.9em;
  white-space: nowrap;
  pointer-events: none;
  transition: opacity 0.2s;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}
.context-chip:hover .chip-tooltip {
  visibility: visible;
  opacity: 1;
}
.context-freq {
  background: var(--badge-bg, #eee);
  color: var(--secondary-text);
  font-size: 0.8em;
  padding: 1px 5px;
  border-radius: 8px;
  min-width: 16px;
  text-align: center;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}

/* RatingBar.tsx */
.bar-container {
  background-color: var(--rating-bg);
  border-radius: 5px;
  height: 30px;
  width: 100%;
  position: relative;
  overflow: hidden;
}
.bar {
  height: 100%;
  border-radius: 5px;
  transition: width 0.3s ease;
}
.bar-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  padding-left: 8px;
  font-weight: bold;
  font-size: 0.95em;
  pointer-events: none;
}

/* TrendCompare.tsx */
.compare-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
  margin-top: 10px;
}
@media (max-width: 799px) {
  .compare-grid {
    grid-template-columns: 1fr;
  }
}
.compare-col {
  background: var(--card-bg, #fff);
  border: 1px solid var(--border-color, #eee);
  border-radius: 8px;
  padding: 12px;
  display: flex;
  flex-direction: column;
}
.col-header {
  font-size: 0.9em;
  font-weight: bold;
  text-align: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid transparent;
  color: var(--text-color);
}
.col-header select {
  font-size: 1em;
  font-weight: bold;
}
.col-header.unique-a { border-color: var(--primary-color, #0056b3); }
.col-header.intersect { border-color: var(--secondary-text, #6c757d); }
.col-header.unique-b { border-color: var(--danger-color, #dc3545); }

.compare-list {
  list-style: none;
  padding: 0;
  margin: 0;
}
.compare-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--border-color, #f5f5f5);
  font-size: 0.9em;
}
.compare-item:last-child { border-bottom: none; }
.topic-name { font-weight: 500; }
.topic-score {
  font-size: 0.8em;
  color: var(--secondary-text);
  background: var(--bg-color-secondary, #f5f5f5);
  padding: 2px 6px;
  border-radius: 4px;
}
.header-select {
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid var(--border-color, #ddd);
  background: var(--card-bg, #fff);
  color: var(--text-color);
  font-size: 0.9em;
  font-weight: bold;
  max-width: 100%;
}

/* TrendLifecycle.tsx */
.lifecycle-controls {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}
.lifecycle-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border-color, #ddd);
  border-radius: 6px;
  font-size: 0.95em;
  background: var(--card-bg, #fff);
  color: var(--text-color);
}
.lifecycle-btn {
  padding: 8px 16px;
  background: var(--primary-color, #0056b3);
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.2s;
}
.lifecycle-btn:hover {
  background: var(--primary-color-dark, #004494);
}
.chart-container {
  height: 300px;
  display: flex;
  align-items: flex-end;
  gap: 4px;
  padding: 20px 0 70px 0;
  border-bottom: 1px solid var(--border-color, #ccc);
  position: relative;
  width: 100%;
}
@media (max-width: 799px) {
  .chart-container {
    gap: 1px;
  }
}
.chart-bar-wrapper {
  flex: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  position: relative;
  outline: none;
}
.chart-bar-wrapper:hover,
.chart-bar-wrapper:focus-visible {
  z-index: 20;
}
.chart-bar-wrapper:focus-visible .chart-bar {
  box-shadow: 0 0 0 2px var(--text-color);
  opacity: 0.8;
}
.chart-bar {
  width: 100%;
  border-radius: 4px 4px 0 0;
  transition: height 0.5s ease;
  min-height: 65px;
  position: relative;
}
.chart-bar:hover {
  opacity: 0.8;
}
.chart-bar.lateral {
  background-color: var(--badge-bg, #e9ecef);
  border-top: 3px dotted var(--secondary-text);
}
.trend-icon {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  font-size: 0.8em;
  line-height: 1;
  margin-bottom: 4px;
  font-weight: bold;
}
.bar-label {
  font-size: 0.85em;
  font-weight: bold;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  position: absolute;
  bottom: 8px;
  left: 0;
  width: 100%;
  pointer-events: none;
  z-index: 2;
}
@media (max-width: 799px) {
  .bar-label {
    left: 50%;
    bottom: 25px;
    transform: translateX(-50%) rotate(-90deg);
    transform-origin: center;
    text-align: center;
    width: 80px;
  }
}
.bar-tooltip {
  position: absolute;
  bottom: 35px;
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--tooltip-bg, rgba(0,0,0,0.8));
  color: var(--tooltip-text, #fff);
  padding: 10px;
  border-radius: 6px;
  font-size: 0.9em;
  line-height: 1.4;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s;
  z-index: 10;
}
@media (max-width: 799px) {
  .bar-tooltip {
    bottom: 70px;
  }
}
.chart-bar-wrapper:hover .bar-tooltip,
.chart-bar-wrapper:focus-visible .bar-tooltip {
  opacity: 1;
}

/* MetaData.tsx */
.meta-data {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  font-size: 0.85em;
  color: var(--secondary-text);
  margin-top: 0.5em;
  margin-bottom: 2em;
  line-height: 1;
}
.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.meta-icon {
  width: 14px;
  height: 14px;
  stroke: currentColor;
  stroke-width: 2;
  fill: none;
}

/* TrendTopTagCloud.tsx */
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-content: center;
  gap: 12px;
  padding: 20px 10px;
  min-height: 200px;
}
.tag-item {
  transition: transform 0.2s, color 0.2s;
  cursor: pointer;
  line-height: 1;
  color: var(--text-color);
  position: relative;
}
.tag-item:hover {
  transform: scale(1.1);
  color: var(--primary-color, #0056b3);
  z-index: 10;
}
.tag-item.active {
  color: var(--primary-color, #0056b3);
  font-weight: bold;
  background-color: var(--bg-color-secondary, #f0f0f0);
  padding: 4px 10px;
  border-radius: 12px;
  transform: scale(1.1);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}
.cloud-context-container {
  margin-top: 10px;
  animation: fadeIn 0.3s ease;
}
.cloud-tooltip {
  visibility: hidden;
  opacity: 0;
  background-color: var(--tooltip-bg, rgba(0,0,0,0.8));
  color: var(--tooltip-text, #fff);
  text-align: center;
  border-radius: 6px;
  padding: 8px;
  position: absolute;
  z-index: 20;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 8px;
  font-size: 12px;
  line-height: 1.4;
  pointer-events: none;
  white-space: nowrap;
  transition: opacity 0.2s;
  font-weight: normal;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}
.tag-item:hover .cloud-tooltip {
  visibility: visible;
  opacity: 1;
}

/* TabTrend.tsx */
.trend-container {
  display: flex;
  flex-direction: column;
  background: var(--card-bg, #fff);
  min-height: 400px;
}
.trend-footer { display: none; }
@media (max-width: 799px) {
  .trend-container {
    height: auto;
    min-height: 0;
    display: block;
  }
  .trend-content {
    overflow-y: visible;
    flex: none;
  }
  .trend-footer { display: block; margin-top: 20px; }
}

.trend-header {
  background: var(--card-bg, #fff);
  flex-shrink: 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  z-index: 5;
  position: relative;
}
@media (max-width: 799px) {
  .trend-header {
    position: sticky;
    top: 105px;
    z-index: 900;
  }
}

/* 1. Top Navigation Tabs */
.nav-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color, #eee);
  background: var(--card-bg, #fff);
  overflow-x: auto;
  scrollbar-width: none;
}
.nav-tabs::-webkit-scrollbar { display: none; }
.nav-tab {
  flex: 1;
  padding: 14px 10px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 0.95em;
  font-weight: 500;
  color: var(--secondary-text);
  border-bottom: 3px solid transparent;
  transition: all 0.2s;
  text-align: center;
  white-space: nowrap;
  min-width: fit-content;
}
.nav-tab:hover {
  color: var(--primary-color, #0056b3);
  background-color: var(--hover-bg, rgba(0,0,0,0.02));
}
.nav-tab.active {
  color: var(--primary-color, #0056b3);
  border-bottom-color: var(--primary-color, #0056b3);
  font-weight: 600;
}
@media (max-width: 799px) {
  .nav-tab {
    padding: 12px 5px;
    font-size: 0.9em;
  }
}

/* 2. Filter Bar (Time Selection) */
.filter-bar {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color, #f0f0f0);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  background: var(--bg-color-secondary, #fafafa);
  gap: 10px;
}
.filter-label {
  font-size: 0.85em;
  font-weight: 600;
  color: var(--secondary-text);
}
.time-selector {
  display: flex;
  gap: 2px;
  background: var(--border-color, #e0e0e0);
  padding: 2px;
  border-radius: 6px;
}
.time-btn {
  border: none;
  background: transparent;
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 0.8em;
  cursor: pointer;
  color: var(--secondary-text);
  font-weight: 500;
  transition: all 0.1s;
}
.time-btn:hover {
  color: var(--text-color);
}
.time-btn.active {
  background: var(--card-bg, #fff);
  color: var(--text-color);
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
  font-weight: 600;
}

/* 3. Content Area */
.trend-content {
  padding: 16px;
  flex: 1;
  overflow-y: auto;
}
.domain-select {
  padding: 6px 10px;
  border-radius: 6px;
  border: 1px solid var(--border-color, #ddd);
  background: var(--card-bg, #fff);
  color: var(--text-color);
  font-size: 0.9em;
}
`;