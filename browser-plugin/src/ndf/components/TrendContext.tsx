import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { TrendContextMetric, TrendRepo } from './TrendRepo';

interface TrendContextProps {
  topic: string;
  className?: string;
}

const contextCss = `
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
`;

export const TrendContext = ({ topic, className }: TrendContextProps) => {
  const { t } = useTranslation();
  const [items, setItems] = useState<TrendContextMetric[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    setLoading(true);
    TrendRepo.getTrendContext(topic).then((data) => {
      if (mounted) {
        setItems(data);
        setLoading(false);
      }
    });
    return () => { mounted = false; };
  }, [topic]);

  if (loading) return <div className={`trend-context ${className || ''}`}>Loading context...</div>;
  if (items.length === 0) return null;

  return (
    <div className={`trend-context ${className || ''}`}>
      <style>{contextCss}</style>
      <div className="context-header">
        Context: How is "{topic}" being described?
      </div>
      <div className="context-list">
        {items.map((item) => (
          <span key={item.context_word} className="context-chip">
            {item.context_word}
            <div className="chip-tooltip">
              {item.type} - Frequency: {item.frequency}
            </div>
          </span>
        ))}
      </div>
    </div>
  );
};