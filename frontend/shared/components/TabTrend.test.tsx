import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { TabTrend } from './TabTrend';

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

jest.mock('@frontend-shared/logger', () => ({
  info: jest.fn(),
  error: jest.fn(),
}));

jest.mock('./Footer', () => ({
  Footer: () => <div>Footer</div>,
}));

jest.mock('./TrendCompare', () => ({
  TrendCompare: () => <div>TrendCompare</div>,
}));

jest.mock('./TrendSearch', () => ({
  TrendSearch: () => <div>TrendSearch</div>,
}));

jest.mock('./TrendTagCloud', () => ({
  TrendTagCloud: () => <div>TrendTagCloud</div>,
}));

describe('TabTrend', () => {
  it('renders without crashing', () => {
    const api = {} as NewsDeframerApi;
    const { container } = render(
      <TabTrend
        api={api}
        domain={{ domain: 'example.com', language: 'en' }}
        availableDomains={[{ domain: 'example.com', language: 'en' }]}
        searchEngineUrl="https://search.example.com"
      />,
    );
    expect(container).not.toBeNull();
  });
});
