import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';
import type { User } from '@/shared/types/user.types';

export interface UpdateUserData {
  role?: string;
  is_verified?: boolean;
}

export interface AdminDashboardStats {
  total_users: number;
  total_orders: number;
  total_revenue: number;
  total_products: number;
  new_users_today: number;
  orders_today: number;
}

export const adminUserApi = {
  getUsers: async (params: { page: number; page_size: number; role?: string }): Promise<PaginatedResponse<User>> => {
    const response = await apiClient.get<PaginatedResponse<User>>('/admin/users', { params });
    return response.data;
  },

  updateUser: async (id: string, data: UpdateUserData): Promise<User> => {
    const response = await apiClient.patch<ApiResponse<User>>(`/admin/users/${id}`, data);
    return response.data.data;
  },

  deleteUser: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/users/${id}`);
  },

  getDashboardStats: async (): Promise<AdminDashboardStats> => {
    const response = await apiClient.get<ApiResponse<AdminDashboardStats>>('/admin/dashboard/stats');
    return response.data.data;
  },
};
