import { render } from '@testing-library/react';

import { TrendTagCloud } from './TrendTagCloud';

HTMLCanvasElement.prototype.getContext = jest.fn(() => null);

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
    getTopTrendByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
  })),
}));

jest.mock('./TrendDetails', () => ({
  TrendDetails: () => <div>TrendDetails</div>,
}));

describe('TrendTagCloud', () => {
  it('renders without crashing', () => {
    const { container } = render(
      <TrendTagCloud
        domain={{ domain: 'example.com', language: 'en' }}
        days={7}
        searchEngineUrl="https://search.example.com"
        activeTab="lifecycle"
        setActiveTab={jest.fn()}
      />,
    );
    expect(container).not.toBeNull();
  });
});
