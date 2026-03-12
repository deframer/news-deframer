import { formatShortDate, formatTime } from './formatTime';

describe('formatTime', () => {
  beforeEach(() => {
    jest.useFakeTimers();
    jest.setSystemTime(new Date('2026-03-08T12:00:00Z'));
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('formats English relative time', () => {
    expect(formatTime('2026-03-08T11:59:30Z', 'en')).toBe('30 seconds ago');
  });

  it('formats German relative time', () => {
    expect(formatTime('2026-03-08T11:59:30Z', 'de')).toBe('vor 30 Sekunden');
  });

  it('keeps relative time for older English dates', () => {
    expect(formatTime('2026-03-04T12:00:00Z', 'en')).toBe('4 days ago');
  });

  it('keeps relative time for older German dates', () => {
    expect(formatTime('2026-03-04T12:00:00Z', 'de')).toBe('vor 4 Tagen');
  });

  it('formats short dates like lifecycle/tag cloud contexts', () => {
    expect(formatShortDate('2026-03-12T12:20:43+01:00', 'en')).toBe('Mar 12');
  });

  it('returns empty string for invalid dates', () => {
    expect(formatTime('not-a-date', 'en')).toBe('');
  });
});
