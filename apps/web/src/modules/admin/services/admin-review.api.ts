import apiClient from '@/shared/lib/api-client';
import type { ApiResponse, PaginatedResponse } from '@/shared/types/api.types';

export interface Review {
  id: string;
  product_id: string;
  product_name: string;
  user_id: string;
  user_name: string;
  rating: number;
  comment: string;
  status: 'pending' | 'approved' | 'rejected';
  created_at: string;
}

export const adminReviewApi = {
  getReviews: async (params?: {
    page?: number;
    page_size?: number;
  }): Promise<PaginatedResponse<Review>> => {
    const response = await apiClient.get<PaginatedResponse<Review>>('/reviews', { params });
    return response.data;
  },

  approveReview: async (id: string): Promise<Review> => {
    const response = await apiClient.patch<ApiResponse<Review>>(`/admin/reviews/${id}/approve`);
    return response.data.data;
  },
};
