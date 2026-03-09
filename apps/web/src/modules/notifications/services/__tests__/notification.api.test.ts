import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
    patch: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { notificationApi } from '../notification.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
  patch: ReturnType<typeof vi.fn>;
};

describe('notificationApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getNotifications', () => {
    it('sends GET to /notifications with pagination', async () => {
      const mockNotifications = [
        { id: 'n1', type: 'order', title: 'Order shipped', message: 'Your order has shipped', is_read: false, created_at: '' },
        { id: 'n2', type: 'promotion', title: 'Sale!', message: '50% off today', is_read: true, created_at: '' },
      ];
      mockApiClient.get.mockResolvedValue({
        data: { data: mockNotifications, total: 2, page: 1, page_size: 20 },
      });

      const result = await notificationApi.getNotifications({ page: 1, page_size: 20 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/notifications', { params: { page: 1, page_size: 20 } });
      expect(result.data).toHaveLength(2);
      expect(result.data[0].type).toBe('order');
    });
  });

  describe('getUnreadCount', () => {
    it('sends GET to /notifications/unread-count and returns count', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: { count: 5 } } });

      const result = await notificationApi.getUnreadCount();

      expect(mockApiClient.get).toHaveBeenCalledWith('/notifications/unread-count');
      expect(result).toBe(5);
    });

    it('returns 0 when no unread notifications', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: { count: 0 } } });

      const result = await notificationApi.getUnreadCount();

      expect(result).toBe(0);
    });
  });

  describe('markAsRead', () => {
    it('sends PATCH to /notifications/:id/read', async () => {
      mockApiClient.patch.mockResolvedValue({});

      await notificationApi.markAsRead('n1');

      expect(mockApiClient.patch).toHaveBeenCalledWith('/notifications/n1/read');
    });
  });

  describe('markAllAsRead', () => {
    it('sends PATCH to /notifications/read-all', async () => {
      mockApiClient.patch.mockResolvedValue({});

      await notificationApi.markAllAsRead();

      expect(mockApiClient.patch).toHaveBeenCalledWith('/notifications/read-all');
    });
  });
});
