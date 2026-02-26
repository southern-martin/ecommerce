import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';
import type { Product } from '@/modules/shop/types/shop.types';

export interface AIChatMessage {
  role: 'user' | 'assistant';
  content: string;
}

export interface AIChatResponse {
  message: string;
  suggested_products?: Product[];
}

export const aiApi = {
  chat: async (messages: AIChatMessage[]): Promise<AIChatResponse> => {
    const response = await apiClient.post<ApiResponse<AIChatResponse>>('/ai/chat', { messages });
    return response.data.data;
  },

  getRecommendations: async (params?: {
    product_id?: string;
    category?: string;
    limit?: number;
  }): Promise<Product[]> => {
    const response = await apiClient.get<ApiResponse<Product[]>>('/ai/recommendations', { params });
    return response.data.data;
  },
};
