import { AnalyzedItem, NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { render } from '@testing-library/react';

import { PortalPage } from '../pages/PortalPage';

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

describe('PortalPage', () => {
  const mockItems: AnalyzedItem[] = [
    {
      url: 'http://example.com/article1',
      hash: 'test-hash-1',
      title_original: 'Original Title 1',
      description_original: 'Original Description 1',
      rating: 0.85,
    },
    {
      url: 'http://example.com/article2',
      hash: 'test-hash-2',
      title_original: 'Original Title 2',
      description_original: 'Original Description 2',
      rating: 0.45,
    },
  ];

  it('renders without crashing', () => {
    const api = {} as NewsDeframerApi;
    const { container } = render(
      <PortalPage 
        api={api}
        items={mockItems} 
        domain={{ domain: 'example.com', portal_url: 'https://example.com', language: 'en' }}
        availableDomains={[]}
        searchEngineUrl="https://search.brave.com"
      />
    );
    expect(container).not.toBeNull();
  });
});
