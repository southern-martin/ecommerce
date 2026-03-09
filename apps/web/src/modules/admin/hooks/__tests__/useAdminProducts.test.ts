import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const mockAdminProductApi = vi.hoisted(() => ({
  getCategories: vi.fn(),
  createCategory: vi.fn(),
  updateCategory: vi.fn(),
  deleteCategory: vi.fn(),
  getAttributes: vi.fn(),
  createAttribute: vi.fn(),
  updateAttribute: vi.fn(),
  deleteAttribute: vi.fn(),
  getAttributeGroups: vi.fn(),
  createAttributeGroup: vi.fn(),
  updateAttributeGroup: vi.fn(),
  deleteAttributeGroup: vi.fn(),
  getGroupAttributes: vi.fn(),
  addAttributeToGroup: vi.fn(),
  removeAttributeFromGroup: vi.fn(),
}));

vi.mock('../../services/admin-product.api', () => ({
  adminProductApi: mockAdminProductApi,
}));

import {
  useCategories,
  useCreateCategory,
  useAdminAttributes,
  useAttributeGroups,
  useGroupAttributes,
} from '../useAdminProducts';

function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false, gcTime: 0 } },
  });
  return ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);
}

describe('useCategories', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch categories', async () => {
    const categories = [{ id: 'c1', name: 'Electronics', slug: 'electronics', created_at: '2024-01-01' }];
    mockAdminProductApi.getCategories.mockResolvedValue(categories);

    const { result } = renderHook(() => useCategories(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(categories);
  });
});

describe('useCreateCategory', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should call createCategory mutation', async () => {
    const created = { id: 'c2', name: 'Books', slug: 'books', created_at: '2024-01-01' };
    mockAdminProductApi.createCategory.mockResolvedValue(created);

    const { result } = renderHook(() => useCreateCategory(), { wrapper: createWrapper() });

    result.current.mutate({ name: 'Books' });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminProductApi.createCategory).toHaveBeenCalledWith({ name: 'Books' });
  });
});

describe('useAdminAttributes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch attributes', async () => {
    const attributes = [{ id: 'a1', name: 'Color', type: 'select', required: true, filterable: true, created_at: '2024-01-01' }];
    mockAdminProductApi.getAttributes.mockResolvedValue(attributes);

    const { result } = renderHook(() => useAdminAttributes(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(attributes);
  });
});

describe('useAttributeGroups', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch attribute groups', async () => {
    const groups = [{ id: 'g1', name: 'Physical', slug: 'physical', sort_order: 1, created_at: '2024-01-01', updated_at: '2024-01-01' }];
    mockAdminProductApi.getAttributeGroups.mockResolvedValue(groups);

    const { result } = renderHook(() => useAttributeGroups(), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual(groups);
  });
});

describe('useGroupAttributes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch group attributes when groupId is provided', async () => {
    const attrs = [{ id: 'a1', name: 'Weight', type: 'number', required: false, filterable: false, created_at: '2024-01-01' }];
    mockAdminProductApi.getGroupAttributes.mockResolvedValue(attrs);

    const { result } = renderHook(() => useGroupAttributes('g1'), { wrapper: createWrapper() });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(mockAdminProductApi.getGroupAttributes).toHaveBeenCalledWith('g1');
    expect(result.current.data).toEqual(attrs);
  });
});
