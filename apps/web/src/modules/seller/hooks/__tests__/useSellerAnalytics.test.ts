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
import { useSellerAnalytics } from '../useSellerAnalytics';

const mockGetOrders = sellerOrderApi.getOrders as ReturnType<typeof vi.fn>;

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useSellerAnalytics', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should return null analytics when no order data', async () => {
    mockGetOrders.mockResolvedValue({ data: undefined, total: 0, page: 1, page_size: 100 });

    const { result } = renderHook(() => useSellerAnalytics(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isLoading).toBe(false));

    expect(result.current.analytics).toBeNull();
  });

  it('should compute status distribution from orders', async () => {
    mockGetOrders.mockResolvedValue({
      data: [
        { id: 'o1', status: 'pending', total: 1000, items: [], created_at: new Date().toISOString() },
        { id: 'o2', status: 'pending', total: 2000, items: [], created_at: new Date().toISOString() },
        { id: 'o3', status: 'shipped', total: 3000, items: [], created_at: new Date().toISOString() },
      ],
      total: 3,
      page: 1,
      page_size: 100,
    });

    const { result } = renderHook(() => useSellerAnalytics(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.analytics).not.toBeNull());

    const statusDist = result.current.analytics!.statusDistribution;
    expect(statusDist).toContainEqual({ name: 'Pending', value: 2 });
    expect(statusDist).toContainEqual({ name: 'Shipped', value: 1 });
  });

  it('should compute top products from order items', async () => {
    mockGetOrders.mockResolvedValue({
      data: [
        {
          id: 'o1',
          status: 'completed',
          total: 5000,
          created_at: new Date().toISOString(),
          items: [
            { id: 'i1', product_id: 'p1', name: 'Widget', price: 1000, quantity: 3, image_url: '' },
            { id: 'i2', product_id: 'p2', name: 'Gadget', price: 2000, quantity: 1, image_url: '' },
          ],
        },
      ],
      total: 1,
      page: 1,
      page_size: 100,
    });

    const { result } = renderHook(() => useSellerAnalytics(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.analytics).not.toBeNull());

    const topProducts = result.current.analytics!.topProducts;
    expect(topProducts[0]).toEqual({ name: 'Widget', revenue: 3000 });
    expect(topProducts[1]).toEqual({ name: 'Gadget', revenue: 2000 });
  });

  it('should compute revenue by day with 7 days', async () => {
    mockGetOrders.mockResolvedValue({
      data: [],
      total: 0,
      page: 1,
      page_size: 100,
    });

    const { result } = renderHook(() => useSellerAnalytics(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.analytics).not.toBeNull());

    expect(result.current.analytics!.revenueByDay).toHaveLength(7);
  });
});
