import { describe, it, expect, vi, beforeEach } from 'vitest';
import { cartApi } from '../cart.api';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
  post: ReturnType<typeof vi.fn>;
  patch: ReturnType<typeof vi.fn>;
  delete: ReturnType<typeof vi.fn>;
};

const mockCart = {
  id: 'cart-1',
  items: [
    {
      id: 'item-1',
      product_id: 'prod-1',
      name: 'Test Product',
      slug: 'test-product',
      image_url: '',
      price: 1999,
      quantity: 2,
    },
  ],
  subtotal: 3998,
  item_count: 2,
};

describe('cartApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getCart', () => {
    it('fetches cart from /cart', async () => {
      mockApiClient.get.mockResolvedValue({ data: { data: mockCart } });

      const result = await cartApi.getCart();

      expect(mockApiClient.get).toHaveBeenCalledWith('/cart');
      expect(result.id).toBe('cart-1');
      expect(result.items).toHaveLength(1);
    });
  });

  describe('addToCart', () => {
    it('sends POST with product data', async () => {
      mockApiClient.post.mockResolvedValue({ data: { data: mockCart } });

      const result = await cartApi.addToCart('prod-1', 2);

      expect(mockApiClient.post).toHaveBeenCalledWith('/cart/items', {
        product_id: 'prod-1',
        quantity: 2,
        variant_id: undefined,
      });
      expect(result.id).toBe('cart-1');
    });

    it('includes variant_id when provided', async () => {
      mockApiClient.post.mockResolvedValue({ data: { data: mockCart } });

      await cartApi.addToCart('prod-1', 1, 'variant-xl');

      expect(mockApiClient.post).toHaveBeenCalledWith('/cart/items', {
        product_id: 'prod-1',
        quantity: 1,
        variant_id: 'variant-xl',
      });
    });
  });

  describe('updateQuantity', () => {
    it('sends PATCH with item_id and quantity', async () => {
      mockApiClient.patch.mockResolvedValue({ data: { data: mockCart } });

      const result = await cartApi.updateQuantity('item-1', 5);

      expect(mockApiClient.patch).toHaveBeenCalledWith('/cart/items', {
        item_id: 'item-1',
        quantity: 5,
      });
      expect(result.id).toBe('cart-1');
    });
  });

  describe('removeFromCart', () => {
    it('sends DELETE with item_id in body', async () => {
      const emptyCart = { ...mockCart, items: [], subtotal: 0, item_count: 0 };
      mockApiClient.delete.mockResolvedValue({ data: { data: emptyCart } });

      const result = await cartApi.removeFromCart('item-1');

      expect(mockApiClient.delete).toHaveBeenCalledWith('/cart/items', {
        data: { item_id: 'item-1' },
      });
      expect(result.items).toHaveLength(0);
    });
  });
});
