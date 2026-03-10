import { render } from '@testing-library/react';

import { TrendLifecycleChart } from './TrendLifecycleChart';

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
    getLifecycleByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
  })),
}));

jest.mock('./ArticleList', () => ({
  ArticleList: () => <div>ArticleList</div>,
}));

describe('TrendLifecycleChart', () => {
  it('renders without crashing', () => {
    const { container } = render(<TrendLifecycleChart domain={{ domain: 'example.com', language: 'en' }} days={7} term="topic" />);
    expect(container).not.toBeNull();
  });
});
