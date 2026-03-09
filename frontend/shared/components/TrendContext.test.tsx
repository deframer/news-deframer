import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
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

describe('TrendContextChart', () => {
  it('renders without crashing', () => {
    const api = {
      getContextByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
    } as unknown as NewsDeframerApi;
    const { container } = render(<TrendContextChart api={api} topic="topic" days={7} domain={{ domain: 'example.com', language: 'en' }} />);
    expect(container).not.toBeNull();
  });
});
