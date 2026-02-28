import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminProductMgmtApi } from '../services/admin-product-mgmt.api';
import type {
  AdminProductFilter,
  CreateProductPayload,
  AdminUpdateProductPayload,
} from '../services/admin-product-mgmt.api';

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
    mutationFn: ({ id, data }: { id: string; data: AdminUpdateProductPayload }) =>
      adminProductMgmtApi.updateProduct(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-list'] });
      queryClient.invalidateQueries({ queryKey: ['admin-product', variables.id] });
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

export function useAdminProduct(id: string) {
  return useQuery({
    queryKey: ['admin-product', id],
    queryFn: () => adminProductMgmtApi.getProduct(id),
    enabled: !!id,
  });
}

export function useAdminProductOptions(productId: string) {
  return useQuery({
    queryKey: ['admin-product-options', productId],
    queryFn: () => adminProductMgmtApi.getOptions(productId),
    enabled: !!productId,
  });
}

export function useAdminAddOption() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, data }: { productId: string; data: { name: string; sort_order?: number; values: { value: string; color_hex?: string; sort_order?: number }[] } }) =>
      adminProductMgmtApi.addOption(productId, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-options', variables.productId] });
      queryClient.invalidateQueries({ queryKey: ['admin-product', variables.productId] });
    },
  });
}

export function useAdminRemoveOption() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, optionId }: { productId: string; optionId: string }) =>
      adminProductMgmtApi.removeOption(productId, optionId),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-options', variables.productId] });
      queryClient.invalidateQueries({ queryKey: ['admin-product', variables.productId] });
    },
  });
}

export function useAdminProductVariants(productId: string) {
  return useQuery({
    queryKey: ['admin-product-variants', productId],
    queryFn: () => adminProductMgmtApi.getVariants(productId),
    enabled: !!productId,
  });
}

export function useAdminGenerateVariants() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (productId: string) => adminProductMgmtApi.generateVariants(productId),
    onSuccess: (_data, productId) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-variants', productId] });
      queryClient.invalidateQueries({ queryKey: ['admin-product', productId] });
    },
  });
}

export function useAdminUpdateVariant() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, variantId, data }: { productId: string; variantId: string; data: Record<string, unknown> }) =>
      adminProductMgmtApi.updateVariant(productId, variantId, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-variants', variables.productId] });
    },
  });
}

export function useAdminUpdateVariantStock() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, variantId, delta }: { productId: string; variantId: string; delta: number }) =>
      adminProductMgmtApi.updateVariantStock(productId, variantId, delta),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-variants', variables.productId] });
    },
  });
}

export function useAdminProductAttributes(productId: string) {
  return useQuery({
    queryKey: ['admin-product-attributes', productId],
    queryFn: () => adminProductMgmtApi.getProductAttributes(productId),
    enabled: !!productId,
  });
}

export function useAdminSetProductAttributes() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, attributes }: { productId: string; attributes: { attribute_id: string; value: string; values?: string[] }[] }) =>
      adminProductMgmtApi.setProductAttributes(productId, attributes),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-product-attributes', variables.productId] });
    },
  });
}
