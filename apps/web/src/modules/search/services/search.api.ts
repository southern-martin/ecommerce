import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';
import type { Product, FilterState } from '@/modules/shop/types/shop.types';

export interface SearchSuggestion {
  text: string;
  type: 'product' | 'category' | 'brand';
}

export const searchApi = {
  search: async (
    query: string,
    filters?: Partial<FilterState>
  ): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get<PaginatedResponse<Product>>('/search', {
      params: { q: query, ...filters },
    });
    return response.data;
  },

  suggest: async (query: string): Promise<SearchSuggestion[]> => {
    const response = await apiClient.get<ApiResponse<SearchSuggestion[]>>(
      '/search/suggestions',
      { params: { q: query } }
    );
    return response.data.data;
  },
};
