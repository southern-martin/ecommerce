import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminProductMgmtApi } from '../services/admin-product-mgmt.api';
import type { AdminProductFilter, CreateProductPayload } from '../services/admin-product-mgmt.api';

export function useAdminProductList(filters: AdminProductFilter = {}) {
  return useQuery({
    queryKey: ['admin-product-list', filters],
    queryFn: () => adminProductMgmtApi.listProducts(filters),
    placeholderData: keepPreviousData,
  });
}

export function useAdminCreateProduct() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateProductPayload) => adminProductMgmtApi.createProduct(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-list'] });
    },
  });
}

export function useAdminUpdateProduct() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateProductPayload> }) =>
      adminProductMgmtApi.updateProduct(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-list'] });
    },
  });
}

export function useAdminDeleteProduct() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminProductMgmtApi.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-list'] });
    },
  });
}
