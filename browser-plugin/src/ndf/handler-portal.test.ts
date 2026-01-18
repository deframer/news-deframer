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

  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
    jest.clearAllMocks();
  });

  test('should apply a green border to the body and fetch items', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockResolvedValue([]);
    await handlePortal(mockClient);
    expect(document.body.style.border).toBe('15px solid green');
    expect(mockClient.getSite).toHaveBeenCalledWith('example.com');
  });

  test('should log an error if root domain cannot be determined', async () => {
    (getDomain as jest.Mock).mockReturnValue(null);
    const logSpy = jest.spyOn(log, 'error').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Could not determine root domain.');
    logSpy.mockRestore();
  });

  test('should log info if no items are found', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockResolvedValue([]);
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('No items found for example.com.');
    logSpy.mockRestore();
  });

  test('should log an error if fetching fails', async () => {
    (getDomain as jest.Mock).mockReturnValue('example.com');
    mockClient.getSite = jest.fn().mockRejectedValue(new Error('API Error'));
    const logSpy = jest.spyOn(log, 'error').mockImplementation(() => {});
    await handlePortal(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Failed to fetch items for example.com:', expect.any(Error));
    logSpy.mockRestore();
  });
});
