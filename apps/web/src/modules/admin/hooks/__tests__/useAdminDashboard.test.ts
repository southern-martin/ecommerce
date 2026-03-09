import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockAdminUserApi = vi.hoisted(() => ({
  getUsers: vi.fn(),
  updateUser: vi.fn(),
  deleteUser: vi.fn(),
  getDashboardStats: vi.fn(),
}));

vi.mock('../../services/admin-user.api', () => ({
  adminUserApi: mockAdminUserApi,
}));

import { useAdminDashboard } from '../useAdminDashboard';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useAdminDashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch dashboard stats', async () => {
    const stats = {
      total_users: 150,
      total_orders: 75,
      total_revenue: 500000,
      total_products: 40,
      new_users_today: 8,
      orders_today: 12,
    };
    mockAdminUserApi.getDashboardStats.mockResolvedValue(stats);

    const { result } = renderHook(() => useAdminDashboard(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(stats);
    expect(mockAdminUserApi.getDashboardStats).toHaveBeenCalledTimes(1);
  });

  it('should handle error state', async () => {
    mockAdminUserApi.getDashboardStats.mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useAdminDashboard(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isError).toBe(true));

    expect(result.current.error).toBeDefined();
  });

  it('should use admin-dashboard query key', async () => {
    mockAdminUserApi.getDashboardStats.mockResolvedValue({
      total_users: 0,
      total_orders: 0,
      total_revenue: 0,
      total_products: 0,
      new_users_today: 0,
      orders_today: 0,
    });

    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false, gcTime: 0 } },
    });
    const wrapper = ({ children }: { children: React.ReactNode }) =>
      React.createElement(QueryClientProvider, { client: queryClient }, children);

    const { result } = renderHook(() => useAdminDashboard(), { wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    const queryState = queryClient.getQueryState(['admin-dashboard']);
    expect(queryState).toBeDefined();
  });
});
