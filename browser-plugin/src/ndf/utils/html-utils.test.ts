import { stripHtml } from './html-utils';

describe('stripHtml', () => {
  it('should return empty string for empty input', () => {
    expect(stripHtml('')).toBe('');
  });

  it('should return original string if no html tags are present', () => {
    const text = 'Just some plain text';
    expect(stripHtml(text)).toBe(text);
  });

  it('should strip simple html tags', () => {
    const input = '<p>Hello <b>World</b></p>';
    expect(stripHtml(input)).toBe('Hello World');
  });

  it('should handle mixed content similar to the reported issue', () => {
    const input = 'Some text here. (<a href="http://example.com">Read more</a>)';
    expect(stripHtml(input)).toBe('Some text here. (Read more)');
  });

  it('should handle complex nested structures', () => {
    const input = '<div><h1>Title</h1><p>Body with <span style="color: red">style</span></p></div>';
    expect(stripHtml(input)).toBe('TitleBody with style');
  });

  it('should decode entities', () => {
    const input = 'Tom &amp; Jerry';
    expect(stripHtml(input)).toBe('Tom & Jerry');
  });
});
