import apiClient from '@/shared/lib/api-client';

export interface ProductOptionValue {
  id: string;
  option_id: string;
  value: string;
  color_hex?: string;
  sort_order: number;
}

export interface ProductOption {
  id: string;
  product_id?: string;
  name: string;
  sort_order: number;
  values: ProductOptionValue[];
}

export interface VariantOptionValue {
  variant_id: string;
  option_id: string;
  option_value_id: string;
  option_name: string;
  value: string;
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
  is_default: boolean;
  is_active: boolean;
  weight_grams: number;
  barcode: string;
  image_urls: string[];
  option_values: VariantOptionValue[];
  created_at?: string;
  updated_at?: string;
}

export const sellerVariantApi = {
  addOption: async (productId: string, data: { name: string; values: string[] }): Promise<ProductOption> => {
    const response = await apiClient.post(`/seller/products/${productId}/options`, data);
    return response.data.data ?? response.data;
  },

  removeOption: async (productId: string, optionId: string): Promise<void> => {
    await apiClient.delete(`/seller/products/${productId}/options/${optionId}`);
  },

  getOptions: async (productId: string): Promise<ProductOption[]> => {
    const response = await apiClient.get(`/seller/products/${productId}/options`);
    return response.data.options || response.data.data || [];
  },

  getVariants: async (productId: string): Promise<Variant[]> => {
    const response = await apiClient.get(`/seller/products/${productId}/variants`);
    return response.data.variants || response.data.data || [];
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
