import { fireEvent, render, screen, waitFor } from '@testing-library/react';

import { ArticleList } from './ArticleList';

const mockGetSettings = jest.fn().mockResolvedValue({});
const mockGetArticlesByTrend = jest.fn();

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
  getSettings: () => mockGetSettings(),
}));

jest.mock('../client', () => ({
  NewsDeframerClient: jest.fn().mockImplementation(() => ({
    getArticlesByTrend: (...args: unknown[]) => mockGetArticlesByTrend(...args),
  })),
}));

describe('ArticleList', () => {
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
    mockGetSettings.mockClear();
    mockGetArticlesByTrend.mockReset();
    (HTMLElement.prototype.scrollIntoView as jest.Mock).mockClear();
  });

  it('renders without crashing', () => {
    mockGetArticlesByTrend.mockImplementation(() => new Promise(() => {}));
    const { container } = render(<ArticleList term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);
    expect(container).not.toBeNull();
  });

  it('scrolls into view after the first article load', async () => {
    mockGetArticlesByTrend.mockResolvedValue([
      { url: 'https://example.com/a', title: 'A', authors: ['Jane'], pub_date: '2026-03-08T00:00:00Z', rating: 0.4 },
    ]);

    render(<ArticleList term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);

    await waitFor(() => {
      expect(screen.getByText('A')).toBeInTheDocument();
    });

    expect(HTMLElement.prototype.scrollIntoView).toHaveBeenCalledTimes(1);
  });

  it('scrolls into view when moving to the next page', async () => {
    mockGetArticlesByTrend
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

    render(<ArticleList term="topic" domain={{ domain: 'example.com', language: 'en' }} days={7} />);

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
