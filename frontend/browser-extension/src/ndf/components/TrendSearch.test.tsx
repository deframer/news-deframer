import { render } from '@testing-library/react';

import { TrendSearch } from './TrendSearch';

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

jest.mock('./TrendDetails', () => ({
  TrendDetails: () => <div>TrendDetails</div>,
}));

describe('TrendSearch', () => {
  it('renders without crashing', () => {
    const { container } = render(
      <TrendSearch
        domain={{ domain: 'example.com', language: 'en' }}
        days={7}
        activeTab="lifecycle"
        setActiveTab={jest.fn()}
      />,
    );
    expect(container).not.toBeNull();
  });
});
