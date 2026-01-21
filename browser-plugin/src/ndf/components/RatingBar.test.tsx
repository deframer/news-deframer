import { render } from '@testing-library/react';
import React from 'react';

import { RatingBar } from './RatingBar';

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

describe('RatingBar', () => {
  it('renders without crashing', () => {
    const { container } = render(<RatingBar value={0.75} label="Test Label" />);
    expect(container).not.toBeNull();
  });
});
