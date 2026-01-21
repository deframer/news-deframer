import { render } from '@testing-library/react';
import React from 'react';

import { AnalyzedItem } from '../client';
import { PortalPage } from '../pages/PortalPage';

jest.mock('../../shared/logger', () => ({
  info: jest.fn(),
  error: jest.fn(),
}));

describe('PortalPage', () => {
  const mockItems: AnalyzedItem[] = [
    {
      url: 'http://example.com/article1',
      title_original: 'Original Title 1',
      description_original: 'Original Description 1',
      rating: 0.85,
    },
    {
      url: 'http://example.com/article2',
      title_original: 'Original Title 2',
      description_original: 'Original Description 2',
      rating: 0.45,
    },
  ];

  it('renders without crashing', () => {
    const { container } = render(<PortalPage items={mockItems} />);
    expect(container).not.toBeNull();
  });
});
