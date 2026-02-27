import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Product, ProductVariant, ProductOption, ProductAttribute, FilterState } from '../types/shop.types';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function mapBackendProduct(raw: any): Product {
  // Map variants from backend shape
  const variants: ProductVariant[] = (raw.variants || []).map((v: any) => ({
    id: v.id,
    sku: v.sku || '',
    name: v.name || '',
    price_cents: v.price_cents || 0,
    compare_at_cents: v.compare_at_cents,
    cost_cents: v.cost_cents,
    stock: v.stock || 0,
    is_default: v.is_default || false,
    is_active: v.is_active !== false,
    weight_grams: v.weight_grams,
    barcode: v.barcode,
    image_urls: v.image_urls,
    option_values: (v.option_values || []).map((ov: any) => ({
      variant_id: ov.variant_id,
      option_id: ov.option_id,
      option_value_id: ov.option_value_id,
      option_name: ov.option_name,
      value: ov.value,
    })),
  }));

  // Map options from backend shape
  const options: ProductOption[] = (raw.options || []).map((o: any) => ({
    id: o.id,
    product_id: o.product_id,
    name: o.name,
    sort_order: o.sort_order || 0,
    values: (o.values || []).map((v: any) => ({
      id: v.id,
      option_id: v.option_id,
      value: v.value,
      color_hex: v.color_hex,
      sort_order: v.sort_order || 0,
    })),
  }));

  // Map attributes
  const attributes: ProductAttribute[] = (raw.attributes || []).map((a: any) => ({
    id: a.id,
    product_id: a.product_id,
    attribute_id: a.attribute_id,
    attribute_name: a.attribute_name,
    value: a.value,
    values: a.values,
  }));

  // Compute stock from variants if available
  const totalStock = variants.length > 0
    ? variants.filter(v => v.is_active).reduce((sum, v) => sum + v.stock, 0)
    : 100; // default if no variants

  return {
    id: raw.id,
    name: raw.name,
    slug: raw.slug,
    description: raw.description || '',
    price: raw.base_price_cents || 0,
    compare_at_price: undefined,
    images: (raw.image_urls || []).map((url: string, i: number) => ({
      id: `img-${i}`,
      url,
      alt: raw.name,
      is_primary: i === 0,
    })),
    category: { id: raw.category_id || '', name: '', slug: '' },
    rating: raw.rating_avg || 0,
    review_count: raw.rating_count || 0,
    in_stock: totalStock > 0,
    stock_quantity: totalStock,
    variants: variants.length > 0 ? variants : undefined,
    options: options.length > 0 ? options : undefined,
    attributes: attributes.length > 0 ? attributes : undefined,
    seller: { id: raw.seller_id || '', name: '' },
    created_at: raw.created_at || '',
  };
}

export const productApi = {
  getProducts: async (filters: Partial<FilterState>): Promise<PaginatedResponse<Product>> => {
    const response = await apiClient.get('/products', { params: filters });
    const raw = response.data;
    const products = (raw.products || []).map(mapBackendProduct);
    return {
      data: products,
      total: raw.total || 0,
      page: raw.page || 1,
      page_size: raw.pageSize || 20,
    };
  },

  getProductBySlug: async (slug: string): Promise<Product> => {
    const response = await apiClient.get(`/products/slug/${slug}`);
    const raw = response.data.data || response.data;
    return mapBackendProduct(raw);
  },

  getFeaturedProducts: async (): Promise<Product[]> => {
    const response = await apiClient.get('/products', { params: { page_size: 8, sort: 'newest' } });
    return (response.data.products || []).map(mapBackendProduct);
  },

  getTrendingProducts: async (): Promise<Product[]> => {
    const response = await apiClient.get('/products', { params: { page_size: 8, sort: 'newest' } });
    return (response.data.products || []).map(mapBackendProduct);
  },
};
