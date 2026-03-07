import { describe, it, expect, vi, beforeEach } from 'vitest';
import { categoryApi } from '../category.api';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
};

const mockCategories = [
  { id: 'cat-1', name: 'Electronics', slug: 'electronics' },
  { id: 'cat-2', name: 'Fashion', slug: 'fashion' },
];

describe('categoryApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCategories', () => {
    it('fetches categories from /categories', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { categories: mockCategories },
      });

      const result = await categoryApi.getCategories();

      expect(mockApiClient.get).toHaveBeenCalledWith('/categories');
      expect(result).toHaveLength(2);
      expect(result[0].name).toBe('Electronics');
    });

    it('handles alternative response shape (data.data)', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { data: mockCategories },
      });

      const result = await categoryApi.getCategories();
      expect(result).toHaveLength(2);
    });

    it('returns empty array when no categories', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { categories: [] },
      });

      const result = await categoryApi.getCategories();
      expect(result).toHaveLength(0);
    });
  });

  describe('getCategoryBySlug', () => {
    it('fetches a single category by slug', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { data: mockCategories[0] },
      });

      const result = await categoryApi.getCategoryBySlug('electronics');

      expect(mockApiClient.get).toHaveBeenCalledWith('/categories/electronics');
      expect(result.name).toBe('Electronics');
    });
  });
});
