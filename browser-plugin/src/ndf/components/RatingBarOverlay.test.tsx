import { render } from '@testing-library/react';
import React from 'react';

import { RatingBarOverlay } from './RatingBarOverlay';

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

describe('RatingBarOverlay', () => {
  it('renders without crashing', () => {
    const { container } = render(
      <RatingBarOverlay value={0.75} reason="Test Reason" />,
    );
    expect(container).not.toBeNull();
  });
});
