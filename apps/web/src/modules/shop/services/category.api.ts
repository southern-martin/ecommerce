import apiClient from '@/shared/lib/api-client';
import type { Category } from '../types/shop.types';

export const categoryApi = {
  getCategories: async (): Promise<Category[]> => {
    const response = await apiClient.get('/categories');
    // Backend returns {categories: [...]}
    return response.data.categories || response.data.data || [];
  },

  getCategoryBySlug: async (slug: string): Promise<Category> => {
    const response = await apiClient.get(`/categories/${slug}`);
    return response.data.data || response.data;
  },
};
