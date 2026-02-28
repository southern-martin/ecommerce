import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';

export type ProductType = 'simple' | 'configurable' | string;

export interface SellerProduct {
  id: string;
  name: string;
  slug: string;
  description: string;
  seller_id: string;
  category_id: string;
  attribute_group_id: string;
  base_price_cents: number;
  currency: string;
  status: string;
  product_type: ProductType;
  has_variants: boolean;
  stock_quantity: number;
  image_urls: string[];
  tags: string[];
  options?: SellerProductOption[];
  variants?: SellerVariant[];
  attributes?: SellerProductAttribute[];
  created_at: string;
  updated_at: string;
}

export interface SellerProductOption {
  id: string;
  product_id: string;
  name: string;
  sort_order: number;
  values: SellerProductOptionValue[];
}

export interface SellerProductOptionValue {
  id: string;
  option_id: string;
  value: string;
  color_hex?: string;
  sort_order: number;
}

export interface SellerVariant {
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
  option_values: { variant_id: string; option_id: string; option_value_id: string; option_name: string; value: string }[];
  created_at: string;
  updated_at: string;
}

export interface SellerProductAttribute {
  id: string;
  product_id: string;
  attribute_id: string;
  attribute_name: string;
  value: string;
  values?: string[];
}

export interface AttributeGroupSummary {
  id: string;
  name: string;
  slug: string;
  description?: string;
  attribute_count?: number;
}

export interface CreateProductData {
  name: string;
  description: string;
  base_price_cents: number;
  currency?: string;
  category_id: string;
  attribute_group_id?: string;
  product_type?: ProductType;
  stock_quantity?: number;
  tags?: string[];
  image_urls?: string[];
  attributes?: { attribute_id: string; value: string; values?: string[] }[];
}

export type UpdateProductData = Partial<CreateProductData> & {
  status?: string;
};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function mapRawProduct(raw: any): SellerProduct {
  return {
    id: raw.id,
    name: raw.name,
    slug: raw.slug,
    description: raw.description || '',
    seller_id: raw.seller_id,
    category_id: raw.category_id || '',
    attribute_group_id: raw.attribute_group_id || '',
    base_price_cents: raw.base_price_cents || 0,
    currency: raw.currency || 'USD',
    status: raw.status || 'draft',
    product_type: raw.product_type || 'simple',
    has_variants: raw.has_variants || false,
    stock_quantity: raw.stock_quantity ?? 0,
    image_urls: raw.image_urls || [],
    tags: raw.tags || [],
    options: raw.options || [],
    variants: raw.variants || [],
    attributes: raw.attributes || [],
    created_at: raw.created_at || '',
    updated_at: raw.updated_at || '',
  };
}

export const sellerProductApi = {
  getProducts: async (params: { page: number; page_size: number }): Promise<PaginatedResponse<SellerProduct>> => {
    const response = await apiClient.get('/seller/products', { params });
    const raw = response.data;
    const products = (raw.products || raw.data || []).map(mapRawProduct);
    return {
      data: products,
      total: raw.total || 0,
      page: raw.page || 1,
      page_size: raw.pageSize || raw.page_size || 20,
    };
  },

  getProductById: async (id: string): Promise<SellerProduct> => {
    const response = await apiClient.get(`/seller/products/${id}`);
    const raw = response.data.data || response.data;
    return mapRawProduct(raw);
  },

  createProduct: async (data: CreateProductData): Promise<SellerProduct> => {
    const response = await apiClient.post('/seller/products', data);
    const raw = response.data.data || response.data;
    return mapRawProduct(raw);
  },

  updateProduct: async (id: string, data: UpdateProductData): Promise<SellerProduct> => {
    const response = await apiClient.patch(`/seller/products/${id}`, data);
    const raw = response.data.data || response.data;
    return mapRawProduct(raw);
  },

  deleteProduct: async (id: string): Promise<void> => {
    await apiClient.delete(`/seller/products/${id}`);
  },

  setProductAttributes: async (
    productId: string,
    attributes: { attribute_id: string; value: string; values?: string[] }[]
  ): Promise<void> => {
    await apiClient.put(`/seller/products/${productId}/attributes`, { attributes });
  },

  getProductAttributes: async (productId: string): Promise<SellerProductAttribute[]> => {
    const response = await apiClient.get(`/seller/products/${productId}/attributes`);
    return response.data.attributes || [];
  },

  getAttributeGroups: async (): Promise<AttributeGroupSummary[]> => {
    const response = await apiClient.get('/attribute-groups');
    return response.data.attribute_groups || [];
  },
};
