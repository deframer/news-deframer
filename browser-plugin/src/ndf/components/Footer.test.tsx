import { render, screen } from '@testing-library/react';
import React from 'react';

import { Footer } from './Footer';

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

describe('Footer', () => {
  it('renders the link with the correct href and text', () => {
    render(<Footer />);
    const linkElement = screen.getByText('News Deframer');
    expect(linkElement).toBeInTheDocument();
    expect(linkElement).toHaveAttribute(
      'href',
      'https://deframer.github.io/',
    );
  });
});
