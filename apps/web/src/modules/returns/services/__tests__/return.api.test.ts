import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { returnApi } from '../return.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
  post: ReturnType<typeof vi.fn>;
};

describe('returnApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getReturns', () => {
    it('sends GET to /returns with pagination params', async () => {
      const mockReturns = [
        { id: 'r1', order_id: 'o1', order_number: 'ORD-001', reason: 'defective', status: 'pending', items: [], refund_amount: 2000, created_at: '', updated_at: '' },
      ];
      mockApiClient.get.mockResolvedValue({
        data: { data: mockReturns, total: 1, page: 1, page_size: 10 },
      });

      const result = await returnApi.getReturns({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/returns', { params: { page: 1, page_size: 10 } });
      expect(result.data).toHaveLength(1);
      expect(result.data[0].reason).toBe('defective');
    });
  });

  describe('createReturn', () => {
    it('sends POST to /returns with return data', async () => {
      const returnData = {
        order_id: 'o1',
        reason: 'wrong_item',
        description: 'Received the wrong color',
        items: [{ product_id: 'p1', quantity: 1 }],
      };
      const mockReturn = {
        id: 'r2',
        ...returnData,
        order_number: 'ORD-001',
        status: 'pending',
        items: [{ id: 'ri1', product_id: 'p1', product_name: 'Widget', quantity: 1, image_url: '' }],
        refund_amount: 1500,
        created_at: '2025-01-01',
        updated_at: '2025-01-01',
      };
      mockApiClient.post.mockResolvedValue({ data: { data: mockReturn } });

      const result = await returnApi.createReturn(returnData);

      expect(mockApiClient.post).toHaveBeenCalledWith('/returns', returnData);
      expect(result.id).toBe('r2');
      expect(result.status).toBe('pending');
    });
  });

  describe('getReturnById', () => {
    it('sends GET to /returns/:id and returns return request', async () => {
      const mockReturn = {
        id: 'r1',
        order_id: 'o1',
        order_number: 'ORD-001',
        reason: 'defective',
        description: 'Screen cracked',
        status: 'approved',
        items: [],
        refund_amount: 5000,
        created_at: '2025-01-01',
        updated_at: '2025-01-02',
      };
      mockApiClient.get.mockResolvedValue({ data: { data: mockReturn } });

      const result = await returnApi.getReturnById('r1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/returns/r1');
      expect(result.id).toBe('r1');
      expect(result.status).toBe('approved');
      expect(result.refund_amount).toBe(5000);
    });
  });
});
