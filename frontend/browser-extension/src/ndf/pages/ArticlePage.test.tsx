import { fireEvent, render, screen } from '@testing-library/react';

import { PROMPTS_URL, REFERENCE_TAB_TARGET } from '../../shared/links';
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
    llm_model: 'mistralai/ministral-3-3b',
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

  it('shows the analysis source before the action buttons', () => {
    const { container } = render(<ArticlePage item={mockItem} />);

    expect(screen.getByText('article.analysis_source_llm_model_label')).toBeInTheDocument();
    expect(screen.getByText('mistralai/ministral-3-3b')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'article.analysis_source_prompts_link' })).toHaveAttribute('href', PROMPTS_URL);
    expect(screen.getByRole('link', { name: 'article.analysis_source_prompts_link' })).toHaveAttribute('target', REFERENCE_TAB_TARGET);

    const analysisSource = container.querySelector('.analysis-source');
    const actionButtons = container.querySelector('.action-buttons');

    expect(analysisSource).not.toBeNull();
    expect(actionButtons).not.toBeNull();
    expect(analysisSource!.compareDocumentPosition(actionButtons!)).toBe(Node.DOCUMENT_POSITION_FOLLOWING);
  });

  it('hides the model line when no model is present', () => {
    const itemWithoutModel = { ...mockItem };
    delete itemWithoutModel.llm_model;

    render(<ArticlePage item={itemWithoutModel} />);

    expect(screen.queryByText('article.analysis_source_llm_model_label')).not.toBeInTheDocument();
    expect(screen.queryByText('mistralai/ministral-3-3b')).not.toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'article.analysis_source_prompts_link' })).toHaveAttribute('href', PROMPTS_URL);
    expect(screen.getByRole('link', { name: 'article.analysis_source_prompts_link' })).toHaveAttribute('target', REFERENCE_TAB_TARGET);
  });
});
