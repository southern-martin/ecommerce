import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

vi.mock('../../services/seller-product.api', () => ({
  sellerProductApi: {
    getProducts: vi.fn(),
    createProduct: vi.fn(),
    updateProduct: vi.fn(),
    deleteProduct: vi.fn(),
  },
}));

import { sellerProductApi } from '../../services/seller-product.api';
import { useSellerProducts, useCreateProduct, useUpdateProduct, useDeleteProduct } from '../useSellerProducts';

const mockGetProducts = sellerProductApi.getProducts as ReturnType<typeof vi.fn>;
const mockCreateProduct = sellerProductApi.createProduct as ReturnType<typeof vi.fn>;
const mockDeleteProduct = sellerProductApi.deleteProduct as ReturnType<typeof vi.fn>;

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useSellerProducts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch products with default pagination', async () => {
    const mockData = { data: [{ id: 'p1', name: 'Widget' }], total: 1, page: 1, page_size: 10 };
    mockGetProducts.mockResolvedValue(mockData);

    const { result } = renderHook(() => useSellerProducts(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockGetProducts).toHaveBeenCalledWith({ page: 1, page_size: 10 });
    expect(result.current.data).toEqual(mockData);
  });

  it('should pass custom page and pageSize', async () => {
    mockGetProducts.mockResolvedValue({ data: [], total: 0, page: 2, page_size: 5 });

    const { result } = renderHook(() => useSellerProducts(2, 5), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockGetProducts).toHaveBeenCalledWith({ page: 2, page_size: 5 });
  });
});

describe('useCreateProduct', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should call createProduct on mutate', async () => {
    const newProduct = { id: 'p3', name: 'New' };
    mockCreateProduct.mockResolvedValue(newProduct);

    const { result } = renderHook(() => useCreateProduct(), { wrapper: createWrapper() });

    result.current.mutate({ name: 'New', description: 'desc', base_price_cents: 100, category_id: 'c1' });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockCreateProduct).toHaveBeenCalledWith({ name: 'New', description: 'desc', base_price_cents: 100, category_id: 'c1' });
  });
});

describe('useDeleteProduct', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should call deleteProduct on mutate', async () => {
    mockDeleteProduct.mockResolvedValue(undefined);

    const { result } = renderHook(() => useDeleteProduct(), { wrapper: createWrapper() });

    result.current.mutate('p1');

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockDeleteProduct).toHaveBeenCalledWith('p1');
  });
});
