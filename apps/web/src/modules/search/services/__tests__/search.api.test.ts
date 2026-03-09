import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { searchApi } from '../search.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
};

describe('searchApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('search', () => {
    it('sends GET to /search with query param', async () => {
      const mockResponse = {
        data: { data: [{ id: 'p1', name: 'Widget' }], total: 1, page: 1, page_size: 20 },
      };
      mockApiClient.get.mockResolvedValue(mockResponse);

      const result = await searchApi.search('widget');

      expect(mockApiClient.get).toHaveBeenCalledWith('/search', {
        params: { q: 'widget' },
      });
      expect(result.data).toHaveLength(1);
    });

    it('passes additional filters as params', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { data: [], total: 0, page: 1, page_size: 20 },
      });

      await searchApi.search('shoes', { category: 'footwear' } as any);

      expect(mockApiClient.get).toHaveBeenCalledWith('/search', {
        params: { q: 'shoes', category: 'footwear' },
      });
    });

    it('returns paginated response', async () => {
      const mockData = {
        data: { data: [{ id: 'p1' }, { id: 'p2' }], total: 50, page: 2, page_size: 20 },
      };
      mockApiClient.get.mockResolvedValue(mockData);

      const result = await searchApi.search('test');

      expect(result.total).toBe(50);
      expect(result.page).toBe(2);
    });
  });

  describe('suggest', () => {
    it('sends GET to /search/suggestions with query param', async () => {
      const suggestions = [
        { text: 'widget pro', type: 'product' },
        { text: 'widgets category', type: 'category' },
      ];
      mockApiClient.get.mockResolvedValue({ data: { data: suggestions } });

      const result = await searchApi.suggest('wid');

      expect(mockApiClient.get).toHaveBeenCalledWith('/search/suggestions', {
        params: { q: 'wid' },
      });
      expect(result).toHaveLength(2);
      expect(result[0].text).toBe('widget pro');
      expect(result[0].type).toBe('product');
    });
  });
});
