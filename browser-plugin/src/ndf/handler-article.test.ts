import log from '../shared/logger';
import { handleArticle } from './handler-article';

describe('Article Handler', () => {
  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
  });

  test('should apply a blue border to the body', () => {
    handleArticle();
    expect(document.body.style.border).toBe('15px solid blue');
  });

  test('should log a message to the console', () => {
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    handleArticle();
    expect(logSpy).toHaveBeenCalledWith('Article page detected.');
    logSpy.mockRestore();
  });
});
