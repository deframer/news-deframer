import { render, screen } from '@testing-library/react';
import React from 'react';

import { Footer } from './Footer';

describe('Footer', () => {
  it('renders the link with the correct href and text', () => {
    render(<Footer />);
    const linkElement = screen.getByText('News Deframer');
    expect(linkElement).toBeInTheDocument();
    expect(linkElement).toHaveAttribute(
      'href',
      'https://github.com/egandro/news-deframer',
    );
  });
});
