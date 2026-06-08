import { getBypassStorageKey } from './session-bypass';

describe('session bypass keying', () => {
  it('scopes bypass by tab id and origin', () => {
    expect(getBypassStorageKey(7, 'https://bank.example.com/login?step=1')).toBe(
      'ndf-bypass:7:https://bank.example.com'
    );
  });

  it('preserves explicit ports in the origin', () => {
    expect(getBypassStorageKey(12, 'https://example.com:8443/news')).toBe(
      'ndf-bypass:12:https://example.com:8443'
    );
  });
});
