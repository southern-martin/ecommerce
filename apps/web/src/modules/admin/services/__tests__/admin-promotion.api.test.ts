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

import { adminPromotionApi } from '../admin-promotion.api';

describe('adminPromotionApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCoupons', () => {
    it('should fetch coupons and unwrap response', async () => {
      const coupons = [{ id: 'cp1', code: 'SAVE10', type: 'percentage', value: 10, used_count: 0, starts_at: '2024-01-01', expires_at: '2024-12-31', is_active: true, created_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { data: coupons } });

      const result = await adminPromotionApi.getCoupons();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/promotions/coupons');
      expect(result).toEqual(coupons);
    });
  });

  describe('createCoupon', () => {
    it('should post coupon data and return created coupon', async () => {
      const data = { code: 'NEW20', type: 'fixed_amount' as const, value: 2000, starts_at: '2024-01-01', expires_at: '2024-12-31' };
      const created = { id: 'cp2', ...data, used_count: 0, is_active: true, created_at: '2024-01-01' };
      mockApiClient.post.mockResolvedValue({ data: { data: created } });

      const result = await adminPromotionApi.createCoupon(data);

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/promotions/coupons', data);
      expect(result).toEqual(created);
    });
  });

  describe('updateCoupon', () => {
    it('should patch coupon by id', async () => {
      const updated = { id: 'cp1', code: 'SAVE15', type: 'percentage', value: 15 };
      mockApiClient.patch.mockResolvedValue({ data: { data: updated } });

      const result = await adminPromotionApi.updateCoupon('cp1', { value: 15 });

      expect(mockApiClient.patch).toHaveBeenCalledWith('/admin/promotions/coupons/cp1', { value: 15 });
      expect(result).toEqual(updated);
    });
  });

  describe('deleteCoupon', () => {
    it('should delete coupon by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminPromotionApi.deleteCoupon('cp1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/promotions/coupons/cp1');
    });
  });

  describe('getFlashSales', () => {
    it('should fetch flash sales and unwrap response', async () => {
      const sales = [{ id: 'fs1', name: 'Summer Sale', discount_percentage: 30, starts_at: '2024-06-01', ends_at: '2024-06-30', is_active: true, created_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { data: sales } });

      const result = await adminPromotionApi.getFlashSales();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/promotions/flash-sales');
      expect(result).toEqual(sales);
    });
  });

  describe('createFlashSale', () => {
    it('should post flash sale data', async () => {
      const data = { name: 'Winter Sale', discount_percentage: 25, starts_at: '2024-12-01', ends_at: '2024-12-31' };
      const created = { id: 'fs2', ...data, is_active: true, created_at: '2024-01-01' };
      mockApiClient.post.mockResolvedValue({ data: { data: created } });

      const result = await adminPromotionApi.createFlashSale(data);

      expect(mockApiClient.post).toHaveBeenCalledWith('/admin/promotions/flash-sales', data);
      expect(result).toEqual(created);
    });
  });

  describe('getBundles', () => {
    it('should fetch bundles and unwrap response', async () => {
      const bundles = [{ id: 'b1', name: 'Starter Pack', discount_percentage: 15, is_active: true, created_at: '2024-01-01' }];
      mockApiClient.get.mockResolvedValue({ data: { data: bundles } });

      const result = await adminPromotionApi.getBundles();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/promotions/bundles');
      expect(result).toEqual(bundles);
    });
  });

  describe('deleteBundle', () => {
    it('should delete bundle by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminPromotionApi.deleteBundle('b1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/promotions/bundles/b1');
    });
  });
});
