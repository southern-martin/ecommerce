import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Order } from '@/modules/checkout/services/order.api';

export interface SellerDashboardStats {
  total_revenue: number;
  total_orders: number;
  total_products: number;
  pending_orders: number;
  revenue_trend: number;
  orders_trend: number;
}

export const sellerOrderApi = {
  getOrders: async (params: { page: number; page_size: number; status?: string }): Promise<PaginatedResponse<Order>> => {
    const response = await apiClient.get<PaginatedResponse<Order>>('/seller/orders', { params });
    return response.data;
  },

  getDashboardStats: async (): Promise<SellerDashboardStats> => {
    const response = await apiClient.get<{ data: SellerDashboardStats }>('/seller/dashboard/stats');
    return response.data.data;
  },
};
