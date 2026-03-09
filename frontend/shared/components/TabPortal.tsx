import { AnalyzedItem } from '@frontend-shared/ndf-api-interfaces';
import { ArticleTile } from './ArticleTile';

interface TabPortalProps {
  items: AnalyzedItem[];
}

export const TabPortal = ({ items }: TabPortalProps) => {
  return (
    <div className="grid">
      {items.map((item) => (
        <ArticleTile item={item} key={item.url} />
      ))}
    </div>
  );
};
