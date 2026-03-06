import { describe, it, expect, beforeEach } from 'vitest';
import { useCartStore } from '../cart.store';
import type { CartItem } from '../../types/cart.types';

const mockItem: CartItem = {
  id: 'item-1',
  product_id: 'prod-1',
  product_name: 'Test Product',
  quantity: 2,
  price_cents: 1999,
  image_url: 'https://example.com/img.jpg',
};

const mockItem2: CartItem = {
  id: 'item-2',
  product_id: 'prod-2',
  product_name: 'Another Product',
  quantity: 1,
  price_cents: 500,
};

describe('useCartStore', () => {
  beforeEach(() => {
    useCartStore.setState({ items: [] });
  });

  describe('addItem', () => {
    it('adds a new item to empty cart', () => {
      useCartStore.getState().addItem(mockItem);
      expect(useCartStore.getState().items).toHaveLength(1);
      expect(useCartStore.getState().items[0]).toEqual(mockItem);
    });

    it('merges quantity for duplicate item', () => {
      useCartStore.getState().addItem(mockItem);
      useCartStore.getState().addItem({ ...mockItem, quantity: 3 });
      expect(useCartStore.getState().items).toHaveLength(1);
      expect(useCartStore.getState().items[0].quantity).toBe(5);
    });

    it('adds multiple different items', () => {
      useCartStore.getState().addItem(mockItem);
      useCartStore.getState().addItem(mockItem2);
      expect(useCartStore.getState().items).toHaveLength(2);
    });
  });

  describe('removeItem', () => {
    it('removes an existing item', () => {
      useCartStore.setState({ items: [mockItem, mockItem2] });
      useCartStore.getState().removeItem('item-1');
      expect(useCartStore.getState().items).toHaveLength(1);
      expect(useCartStore.getState().items[0].id).toBe('item-2');
    });

    it('does nothing for non-existing item', () => {
      useCartStore.setState({ items: [mockItem] });
      useCartStore.getState().removeItem('non-existent');
      expect(useCartStore.getState().items).toHaveLength(1);
    });
  });

  describe('updateQuantity', () => {
    it('updates quantity of an item', () => {
      useCartStore.setState({ items: [mockItem] });
      useCartStore.getState().updateQuantity('item-1', 5);
      expect(useCartStore.getState().items[0].quantity).toBe(5);
    });

    it('removes item when quantity is zero', () => {
      useCartStore.setState({ items: [mockItem] });
      useCartStore.getState().updateQuantity('item-1', 0);
      expect(useCartStore.getState().items).toHaveLength(0);
    });

    it('removes item when quantity is negative', () => {
      useCartStore.setState({ items: [mockItem] });
      useCartStore.getState().updateQuantity('item-1', -1);
      expect(useCartStore.getState().items).toHaveLength(0);
    });
  });

  describe('clearCart', () => {
    it('removes all items', () => {
      useCartStore.setState({ items: [mockItem, mockItem2] });
      useCartStore.getState().clearCart();
      expect(useCartStore.getState().items).toHaveLength(0);
    });
  });

  describe('subtotal', () => {
    it('calculates correct subtotal in cents', () => {
      useCartStore.setState({ items: [mockItem, mockItem2] });
      // mockItem: 1999 * 2 = 3998
      // mockItem2: 500 * 1 = 500
      expect(useCartStore.getState().subtotal()).toBe(4498);
    });

    it('returns 0 for empty cart', () => {
      expect(useCartStore.getState().subtotal()).toBe(0);
    });
  });

  describe('itemCount', () => {
    it('sums quantities across all items', () => {
      useCartStore.setState({ items: [mockItem, mockItem2] });
      // mockItem: qty 2 + mockItem2: qty 1 = 3
      expect(useCartStore.getState().itemCount()).toBe(3);
    });

    it('returns 0 for empty cart', () => {
      expect(useCartStore.getState().itemCount()).toBe(0);
    });
  });
});
