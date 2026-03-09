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

import { sellerOrderApi } from '../seller-order.api';

describe('sellerOrderApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getOrders', () => {
    it('should fetch orders with pagination params', async () => {
      const mockResponse = { data: [], total: 0, page: 1, page_size: 10 };
      mockApiClient.get.mockResolvedValue({ data: mockResponse });

      const result = await sellerOrderApi.getOrders({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/orders', { params: { page: 1, page_size: 10 } });
      expect(result).toEqual(mockResponse);
    });

    it('should pass status filter when provided', async () => {
      const mockResponse = { data: [], total: 0, page: 1, page_size: 10 };
      mockApiClient.get.mockResolvedValue({ data: mockResponse });

      await sellerOrderApi.getOrders({ page: 1, page_size: 10, status: 'pending' });

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/orders', {
        params: { page: 1, page_size: 10, status: 'pending' },
      });
    });

    it('should return order data from response', async () => {
      const orders = [
        { id: 'o1', order_number: 'ORD-001', status: 'pending', total: 5000, items: [], created_at: '2026-01-01' },
      ];
      mockApiClient.get.mockResolvedValue({ data: { data: orders, total: 1, page: 1, page_size: 10 } });

      const result = await sellerOrderApi.getOrders({ page: 1, page_size: 10 });

      expect(result.data).toHaveLength(1);
      expect(result.total).toBe(1);
    });
  });
});
