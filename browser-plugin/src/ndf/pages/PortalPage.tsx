import React from 'react';
import { getDomain } from 'tldts';

import log from '../../shared/logger';
import { AnalyzedItem } from '../client';
import { ArticleTile } from '../components/ArticleTile';
import { Footer } from '../components/Footer';

const portalPageCss = `
  body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    margin: 0;
    background-color: #f0f2f5;
  }
  .container { padding: 2em; }
  .header {
    background-color: #fff;
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
  .btn {
    padding: 10px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #fff;
    color: #333;
    font-size: 0.95em;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
    text-decoration: none;
    flex-shrink: 0;
  }
  .btn:hover { background-color: #f8f9fa; }

  #btn-hide {
    position: absolute;
    right: 1.5em;
    top: 50%;
    transform: translateY(-50%);
  }

  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 20px; }
`;

interface PortalPageProps {
  items: AnalyzedItem[];
}

export const PortalPage = ({ items }: PortalPageProps) => {
  const rootDomain = getDomain(window.location.hostname) || window.location.hostname;

  const bypassAndReload = () => {
    window.scrollTo(0, 0);
    log.info('Bypassing for this session and reloading.');
    sessionStorage.setItem('__ndf-bypass', 'true');
    window.location.reload();
  };

  return (
    <>
      <style>{portalPageCss}</style>
      <div className="container">
        <div className="header">
          <h1>News Deframer: {rootDomain}</h1>
          <button id="btn-hide" className="btn" onClick={bypassAndReload}>Hide</button>
        </div>
        <div className="grid">
          {items.map((item) => (
            <ArticleTile item={item} key={item.url} />
          ))}
        </div>
        <Footer />
      </div>
    </>
  );
};
