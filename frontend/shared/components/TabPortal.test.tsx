import { render } from '@testing-library/react';

import { TabPortal } from './TabPortal';

jest.mock('./ArticleTile', () => ({
  ArticleTile: () => <div>ArticleTile</div>,
}));

describe('TabPortal', () => {
  it('renders without crashing', () => {
    const { container } = render(<TabPortal items={[]} />);
    expect(container).not.toBeNull();
  });
});
