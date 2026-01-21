import { render } from '@testing-library/react';
import React from 'react';

import { AnalyzedItem } from '../client';
import { ArticleTile } from './ArticleTile';

describe('ArticleTile', () => {
  const mockItem: AnalyzedItem = {
    url: 'http://example.com/article1',
    title_original: 'Original Title',
    description_original: 'Original Description',
    rating: 0.85,
  };

  it('renders without crashing', () => {
    const { container } = render(<ArticleTile item={mockItem} />);
    expect(container).not.toBeNull();
  });
});
