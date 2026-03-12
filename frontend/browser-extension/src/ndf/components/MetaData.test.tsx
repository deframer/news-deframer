import { render, screen } from '@testing-library/react';

import { MetaData } from './MetaData';

const mockUseTranslation = jest.fn();

jest.mock('react-i18next', () => ({
  useTranslation: () => mockUseTranslation(),
  initReactI18next: {
    type: '3rdParty',
    init: jest.fn(),
  },
}));

describe('MetaData', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2026-03-08T12:00:00Z'));
    mockUseTranslation.mockReturnValue({
      t: (_key: string, fallback?: string) => fallback,
      i18n: {
        language: 'en',
      },
    });
  });

  afterEach(() => {
    jest.useRealTimers();
    mockUseTranslation.mockReset();
  });

  it('renders relative time for dates under one minute old', () => {
    render(<MetaData pubDate="2026-03-08T11:59:30Z" />);

    expect(screen.getByText('30 seconds ago')).toBeInTheDocument();
  });

  it('renders localized relative copy for German dates under one minute old', () => {
    mockUseTranslation.mockReturnValue({
      t: (_key: string, fallback?: string) => fallback,
      i18n: {
        language: 'de',
      },
    });

    render(<MetaData pubDate="2026-03-08T11:59:30Z" />);

    expect(screen.getByText('vor 30 Sekunden')).toBeInTheDocument();
  });

  it('renders future timestamps in larger units', () => {
    render(<MetaData pubDate="2026-03-10T23:00:20Z" />);

    expect(screen.getByText('in 2 days')).toBeInTheDocument();
  });

  it('renders the author value when present', () => {
    render(<MetaData author="Clark Kent" />);

    expect(screen.getByText('Clark Kent')).toBeInTheDocument();
  });

  it('renders a separator between time and author', () => {
    render(<MetaData pubDate="2026-03-08T11:59:30Z" author="Clark Kent" />);

    expect(screen.getByText('|')).toBeInTheDocument();
    expect(screen.getByText('Clark Kent')).toBeInTheDocument();
  });
});
