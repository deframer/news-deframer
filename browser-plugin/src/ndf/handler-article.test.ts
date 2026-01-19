/**
 * @jest-environment-options { "url": "http://test.com/" }
 */

import log from '../shared/logger';
import { NewsDeframerClient } from './client';
import { handleArticle } from './handler-article';

describe('Article Handler', () => {
  const mockClient = {
    getItem: jest.fn(),
  } as unknown as NewsDeframerClient;

  const mockStop = jest.fn();

  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
    jest.clearAllMocks();
    jest.restoreAllMocks();

    // Mock window.stop and window.location
    window.stop = mockStop;
    // In JSDOM, window.location.reload is read-only and cannot be mocked.
    // Tests that trigger a reload will instead assert that the call throws the
    // expected "Not implemented" error from the JSDOM environment.

    sessionStorage.clear();
  });

  test('should fetch item for current URL', async () => {
    mockClient.getItem = jest.fn().mockResolvedValue({ hash: '123', url: 'http://test.com/', rating: 5 });
    await handleArticle(mockClient);

    expect(mockClient.getItem).toHaveBeenCalledWith('http://test.com/');
    expect(mockStop).toHaveBeenCalled();
  });

  test('should log a message to the console', async () => {
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    mockClient.getItem = jest.fn().mockResolvedValue({ hash: '123', url: 'http://test.com/', rating: 5 });
    await handleArticle(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Article page detected. Stopping window immediately.');
    logSpy.mockRestore();
  });

  test('should log info and set bypass if no item is found', async () => {
    mockClient.getItem = jest.fn().mockResolvedValue(null);
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});

    // This will call reload() and JSDOM will complain, but the test will continue.
    await handleArticle(mockClient);

    expect(logSpy).toHaveBeenCalledWith('No item found for this URL. Reloading with bypass.');
    expect(sessionStorage.getItem('__ndf-bypass')).toBe('true');
    logSpy.mockRestore();
  });
});
