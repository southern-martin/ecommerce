import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { reviewApi } from '../review.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
  post: ReturnType<typeof vi.fn>;
};

describe('reviewApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getProductReviews', () => {
    it('sends GET to /products/:id/reviews with pagination', async () => {
      const mockReviews = [
        { id: 'rev1', product_id: 'p1', user_id: 'u1', user_name: 'John', rating: 5, title: 'Great!', comment: 'Love it', helpful_count: 3, created_at: '' },
      ];
      mockApiClient.get.mockResolvedValue({
        data: { data: mockReviews, total: 1, page: 1, page_size: 10 },
      });

      const result = await reviewApi.getProductReviews('p1', { page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/products/p1/reviews', { params: { page: 1, page_size: 10 } });
      expect(result.data).toHaveLength(1);
      expect(result.data[0].rating).toBe(5);
    });

    it('returns paginated response structure', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { data: [], total: 42, page: 3, page_size: 10 },
      });

      const result = await reviewApi.getProductReviews('p1', { page: 3, page_size: 10 });

      expect(result.total).toBe(42);
      expect(result.page).toBe(3);
    });
  });

  describe('createReview', () => {
    it('sends POST to /reviews with review data', async () => {
      const reviewData = {
        product_id: 'p1',
        rating: 4,
        title: 'Good product',
        comment: 'Works well for the price',
        images: ['https://example.com/img1.jpg'],
      };
      const mockReview = {
        id: 'rev2',
        ...reviewData,
        user_id: 'u1',
        user_name: 'Jane',
        helpful_count: 0,
        created_at: '2025-03-01',
      };
      mockApiClient.post.mockResolvedValue({ data: { data: mockReview } });

      const result = await reviewApi.createReview(reviewData);

      expect(mockApiClient.post).toHaveBeenCalledWith('/reviews', reviewData);
      expect(result.id).toBe('rev2');
      expect(result.rating).toBe(4);
    });
  });

  describe('markHelpful', () => {
    it('sends POST to /reviews/:id/helpful', async () => {
      mockApiClient.post.mockResolvedValue({});

      await reviewApi.markHelpful('rev1');

      expect(mockApiClient.post).toHaveBeenCalledWith('/reviews/rev1/helpful');
    });
  });
});
