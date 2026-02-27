import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Product } from '@/modules/shop/types/shop.types';

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
    product_type: raw.product_type || 'simple',
    in_stock: raw.status === 'active',
    stock_quantity: raw.stock_quantity ?? 0,
    seller: { id: raw.seller_id || '', name: '' },
    created_at: raw.created_at || '',
    // Extra admin fields stored in the Product object
    _status: raw.status || 'draft',
    _tags: raw.tags || [],
    _currency: raw.currency || 'USD',
    _has_variants: raw.has_variants || false,
    _product_type: raw.product_type || 'simple',
  } as Product & { _status: string; _tags: string[]; _currency: string; _has_variants: boolean; _product_type: string };
}

export interface AdminProductFilter {
  page?: number;
  page_size?: number;
  seller_id?: string;
  category_id?: string;
  status?: string;
  search?: string;
  sort?: string;
}

export interface AdminUpdateProductPayload {
  name?: string;
  description?: string;
  base_price_cents?: number;
  currency?: string;
  status?: string;
  tags?: string[];
  image_urls?: string[];
  category_id?: string;
}

export interface CreateProductPayload {
  name: string;
  description: string;
  base_price_cents: number;
  currency?: string;
  category_id: string;
  tags?: string[];
  image_urls?: string[];
}

export const adminProductMgmtApi = {
  // Uses dedicated admin endpoint — returns ALL products (all statuses, all sellers)
  listProducts: async (filters: AdminProductFilter = {}): Promise<PaginatedResponse<Product>> => {
    const params: Record<string, string | number> = {};
    if (filters.page) params.page = filters.page;
    if (filters.page_size) params.page_size = filters.page_size;
    if (filters.seller_id) params.seller_id = filters.seller_id;
    if (filters.category_id) params.category_id = filters.category_id;
    if (filters.status) params.status = filters.status;
    if (filters.search) params.q = filters.search;
    if (filters.sort) params.sort_by = filters.sort;

    const response = await apiClient.get('/admin/products', { params });
    const raw = response.data;
    const products = (raw.products || []).map(mapBackendProduct);
    return {
      data: products,
      total: raw.total || 0,
      page: raw.page || 1,
      page_size: raw.pageSize || 20,
    };
  },

  // Admin create uses seller endpoint (needs seller_id context)
  createProduct: async (data: CreateProductPayload): Promise<Product> => {
    const response = await apiClient.post('/seller/products', data);
    return mapBackendProduct(response.data);
  },

  // Admin update — uses dedicated admin endpoint (no seller ownership check)
  updateProduct: async (id: string, data: AdminUpdateProductPayload): Promise<Product> => {
    const response = await apiClient.patch(`/admin/products/${id}`, data);
    return mapBackendProduct(response.data);
  },

  // Admin delete — uses dedicated admin endpoint (no seller ownership check)
  deleteProduct: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/products/${id}`);
  },
};
