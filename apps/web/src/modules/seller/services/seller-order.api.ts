import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Order } from '@/modules/checkout/services/order.api';

export const sellerOrderApi = {
  getOrders: async (params: { page: number; page_size: number; status?: string }): Promise<PaginatedResponse<Order>> => {
    const response = await apiClient.get<PaginatedResponse<Order>>('/seller/orders', { params });
    return response.data;
  },
};
