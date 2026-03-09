import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
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

jest.mock('./ArticleList', () => ({
  ArticleList: () => <div>ArticleList</div>,
}));

describe('TrendLifecycleChart', () => {
  it('renders without crashing', () => {
    const api = {
      getLifecycleByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
    } as unknown as NewsDeframerApi;
    const { container } = render(<TrendLifecycleChart api={api} domain={{ domain: 'example.com', language: 'en' }} days={7} term="topic" />);
    expect(container).not.toBeNull();
  });
});
