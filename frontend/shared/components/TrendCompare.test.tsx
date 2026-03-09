import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
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
    const api = {
      getDomainComparison: jest.fn().mockImplementation(() => new Promise(() => {})),
    } as unknown as NewsDeframerApi;
    const { container } = render(
      <TrendCompare
        api={api}
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
