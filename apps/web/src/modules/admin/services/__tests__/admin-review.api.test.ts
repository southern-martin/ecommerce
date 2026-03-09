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

import { adminReviewApi } from '../admin-review.api';

describe('adminReviewApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getReviews', () => {
    it('should fetch reviews with default params', async () => {
      const mockData = { items: [], total: 0, page: 1, page_size: 10 };
      mockApiClient.get.mockResolvedValue({ data: mockData });

      const result = await adminReviewApi.getReviews();

      expect(mockApiClient.get).toHaveBeenCalledWith('/reviews', { params: undefined });
      expect(result).toEqual(mockData);
    });

    it('should fetch reviews with pagination params', async () => {
      const mockData = { items: [], total: 0, page: 2, page_size: 5 };
      mockApiClient.get.mockResolvedValue({ data: mockData });

      const result = await adminReviewApi.getReviews({ page: 2, page_size: 5 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/reviews', {
        params: { page: 2, page_size: 5 },
      });
      expect(result).toEqual(mockData);
    });
  });

  describe('approveReview', () => {
    it('should patch review to approve and return updated review', async () => {
      const review = {
        id: 'r1',
        product_id: 'p1',
        product_name: 'Widget',
        user_id: 'u1',
        user_name: 'John',
        rating: 5,
        comment: 'Great!',
        status: 'approved',
        created_at: '2024-01-01',
      };
      mockApiClient.patch.mockResolvedValue({ data: { data: review } });

      const result = await adminReviewApi.approveReview('r1');

      expect(mockApiClient.patch).toHaveBeenCalledWith('/admin/reviews/r1/approve');
      expect(result).toEqual(review);
    });
  });
});
