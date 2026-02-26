import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Product, FilterState } from '../types/shop.types';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function mapBackendProduct(raw: any): Product {
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
    in_stock: true, // Default to true since backend doesn't have stock on product level
    stock_quantity: 100,
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
