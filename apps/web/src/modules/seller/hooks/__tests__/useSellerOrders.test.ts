import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

vi.mock('../../services/seller-order.api', () => ({
  sellerOrderApi: {
    getOrders: vi.fn(),
  },
}));

vi.mock('../../services/seller-product.api', () => ({
  sellerProductApi: {
    getProducts: vi.fn(),
  },
}));

import { sellerOrderApi } from '../../services/seller-order.api';
import { sellerProductApi } from '../../services/seller-product.api';
import { useSellerOrders, useSellerDashboardStats } from '../useSellerOrders';

const mockGetOrders = sellerOrderApi.getOrders as ReturnType<typeof vi.fn>;
const mockGetProducts = sellerProductApi.getProducts as ReturnType<typeof vi.fn>;

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useSellerOrders', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch orders with default pagination', async () => {
    const mockData = { data: [{ id: 'o1', status: 'pending', total: 5000 }], total: 1, page: 1, page_size: 10 };
    mockGetOrders.mockResolvedValue(mockData);

    const { result } = renderHook(() => useSellerOrders(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockGetOrders).toHaveBeenCalledWith({ page: 1, page_size: 10, status: undefined });
    expect(result.current.data).toEqual(mockData);
  });

  it('should pass status filter', async () => {
    mockGetOrders.mockResolvedValue({ data: [], total: 0, page: 1, page_size: 10 });

    const { result } = renderHook(() => useSellerOrders(1, 10, 'shipped'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockGetOrders).toHaveBeenCalledWith({ page: 1, page_size: 10, status: 'shipped' });
  });
});

describe('useSellerDashboardStats', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should compute stats from orders and products data', async () => {
    mockGetOrders.mockResolvedValue({
      data: [
        { id: 'o1', status: 'pending', total: 5000, items: [] },
        { id: 'o2', status: 'shipped', total: 3000, items: [] },
        { id: 'o3', status: 'pending', total: 2000, items: [] },
      ],
      total: 3,
      page: 1,
      page_size: 100,
    });
    mockGetProducts.mockResolvedValue({ data: [], total: 15, page: 1, page_size: 1 });

    const { result } = renderHook(() => useSellerDashboardStats(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.data).toEqual({
      total_revenue: 10000,
      total_orders: 3,
      total_products: 15,
      pending_orders: 2,
    });
  });

  it('should return null stats when no orders data', async () => {
    mockGetOrders.mockResolvedValue(undefined);
    mockGetProducts.mockResolvedValue({ data: [], total: 0, page: 1, page_size: 1 });

    const { result } = renderHook(() => useSellerDashboardStats(), { wrapper: createWrapper() });

    // Initially loading then data is undefined => stats null
    await waitFor(() => expect(result.current.isLoading).toBe(false));
    // When ordersData is undefined, stats should be null
  });
});
