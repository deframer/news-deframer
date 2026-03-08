import { render } from '@testing-library/react';

import { ArticleList } from './ArticleList';

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
    getArticlesByTrend: jest.fn().mockImplementation(() => new Promise(() => {})),
  })),
}));

describe('ArticleList', () => {
  it('renders without crashing', () => {
    const { container } = render(<ArticleList term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);
    expect(container).not.toBeNull();
  });
});
