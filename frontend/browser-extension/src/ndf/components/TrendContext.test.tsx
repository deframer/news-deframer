import { render } from '@testing-library/react';

import { TrendContextChart } from './TrendContext';

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
    getContextByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
  })),
}));

describe('TrendContextChart', () => {
  it('renders without crashing', () => {
    const { container } = render(<TrendContextChart topic="topic" days={7} domain={{ domain: 'example.com', language: 'en' }} />);
    expect(container).not.toBeNull();
  });
});
