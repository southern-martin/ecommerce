import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockApiClient = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
}));

vi.mock('@/shared/lib/api-client', () => ({
  default: mockApiClient,
}));

import { sellerCouponApi } from '../seller-coupon.api';

describe('sellerCouponApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCoupons', () => {
    it('should fetch coupons with optional pagination params', async () => {
      const response = { data: [{ id: 'c1', code: 'SAVE10' }], total: 1 };
      mockApiClient.get.mockResolvedValue({ data: response });

      const result = await sellerCouponApi.getCoupons({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/coupons', { params: { page: 1, page_size: 10 } });
      expect(result).toEqual(response);
    });

    it('should work without params', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: [] } });

      await sellerCouponApi.getCoupons();

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/coupons', { params: undefined });
    });
  });

  describe('getCoupon', () => {
    it('should fetch a single coupon by id', async () => {
      const coupon = { id: 'c1', code: 'SAVE10', type: 'percentage', value: 10 };
      mockApiClient.get.mockResolvedValue({ data: { data: coupon } });

      const result = await sellerCouponApi.getCoupon('c1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/coupons/c1');
      expect(result).toEqual(coupon);
    });

    it('should handle response without data wrapper', async () => {
      const coupon = { id: 'c1', code: 'SAVE10', type: 'percentage', value: 10 };
      mockApiClient.get.mockResolvedValue({ data: coupon });

      const result = await sellerCouponApi.getCoupon('c1');

      expect(result).toEqual(coupon);
    });
  });

  describe('createCoupon', () => {
    it('should create a coupon and return it', async () => {
      const newCoupon = { code: 'NEW20', type: 'percentage', value: 20 };
      const created = { id: 'c2', ...newCoupon };
      mockApiClient.post.mockResolvedValue({ data: { data: created } });

      const result = await sellerCouponApi.createCoupon(newCoupon);

      expect(mockApiClient.post).toHaveBeenCalledWith('/seller/coupons', newCoupon);
      expect(result.id).toBe('c2');
    });
  });

  describe('updateCoupon', () => {
    it('should update a coupon by id', async () => {
      const updateData = { value: 25 };
      const updated = { id: 'c1', code: 'SAVE10', type: 'percentage', value: 25 };
      mockApiClient.patch.mockResolvedValue({ data: { data: updated } });

      const result = await sellerCouponApi.updateCoupon('c1', updateData);

      expect(mockApiClient.patch).toHaveBeenCalledWith('/seller/coupons/c1', updateData);
      expect(result.value).toBe(25);
    });
  });

  describe('deleteCoupon', () => {
    it('should delete a coupon by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await sellerCouponApi.deleteCoupon('c1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/seller/coupons/c1');
    });
  });
});
