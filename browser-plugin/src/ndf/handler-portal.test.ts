import log from '../shared/logger';
import { handlePortal } from './handler-portal';

describe('Portal Handler', () => {
  beforeEach(() => {
    // Reset the document body for each test
    document.body.innerHTML = '';
    document.body.style.border = '';
  });

  test('should apply a green border to the body', () => {
    handlePortal();
    expect(document.body.style.border).toBe('15px solid green');
  });

  test('should log a message to the console', () => {
    const logSpy = jest.spyOn(log, 'info').mockImplementation(() => {});
    handlePortal();
    expect(logSpy).toHaveBeenCalledWith('Portal page detected.');
    logSpy.mockRestore();
  });
});
