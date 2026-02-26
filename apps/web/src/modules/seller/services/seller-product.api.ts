import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';
import type { Product } from '@/modules/shop/types/shop.types';

export interface CreateProductData {
  name: string;
  description: string;
  price: number;
  compare_at_price?: number;
  category_id: string;
  stock_quantity: number;
  images: { url: string; alt: string; is_primary: boolean }[];
  variants?: { name: string; value: string; price_modifier: number; stock_quantity: number }[];
}

export type UpdateProductData = Partial<CreateProductData>;

export const sellerProductApi = {
  getProducts: async (params: { page: number; page_size: number }): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get<PaginatedResponse<Product>>('/seller/products', { params });
    return response.data;
  },

  getProductById: async (id: string): Promise<Product> => {
    const response = await apiClient.get<ApiResponse<Product>>(`/seller/products/${id}`);
    return response.data.data;
  },

  createProduct: async (data: CreateProductData): Promise<Product> => {
    const response = await apiClient.post<ApiResponse<Product>>('/seller/products', data);
    return response.data.data;
  },

  updateProduct: async (id: string, data: UpdateProductData): Promise<Product> => {
    const response = await apiClient.patch<ApiResponse<Product>>(`/seller/products/${id}`, data);
    return response.data.data;
  },

  deleteProduct: async (id: string): Promise<void> => {
    await apiClient.delete(`/seller/products/${id}`);
  },
};
