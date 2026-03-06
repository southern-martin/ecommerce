import { describe, it, expect } from 'vitest';
import { emailSchema, passwordSchema, phoneSchema, addressSchema, paginationSchema } from '../validators';

describe('emailSchema', () => {
  it('accepts valid email', () => {
    expect(emailSchema.safeParse('test@example.com').success).toBe(true);
  });

  it('rejects invalid email', () => {
    const result = emailSchema.safeParse('not-an-email');
    expect(result.success).toBe(false);
  });

  it('rejects empty string', () => {
    expect(emailSchema.safeParse('').success).toBe(false);
  });
});

describe('passwordSchema', () => {
  it('accepts 8+ character password', () => {
    expect(passwordSchema.safeParse('12345678').success).toBe(true);
  });

  it('rejects short password', () => {
    const result = passwordSchema.safeParse('1234567');
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toContain('8 characters');
    }
  });
});

describe('phoneSchema', () => {
  it('accepts valid phone with + prefix', () => {
    expect(phoneSchema.safeParse('+14155552671').success).toBe(true);
  });

  it('accepts empty string (optional)', () => {
    expect(phoneSchema.safeParse('').success).toBe(true);
  });

  it('accepts undefined (optional)', () => {
    expect(phoneSchema.safeParse(undefined).success).toBe(true);
  });

  it('rejects invalid phone', () => {
    expect(phoneSchema.safeParse('abc').success).toBe(false);
  });
});

describe('addressSchema', () => {
  const validAddress = {
    street: '123 Main St',
    city: 'San Francisco',
    state: 'CA',
    zipCode: '94102',
    country: 'US',
  };

  it('accepts valid address', () => {
    expect(addressSchema.safeParse(validAddress).success).toBe(true);
  });

  it('rejects missing street', () => {
    expect(addressSchema.safeParse({ ...validAddress, street: '' }).success).toBe(false);
  });

  it('rejects missing city', () => {
    expect(addressSchema.safeParse({ ...validAddress, city: '' }).success).toBe(false);
  });

  it('rejects missing all fields', () => {
    expect(addressSchema.safeParse({}).success).toBe(false);
  });
});

describe('paginationSchema', () => {
  it('accepts valid pagination', () => {
    const result = paginationSchema.safeParse({ page: 1, pageSize: 20 });
    expect(result.success).toBe(true);
  });

  it('applies defaults', () => {
    const result = paginationSchema.parse({});
    expect(result.page).toBe(1);
    expect(result.pageSize).toBe(20);
  });

  it('rejects page < 1', () => {
    expect(paginationSchema.safeParse({ page: 0 }).success).toBe(false);
  });

  it('rejects pageSize > 100', () => {
    expect(paginationSchema.safeParse({ pageSize: 101 }).success).toBe(false);
  });
});
