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
      t: (key: string, fallback?: string) => {
        if (key === 'metadata.just_now') return 'a moment ago';
        return fallback ?? key;
      },
      i18n: {
        language: 'en',
      },
    });
  });

  afterEach(() => {
    jest.useRealTimers();
    mockUseTranslation.mockReset();
  });

  it('renders a moment ago for dates under one minute old', () => {
    render(<MetaData pubDate="2026-03-08T11:59:30Z" />);

    expect(screen.getByText('a moment ago')).toBeInTheDocument();
  });

  it('renders localized moment-ago copy for German dates under one minute old', () => {
    mockUseTranslation.mockReturnValue({
      t: (key: string, fallback?: string) => {
        if (key === 'metadata.just_now') return 'vor einem Augenblick';
        return fallback ?? key;
      },
      i18n: {
        language: 'de',
      },
    });

    render(<MetaData pubDate="2026-03-08T11:59:30Z" />);

    expect(screen.getByText('vor einem Augenblick')).toBeInTheDocument();
  });

  it('renders future timestamps in larger units instead of raw seconds', () => {
    render(<MetaData pubDate="2026-03-10T23:00:20Z" />);

    expect(screen.getByText('in 2d')).toBeInTheDocument();
    expect(screen.queryByText(/in .*s$/)).not.toBeInTheDocument();
  });
});
