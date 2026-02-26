import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';
import type { Product } from '@/modules/shop/types/shop.types';

export interface FlashSale {
  id: string;
  title: string;
  description: string;
  products: Product[];
  discount_percentage: number;
  starts_at: string;
  ends_at: string;
  is_active: boolean;
}

export const flashSaleApi = {
  getActiveFlashSales: async (): Promise<FlashSale[]> => {
    const response = await apiClient.get<ApiResponse<FlashSale[]>>('/flash-sales/active');
    return response.data.data;
  },

  getFlashSaleById: async (id: string): Promise<FlashSale> => {
    const response = await apiClient.get<ApiResponse<FlashSale>>(`/flash-sales/${id}`);
    return response.data.data;
  },
};
