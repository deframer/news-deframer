import log from '../shared/logger';
import { NewsDeframerClient } from './client';
import { handleArticle } from './handler-article';

describe('Article Handler', () => {
  const mockClient = {
    getItem: jest.fn(),
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
        href: 'http://test.com/',
        hostname: 'test.com',
        reload: mockReload,
      },
      configurable: true,
    });

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

  test('should log info if no item is found', async () => {
    mockClient.getItem = jest.fn().mockResolvedValue(null);
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    await handleArticle(mockClient);
    expect(logSpy).toHaveBeenCalledWith('No item found for this URL. Reloading with bypass.');
    expect(mockReload).toHaveBeenCalled();
    logSpy.mockRestore();
  });
});
