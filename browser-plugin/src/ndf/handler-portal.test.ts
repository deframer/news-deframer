import { getDomain } from 'tldts';

import log from '../shared/logger';
import { NewsDeframerClient } from './client';
import { handlePortal } from './handler-portal';

jest.mock('tldts', () => ({
  getDomain: jest.fn(),
}));

describe('Portal Handler', () => {
  const mockClient = {
    getSite: jest.fn(),
  } as unknown as NewsDeframerClient;

  const mockReload = jest.fn();
  const mockStop = jest.fn();

  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
    jest.clearAllMocks();

    // Mock window.stop and window.location
    window.stop = mockStop;
    delete (window as Partial<Window & typeof globalThis>).location;
    Object.defineProperty(window, 'location', {
      value: {
        href: 'http://example.com/',
        hostname: 'example.com',
        reload: mockReload,
      },
      configurable: true,
    });
  });

  test('should fetch items for the root domain', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockResolvedValue([{ hash: '123', url: 'http://test.com', rating: 5 }]);
    await handlePortal(mockClient);

    expect(mockClient.getSite).toHaveBeenCalledWith('example.com');
    expect(mockStop).toHaveBeenCalled();
  });

  test('should log an error if root domain cannot be determined', async () => {
    (getDomain as jest.Mock).mockReturnValue(null);
    const logSpy = jest.spyOn(log, 'error').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Could not determine root domain. Reloading with bypass.');
    expect(mockReload).toHaveBeenCalled();
    logSpy.mockRestore();
  });

  test('should log info if no items are found', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockResolvedValue([]);
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('No items found for example.com. Reloading with bypass.');
    expect(mockReload).toHaveBeenCalled();
    logSpy.mockRestore();
  });

  test('should log an error if fetching fails', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockRejectedValue(new Error('API Error'));
    const logSpy = jest.spyOn(log, 'error').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Failed to fetch items for example.com:', expect.any(Error));
    expect(mockReload).toHaveBeenCalled();
    logSpy.mockRestore();
  });
});
