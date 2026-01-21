import { render } from '@testing-library/react';
import React from 'react';

import { RatingBarOverlay } from './RatingBarOverlay';

describe('RatingBarOverlay', () => {
  it('renders without crashing', () => {
    const { container } = render(
      <RatingBarOverlay value={0.75} reason="Test Reason" />,
    );
    expect(container).not.toBeNull();
  });
});
