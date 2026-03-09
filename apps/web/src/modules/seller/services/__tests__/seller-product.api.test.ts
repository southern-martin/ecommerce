import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockApiClient = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
}));

vi.mock('@/shared/lib/api-client', () => ({
  default: mockApiClient,
}));

import { sellerProductApi } from '../seller-product.api';

describe('sellerProductApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getProducts', () => {
    it('should fetch paginated products and map them', async () => {
      const rawProducts = [
        { id: 'p1', name: 'Widget', slug: 'widget', seller_id: 's1', base_price_cents: 1999, status: 'active' },
        { id: 'p2', name: 'Gadget', slug: 'gadget', seller_id: 's1' },
      ];
      mockApiClient.get.mockResolvedValue({ data: { products: rawProducts, total: 2, page: 1, page_size: 10 } });

      const result = await sellerProductApi.getProducts({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/products', { params: { page: 1, page_size: 10 } });
      expect(result.data).toHaveLength(2);
      expect(result.data[0].name).toBe('Widget');
      expect(result.data[0].currency).toBe('USD');
      expect(result.total).toBe(2);
    });

    it('should apply defaults for missing fields via mapRawProduct', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: [{ id: 'p1', name: 'Bare', slug: 'bare', seller_id: 's1' }], total: 1, page: 1 } });

      const result = await sellerProductApi.getProducts({ page: 1, page_size: 10 });
      const product = result.data[0];

      expect(product.description).toBe('');
      expect(product.currency).toBe('USD');
      expect(product.status).toBe('draft');
      expect(product.product_type).toBe('simple');
      expect(product.has_variants).toBe(false);
      expect(product.stock_quantity).toBe(0);
      expect(product.image_urls).toEqual([]);
      expect(product.tags).toEqual([]);
    });
  });

  describe('getProductById', () => {
    it('should fetch a single product by id', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: { id: 'p1', name: 'Widget', slug: 'widget', seller_id: 's1' } } });

      const result = await sellerProductApi.getProductById('p1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/products/p1');
      expect(result.id).toBe('p1');
      expect(result.name).toBe('Widget');
    });
  });

  describe('createProduct', () => {
    it('should post new product data and return mapped product', async () => {
      const createData = { name: 'New Product', description: 'desc', base_price_cents: 500, category_id: 'c1' };
      mockApiClient.post.mockResolvedValue({ data: { data: { id: 'p3', ...createData, slug: 'new-product', seller_id: 's1' } } });

      const result = await sellerProductApi.createProduct(createData);

      expect(mockApiClient.post).toHaveBeenCalledWith('/seller/products', createData);
      expect(result.id).toBe('p3');
      expect(result.name).toBe('New Product');
    });
  });

  describe('updateProduct', () => {
    it('should patch product and return mapped result', async () => {
      const updateData = { name: 'Updated Widget' };
      mockApiClient.patch.mockResolvedValue({ data: { data: { id: 'p1', name: 'Updated Widget', slug: 'widget', seller_id: 's1' } } });

      const result = await sellerProductApi.updateProduct('p1', updateData);

      expect(mockApiClient.patch).toHaveBeenCalledWith('/seller/products/p1', updateData);
      expect(result.name).toBe('Updated Widget');
    });
  });

  describe('deleteProduct', () => {
    it('should delete a product by id', async () => {
      mockApiClient.delete.mockResolvedValue({});

      await sellerProductApi.deleteProduct('p1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/seller/products/p1');
    });
  });

  describe('setProductAttributes', () => {
    it('should put attributes for a product', async () => {
      const attributes = [{ attribute_id: 'a1', value: 'Red' }];
      mockApiClient.put.mockResolvedValue({});

      await sellerProductApi.setProductAttributes('p1', attributes);

      expect(mockApiClient.put).toHaveBeenCalledWith('/seller/products/p1/attributes', { attributes });
    });
  });

  describe('getProductAttributes', () => {
    it('should fetch product attributes', async () => {
      const attrs = [{ id: 'pa1', product_id: 'p1', attribute_id: 'a1', attribute_name: 'Color', value: 'Red' }];
      mockApiClient.get.mockResolvedValue({ data: { attributes: attrs } });

      const result = await sellerProductApi.getProductAttributes('p1');

      expect(mockApiClient.get).toHaveBeenCalledWith('/seller/products/p1/attributes');
      expect(result).toEqual(attrs);
    });
  });

  describe('getAttributeGroups', () => {
    it('should fetch attribute groups', async () => {
      const groups = [{ id: 'g1', name: 'Clothing', slug: 'clothing' }];
      mockApiClient.get.mockResolvedValue({ data: { attribute_groups: groups } });

      const result = await sellerProductApi.getAttributeGroups();

      expect(mockApiClient.get).toHaveBeenCalledWith('/attribute-groups');
      expect(result).toEqual(groups);
    });
  });
});
