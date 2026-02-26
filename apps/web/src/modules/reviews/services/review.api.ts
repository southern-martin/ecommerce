import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';

export interface Review {
  id: string;
  product_id: string;
  user_id: string;
  user_name: string;
  user_avatar?: string;
  rating: number;
  title: string;
  comment: string;
  images?: string[];
  helpful_count: number;
  created_at: string;
}

export interface CreateReviewData {
  product_id: string;
  rating: number;
  title: string;
  comment: string;
  images?: string[];
}

export const reviewApi = {
  getProductReviews: async (
    productId: string,
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<Review>> => {
    const response = await apiClient.get<PaginatedResponse<Review>>(
      `/products/${productId}/reviews`,
      { params }
    );
    return response.data;
  },

  createReview: async (data: CreateReviewData): Promise<Review> => {
    const response = await apiClient.post<ApiResponse<Review>>('/reviews', data);
    return response.data.data;
  },

  markHelpful: async (reviewId: string): Promise<void> => {
    await apiClient.post(`/reviews/${reviewId}/helpful`);
  },
};
