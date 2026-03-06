import { describe, it, expect } from 'vitest';
import { cn, formatPrice, formatDate, formatDateTime, truncate } from '../utils';

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar');
  });

  it('handles conditional classes', () => {
    const shouldHide = false;
    expect(cn('base', shouldHide && 'hidden', 'visible')).toBe('base visible');
  });

  it('deduplicates tailwind classes', () => {
    const result = cn('p-4', 'p-2');
    expect(result).toBe('p-2');
  });
});

describe('formatPrice', () => {
  it('formats cents to dollars', () => {
    expect(formatPrice(1999)).toBe('$19.99');
  });

  it('formats zero', () => {
    expect(formatPrice(0)).toBe('$0.00');
  });

  it('formats large amounts', () => {
    expect(formatPrice(999999)).toBe('$9,999.99');
  });

  it('formats single cent', () => {
    expect(formatPrice(1)).toBe('$0.01');
  });

  it('uses specified currency', () => {
    const result = formatPrice(1000, 'EUR');
    expect(result).toContain('10.00');
  });
});

describe('formatDate', () => {
  it('formats ISO date string', () => {
    const result = formatDate('2026-01-15T10:00:00Z');
    expect(result).toContain('Jan');
    expect(result).toContain('15');
    expect(result).toContain('2026');
  });

  it('formats Date object', () => {
    const result = formatDate(new Date('2026-06-01'));
    expect(result).toContain('2026');
  });
});

describe('formatDateTime', () => {
  it('includes time components', () => {
    const result = formatDateTime('2026-03-05T14:30:00Z');
    expect(result).toContain('Mar');
    expect(result).toContain('5');
    expect(result).toContain('2026');
  });
});

describe('truncate', () => {
  it('returns short strings unchanged', () => {
    expect(truncate('hello', 10)).toBe('hello');
  });

  it('returns exact-length strings unchanged', () => {
    expect(truncate('hello', 5)).toBe('hello');
  });

  it('truncates long strings with ellipsis', () => {
    expect(truncate('hello world', 5)).toBe('hello...');
  });

  it('handles empty string', () => {
    expect(truncate('', 5)).toBe('');
  });
});
