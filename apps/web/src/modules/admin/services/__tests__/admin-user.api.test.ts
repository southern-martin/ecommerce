import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockApiClient = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
}));

vi.mock('@/shared/lib/api-client', () => ({
  default: mockApiClient,
}));

import { adminUserApi } from '../admin-user.api';

describe('adminUserApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getUsers', () => {
    it('should fetch users with pagination params', async () => {
      const mockResponse = {
        data: { items: [], total: 0, page: 1, page_size: 10 },
      };
      mockApiClient.get.mockResolvedValue(mockResponse);

      const result = await adminUserApi.getUsers({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users', {
        params: { page: 1, page_size: 10 },
      });
      expect(result).toEqual(mockResponse.data);
    });

    it('should pass role filter when provided', async () => {
      const mockResponse = { data: { items: [], total: 0 } };
      mockApiClient.get.mockResolvedValue(mockResponse);

      await adminUserApi.getUsers({ page: 1, page_size: 10, role: 'admin' });

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/users', {
        params: { page: 1, page_size: 10, role: 'admin' },
      });
    });
  });

  describe('updateUser', () => {
    it('should patch user and return updated user data', async () => {
      const updatedUser = { id: 'u1', role: 'admin', is_verified: true };
      mockApiClient.patch.mockResolvedValue({ data: { data: updatedUser } });

      const result = await adminUserApi.updateUser('u1', { role: 'admin' });

      expect(mockApiClient.patch).toHaveBeenCalledWith('/admin/users/u1', { role: 'admin' });
      expect(result).toEqual(updatedUser);
    });
  });

  describe('deleteUser', () => {
    it('should delete user by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await adminUserApi.deleteUser('u1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/admin/users/u1');
    });
  });

  describe('getDashboardStats', () => {
    it('should fetch dashboard stats and unwrap response', async () => {
      const stats = {
        total_users: 100,
        total_orders: 50,
        total_revenue: 99900,
        total_products: 25,
        new_users_today: 5,
        orders_today: 3,
      };
      mockApiClient.get.mockResolvedValue({ data: { data: stats } });

      const result = await adminUserApi.getDashboardStats();

      expect(mockApiClient.get).toHaveBeenCalledWith('/admin/dashboard/stats');
      expect(result).toEqual(stats);
    });
  });
});
