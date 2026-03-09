import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { TrendDetails } from './TrendDetails';

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

jest.mock('./TrendContext', () => ({
  TrendContextChart: () => <div>TrendContextChart</div>,
}));

jest.mock('./TrendLifecycleChart', () => ({
  TrendLifecycleChart: () => <div>TrendLifecycleChart</div>,
}));

describe('TrendDetails', () => {
  it('renders without crashing', () => {
    const api = {} as NewsDeframerApi;
    const { container } = render(
      <TrendDetails
        api={api}
        term="topic"
        domain={{ domain: 'example.com', language: 'en' }}
        days={7}
        activeTab="lifecycle"
        setActiveTab={jest.fn()}
      />,
    );
    expect(container).not.toBeNull();
  });
});
