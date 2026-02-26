import apiClient from '@/shared/lib/api-client';

export interface ProductOption {
  id: string;
  name: string;
  values: string[];
  sort_order: number;
}

export interface Variant {
  id: string;
  product_id: string;
  sku: string;
  name: string;
  price_cents: number;
  compare_at_cents: number;
  cost_cents: number;
  stock: number;
  is_active: boolean;
  option_values: { option_name: string; value: string }[];
}

export const sellerVariantApi = {
  addOption: async (productId: string, data: { name: string; values: string[] }): Promise<ProductOption> => {
    const response = await apiClient.post(`/seller/products/${productId}/options`, data);
    return response.data.data ?? response.data;
  },

  removeOption: async (productId: string, optionId: string): Promise<void> => {
    await apiClient.delete(`/seller/products/${productId}/options/${optionId}`);
  },

  generateVariants: async (productId: string): Promise<Variant[]> => {
    const response = await apiClient.post(`/seller/products/${productId}/variants/generate`);
    return response.data.data ?? response.data;
  },

  updateVariant: async (productId: string, variantId: string, data: Partial<Variant>): Promise<Variant> => {
    const response = await apiClient.patch(`/seller/products/${productId}/variants/${variantId}`, data);
    return response.data.data ?? response.data;
  },

  updateVariantStock: async (productId: string, variantId: string, stock: number): Promise<void> => {
    await apiClient.patch(`/seller/products/${productId}/variants/${variantId}/stock`, { stock });
  },
};
