import { render } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { TrendTagCloud } from './TrendTagCloud';

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

jest.mock('@visx/responsive', () => ({
  ParentSize: ({ children }: { children: ({ width, height }: { width: number; height: number }) => React.ReactNode }) => children({ width: 400, height: 300 }),
}));

jest.mock('@visx/scale', () => ({
  scaleLog: () => (value: number) => value,
}));

jest.mock('@visx/text', () => ({
  Text: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}));

jest.mock('@visx/wordcloud', () => ({
  Wordcloud: ({ children }: { children: (words: never[]) => React.ReactNode }) => <div>{children([])}</div>,
}));

jest.mock('./TrendDetails', () => ({
  TrendDetails: () => <div>TrendDetails</div>,
}));

describe('TrendTagCloud', () => {
  it('renders without crashing', () => {
    const api = {
      getTopTrendByDomain: jest.fn().mockImplementation(() => new Promise(() => {})),
    } as unknown as NewsDeframerApi;
    const { container } = render(
      <TrendTagCloud
        api={api}
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
