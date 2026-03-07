import { describe, it, expect, vi, beforeEach } from 'vitest';
import { productApi } from '../product.api';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
};

const mockRawProduct = {
  id: 'prod-1',
  name: 'Test Product',
  slug: 'test-product',
  description: 'A great product',
  base_price_cents: 2999,
  category_id: 'cat-1',
  product_type: 'simple',
  stock_quantity: 50,
  rating_avg: 4.5,
  rating_count: 12,
  image_urls: ['https://example.com/img1.jpg'],
  seller_id: 'seller-1',
  created_at: '2026-01-01T00:00:00Z',
};

describe('productApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getProducts', () => {
    it('fetches products with filters', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          products: [mockRawProduct],
          total: 1,
          page: 1,
          pageSize: 20,
        },
      });

      const result = await productApi.getProducts({ page: 1, page_size: 20, sort: 'newest' });

      expect(mockApiClient.get).toHaveBeenCalledWith('/products', {
        params: { page: 1, page_size: 20, sort: 'newest' },
      });
      expect(result.data).toHaveLength(1);
      expect(result.data[0].name).toBe('Test Product');
      expect(result.data[0].price).toBe(2999);
      expect(result.total).toBe(1);
    });

    it('maps backend product fields correctly', async () => {
      mockApiClient.get.mockResolvedValue({
        data: {
          products: [mockRawProduct],
          total: 1,
          page: 1,
          pageSize: 20,
        },
      });

      const result = await productApi.getProducts({});
      const product = result.data[0];

      expect(product.id).toBe('prod-1');
      expect(product.slug).toBe('test-product');
      expect(product.rating).toBe(4.5);
      expect(product.review_count).toBe(12);
      expect(product.in_stock).toBe(true);
      expect(product.stock_quantity).toBe(50);
      expect(product.images).toHaveLength(1);
      expect(product.images[0].url).toBe('https://example.com/img1.jpg');
      expect(product.images[0].is_primary).toBe(true);
    });

    it('handles empty product list', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { products: [], total: 0, page: 1, pageSize: 20 },
      });

      const result = await productApi.getProducts({});
      expect(result.data).toHaveLength(0);
      expect(result.total).toBe(0);
    });

    it('handles products with variants', async () => {
      const productWithVariants = {
        ...mockRawProduct,
        product_type: 'configurable',
        variants: [
          { id: 'v1', price_cents: 2999, stock: 10, is_active: true, option_values: [] },
          { id: 'v2', price_cents: 3999, stock: 5, is_active: true, option_values: [] },
        ],
      };
      mockApiClient.get.mockResolvedValue({
        data: { products: [productWithVariants], total: 1, page: 1, pageSize: 20 },
      });

      const result = await productApi.getProducts({});
      const product = result.data[0];

      expect(product.variants).toHaveLength(2);
      expect(product.min_price).toBe(2999);
      expect(product.max_price).toBe(3999);
      expect(product.stock_quantity).toBe(15); // 10 + 5
    });
  });

  describe('getProductBySlug', () => {
    it('fetches a single product by slug', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { data: mockRawProduct },
      });

      const result = await productApi.getProductBySlug('test-product');

      expect(mockApiClient.get).toHaveBeenCalledWith('/products/slug/test-product');
      expect(result.name).toBe('Test Product');
      expect(result.slug).toBe('test-product');
    });
  });

  describe('getFeaturedProducts', () => {
    it('fetches featured products with limit 8', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { products: [mockRawProduct] },
      });

      const result = await productApi.getFeaturedProducts();

      expect(mockApiClient.get).toHaveBeenCalledWith('/products', {
        params: { page_size: 8, sort: 'newest' },
      });
      expect(result).toHaveLength(1);
    });
  });

  describe('getTrendingProducts', () => {
    it('fetches trending products', async () => {
      mockApiClient.get.mockResolvedValue({
        data: { products: [mockRawProduct] },
      });

      const result = await productApi.getTrendingProducts();

      expect(mockApiClient.get).toHaveBeenCalledWith('/products', {
        params: { page_size: 8, sort: 'newest' },
      });
      expect(result).toHaveLength(1);
    });
  });
});
