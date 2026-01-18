import log from '../shared/logger';
import { NewsDeframerClient } from './client';
import { handleArticle } from './handler-article';

describe('Article Handler', () => {
  const mockClient = {
    getItem: jest.fn(),
  } as unknown as NewsDeframerClient;

  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
    jest.clearAllMocks();
  });

  test('should apply a blue border to the body and fetch item', async () => {
    mockClient.getItem = jest.fn().mockResolvedValue({ hash: '123', url: 'http://test.com' });
    await handleArticle(mockClient);
    expect(document.body.style.border).toBe('15px solid blue');
    expect(mockClient.getItem).toHaveBeenCalledWith(window.location.href);
  });

  test('should log a message to the console', async () => {
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    await handleArticle(mockClient);
    expect(logSpy).toHaveBeenCalledWith('Article page detected.');
    logSpy.mockRestore();
  });

  test('should log info if no item is found', async () => {
    mockClient.getItem = jest.fn().mockResolvedValue(null);
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    await handleArticle(mockClient);
    expect(logSpy).toHaveBeenCalledWith('No item found for this URL.');
    logSpy.mockRestore();
  });
});
