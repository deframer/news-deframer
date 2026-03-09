import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import { NewsDeframerApi } from '@frontend-shared/ndf-api-interfaces';
import { ArticleList } from './ArticleList';

const getArticlesByTrendMock = jest.fn();

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

describe('ArticleList', () => {
  const api: NewsDeframerApi = {
    getDomains: jest.fn(),
    getItem: jest.fn(),
    getSite: jest.fn(),
    getTopTrendByDomain: jest.fn(),
    getContextByDomain: jest.fn(),
    getLifecycleByDomain: jest.fn(),
    getDomainComparison: jest.fn(),
    getArticlesByTrend: (...args) => getArticlesByTrendMock(...args),
  };

  beforeAll(() => {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: jest.fn().mockImplementation(() => ({
        matches: false,
        addEventListener: jest.fn(),
        removeEventListener: jest.fn(),
      })),
    });

    Object.defineProperty(HTMLElement.prototype, 'scrollIntoView', {
      writable: true,
      value: jest.fn(),
    });
  });

  beforeEach(() => {
    getArticlesByTrendMock.mockReset();
    (HTMLElement.prototype.scrollIntoView as jest.Mock).mockClear();
  });

  it('renders without crashing', () => {
    getArticlesByTrendMock.mockImplementation(() => new Promise(() => {}));
    const { container } = render(<ArticleList api={api} term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);
    expect(container).not.toBeNull();
  });

  it('scrolls into view after the first article load', async () => {
    getArticlesByTrendMock.mockResolvedValue([
      { url: 'https://example.com/a', title: 'A', authors: ['Jane'], pub_date: '2026-03-08T00:00:00Z', rating: 0.4 },
    ]);

    render(<ArticleList api={api} term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);

    await waitFor(() => {
      expect(screen.getByText('A')).toBeInTheDocument();
    });

    expect(HTMLElement.prototype.scrollIntoView).toHaveBeenCalledTimes(1);
  });

  it('scrolls into view when moving to the next page', async () => {
    getArticlesByTrendMock
      .mockResolvedValueOnce(Array.from({ length: 10 }, (_, index) => ({
        url: `https://example.com/${index}`,
        title: `Article ${index}`,
        authors: ['Jane'],
        pub_date: '2026-03-08T00:00:00Z',
        rating: 0.4,
      })))
      .mockResolvedValueOnce([
        { url: 'https://example.com/next', title: 'Next Page', authors: ['John'], pub_date: '2026-03-09T00:00:00Z', rating: 0.5 },
      ]);

    render(<ArticleList api={api} term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);

    await waitFor(() => {
      expect(screen.getByText('Article 0')).toBeInTheDocument();
    });

    (HTMLElement.prototype.scrollIntoView as jest.Mock).mockClear();

    fireEvent.click(screen.getByRole('button', { name: /next/i }));

    await waitFor(() => {
      expect(screen.getByText('Next Page')).toBeInTheDocument();
    });

    expect(HTMLElement.prototype.scrollIntoView).toHaveBeenCalledTimes(1);
  });
});
