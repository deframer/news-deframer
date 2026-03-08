import { render } from '@testing-library/react';

import { TrendCompare } from './TrendCompare';

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

jest.mock('../../shared/settings', () => ({
  getSettings: jest.fn().mockResolvedValue({}),
}));

jest.mock('../client', () => ({
  NewsDeframerClient: jest.fn().mockImplementation(() => ({
    getDomainComparison: jest.fn().mockImplementation(() => new Promise(() => {})),
  })),
}));

jest.mock('./ArticleList', () => ({
  ArticleList: () => <div>ArticleList</div>,
}));

describe('TrendCompare', () => {
  beforeAll(() => {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: jest.fn().mockImplementation(() => ({
        matches: false,
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
      })),
    });
  });

  it('renders without crashing', () => {
    const { container } = render(
      <TrendCompare
        days={7}
        baseItems={[]}
        compareDomain={null}
        availableDomains={[]}
        onSelectDomain={jest.fn()}
        domain={{ domain: 'example.com', language: 'en' }}
        searchEngineUrl="https://search.example.com"
      />,
    );
    expect(container).not.toBeNull();
  });
});
