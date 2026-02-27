import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Product } from '@/modules/shop/types/shop.types';

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
    in_stock: true,
    stock_quantity: 100,
    seller: { id: raw.seller_id || '', name: '' },
    created_at: raw.created_at || '',
  };
}

export interface AdminProductFilter {
  page?: number;
  page_size?: number;
  seller_id?: string;
  category_id?: string;
  search?: string;
  sort?: string;
}

export interface CreateProductPayload {
  name: string;
  description: string;
  price: number;
  compare_at_price?: number;
  category_id: string;
  stock_quantity: number;
  images: { url: string; alt: string; is_primary: boolean }[];
}

export const adminProductMgmtApi = {
  listProducts: async (filters: AdminProductFilter = {}): Promise<PaginatedResponse<Product>> => {
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

  createProduct: async (data: CreateProductPayload): Promise<Product> => {
    const response = await apiClient.post('/seller/products', data);
    return response.data.data;
  },

  updateProduct: async (id: string, data: Partial<CreateProductPayload>): Promise<Product> => {
    const response = await apiClient.patch(`/seller/products/${id}`, data);
    return response.data.data;
  },

  deleteProduct: async (id: string): Promise<void> => {
    await apiClient.delete(`/seller/products/${id}`);
  },
};
