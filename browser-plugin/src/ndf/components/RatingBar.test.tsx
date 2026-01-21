import { render } from '@testing-library/react';
import React from 'react';

import { RatingBar } from './RatingBar';

describe('RatingBar', () => {
  it('renders without crashing', () => {
    const { container } = render(<RatingBar value={0.75} label="Test Label" />);
    expect(container).not.toBeNull();
  });
});
