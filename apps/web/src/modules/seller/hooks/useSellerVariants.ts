import { useMutation, useQueryClient } from '@tanstack/react-query';
import { sellerVariantApi } from '../services/seller-variant.api';

export function useAddOption() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, data }: { productId: string; data: { name: string; values: string[] } }) =>
      sellerVariantApi.addOption(productId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-products'] });
    },
  });
}

export function useRemoveOption() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, optionId }: { productId: string; optionId: string }) =>
      sellerVariantApi.removeOption(productId, optionId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-products'] });
    },
  });
}

export function useGenerateVariants() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (productId: string) => sellerVariantApi.generateVariants(productId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-products'] });
    },
  });
}

export function useUpdateVariant() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, variantId, data }: { productId: string; variantId: string; data: Partial<any> }) =>
      sellerVariantApi.updateVariant(productId, variantId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-products'] });
    },
  });
}

export function useUpdateVariantStock() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ productId, variantId, stock }: { productId: string; variantId: string; stock: number }) =>
      sellerVariantApi.updateVariantStock(productId, variantId, stock),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-products'] });
    },
  });
}
