import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';

export interface Notification {
  id: string;
  type: 'order' | 'promotion' | 'system' | 'review' | 'shipping';
  title: string;
  message: string;
  is_read: boolean;
  action_url?: string;
  created_at: string;
}

export const notificationApi = {
  getNotifications: async (
    params: { page: number; page_size: number }
  ): Promise<PaginatedResponse<Notification>> => {
    const response = await apiClient.get<PaginatedResponse<Notification>>('/notifications', { params });
    return response.data;
  },

  getUnreadCount: async (): Promise<number> => {
    const response = await apiClient.get<ApiResponse<{ count: number }>>('/notifications/unread-count');
    return response.data.data.count;
  },

  markAsRead: async (id: string): Promise<void> => {
    await apiClient.patch(`/notifications/${id}/read`);
  },

  markAllAsRead: async (): Promise<void> => {
    await apiClient.patch('/notifications/read-all');
  },
};
