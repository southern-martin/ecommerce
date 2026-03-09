import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { accountOrderApi } from '../order.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
};

describe('accountOrderApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getOrders', () => {
    it('sends GET to /orders with pagination params', async () => {
      const mockOrders = [
        {
          id: 'o1',
          order_number: 'ORD-001',
          status: 'delivered',
          items: [],
          subtotal_cents: 5000,
          shipping_cents: 500,
          tax_cents: 450,
          discount_cents: 0,
          total_cents: 5950,
          shipping_address: {},
          created_at: '2025-01-01',
        },
      ];
      mockApiClient.get.mockResolvedValue({
        data: { data: mockOrders, total: 1, page: 1, page_size: 10 },
      });

      const result = await accountOrderApi.getOrders({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/orders', { params: { page: 1, page_size: 10 } });
      expect(result.data).toHaveLength(1);
      expect(result.total).toBe(1);
      expect(result.data[0].order_number).toBe('ORD-001');
    });

    it('maps order items correctly', async () => {
      const mockOrders = [
        {
          id: 'o2',
          status: 'pending',
          items: [
            { id: 'i1', product_id: 'p1', product_name: 'Widget', image_url: 'img.jpg', unit_price_cents: 1000, quantity: 2 },
          ],
          subtotal_cents: 2000,
          total_cents: 2000,
          shipping_address: { full_name: 'John Doe', line1: '123 Main St', city: 'NYC', state: 'NY', postal_code: '10001', country_code: 'US' },
          created_at: '2025-02-01',
        },
      ];
      mockApiClient.get.mockResolvedValue({
        data: { data: mockOrders, total: 1, page: 1, page_size: 10 },
      });

      const result = await accountOrderApi.getOrders({ page: 1, page_size: 10 });
      const item = result.data[0].items[0];

      expect(item.name).toBe('Widget');
      expect(item.price).toBe(1000);
      expect(item.quantity).toBe(2);
    });

    it('handles fallback to orders key in response', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { orders: [{ id: 'o3', status: 'shipped', items: [], created_at: '' }] },
      });

      const result = await accountOrderApi.getOrders({ page: 1, page_size: 10 });

      expect(result.data).toHaveLength(1);
      expect(result.data[0].id).toBe('o3');
    });
  });

  describe('getOrderById', () => {
    it('sends GET to /orders/:id and returns mapped order', async () => {
      const mockOrder = {
        id: 'o1',
        order_number: 'ORD-001',
        status: 'delivered',
        items: [],
        subtotal_cents: 5000,
        total_cents: 5000,
        shipping_address: {},
        created_at: '2025-01-01',
      };
      mockApiClient.get.mockResolvedValue({ data: { data: mockOrder } });

      const result = await accountOrderApi.getOrderById('o1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/orders/o1');
      expect(result.id).toBe('o1');
      expect(result.order_number).toBe('ORD-001');
    });
  });
});
