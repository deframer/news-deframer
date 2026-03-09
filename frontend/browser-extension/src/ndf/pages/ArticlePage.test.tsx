import { AnalyzedItem } from '@frontend-shared/ndf-api-interfaces';
import { render } from '@testing-library/react';

import { ArticlePage } from '../pages/ArticlePage';

jest.mock('@frontend-shared/logger', () => ({
  info: jest.fn(),
  error: jest.fn(),
}));

jest.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: {
      language: 'en',
      changeLanguage: jest.fn(),
    },
  }),
  initReactI18next: {
    type: '3rdParty',
    init: jest.fn(),
  },
}));

describe('ArticlePage', () => {
  const mockItem: AnalyzedItem = {
    url: 'http://example.com/article1',
    hash: 'test-hash',
    title_original: 'Original Title',
    description_original: 'Original Description',
    rating: 0.85,
  };

  it('renders without crashing', () => {
    const { container } = render(<ArticlePage item={mockItem} />);
    expect(container).not.toBeNull();
  });
});
