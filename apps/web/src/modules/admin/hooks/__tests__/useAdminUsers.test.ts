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

import { useAdminUsers, useUpdateUser, useDeleteUser } from '../useAdminUsers';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useAdminUsers', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch users with default pagination', async () => {
    const mockData = { items: [{ id: 'u1', email: 'test@test.com' }], total: 1 };
    mockAdminUserApi.getUsers.mockResolvedValue(mockData);

    const { result } = renderHook(() => useAdminUsers(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminUserApi.getUsers).toHaveBeenCalledWith({ page: 1, page_size: 10 });
    expect(result.current.data).toEqual(mockData);
  });

  it('should fetch users with custom page and role', async () => {
    mockAdminUserApi.getUsers.mockResolvedValue({ items: [], total: 0 });

    const { result } = renderHook(() => useAdminUsers(2, 20, 'admin'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminUserApi.getUsers).toHaveBeenCalledWith({ page: 2, page_size: 20, role: 'admin' });
  });
});

describe('useUpdateUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should call updateUser with id and data', async () => {
    const updatedUser = { id: 'u1', role: 'admin' };
    mockAdminUserApi.updateUser.mockResolvedValue(updatedUser);

    const { result } = renderHook(() => useUpdateUser(), { wrapper: createWrapper() });

    result.current.mutate({ id: 'u1', data: { role: 'admin' } });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminUserApi.updateUser).toHaveBeenCalledWith('u1', { role: 'admin' });
  });
});

describe('useDeleteUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should call deleteUser with id', async () => {
    mockAdminUserApi.deleteUser.mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteUser(), { wrapper: createWrapper() });

    result.current.mutate('u1');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminUserApi.deleteUser).toHaveBeenCalledWith('u1');
  });
});
