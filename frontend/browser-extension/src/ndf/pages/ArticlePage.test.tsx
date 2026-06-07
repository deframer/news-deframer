import { fireEvent, render, screen } from '@testing-library/react';

import { AnalyzedItem } from '../client';
import { ArticlePage } from '../pages/ArticlePage';

jest.mock('../../shared/logger', () => ({
  info: jest.fn(),
  error: jest.fn(),
}));

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

describe('ArticlePage', () => {
  const mockItem: AnalyzedItem = {
    url: 'http://example.com/article1',
    hash: 'test-hash',
    title_original: 'Original Title',
    description_original: 'Original Description',
    rating: 0.85,
  };

  it('renders without crashing', () => {
    const { container } = render(<ArticlePage item={mockItem} />);
    expect(container).not.toBeNull();
  });

  it('shows public service media badge when tagged', () => {
    render(<ArticlePage item={{ ...mockItem, tags: ['news', 'public_service_media'] }} />);

    expect(screen.getByText('article.public_service_media')).toBeInTheDocument();
  });

  it('switches back to the article tab when opening original content from sentiments', () => {
    render(<ArticlePage item={mockItem} />);

    fireEvent.click(screen.getByRole('button', { name: 'article.tab_sentiments' }));
    fireEvent.click(screen.getByRole('button', { name: 'article.btn_original_title' }));

    expect(screen.getByRole('button', { name: 'article.tab_article' })).toHaveClass('active');
    expect(screen.getByText('article.original_section')).toBeInTheDocument();
  });

  it('switches back to the article tab when opening details from sentiments', () => {
    render(<ArticlePage item={mockItem} />);

    fireEvent.click(screen.getByRole('button', { name: 'article.tab_sentiments' }));
    fireEvent.click(screen.getByRole('button', { name: 'article.btn_details' }));

    expect(screen.getByRole('button', { name: 'article.tab_article' })).toHaveClass('active');
    expect(screen.getByText('metrics.clickbait')).toBeInTheDocument();
  });
});
