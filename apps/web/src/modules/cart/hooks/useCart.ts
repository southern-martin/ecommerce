import { useCartStore } from '@/shared/stores/cart.store';

/**
 * useCart hook — reads from the Zustand cart store (client-side).
 * Provides a cart object shape compatible with CartPage and CheckoutPage.
 */
export function useCart() {
  const items = useCartStore((s) => s.items);
  const addItem = useCartStore((s) => s.addItem);
  const removeItem = useCartStore((s) => s.removeItem);
  const updateQty = useCartStore((s) => s.updateQuantity);
  const subtotalFn = useCartStore((s) => s.subtotal);
  const itemCountFn = useCartStore((s) => s.itemCount);

  const subtotal = subtotalFn();
  const itemCount = itemCountFn();

  // Map Zustand CartItem → shape expected by CartPage components
  const cart = {
    id: 'local',
    items: items.map((i) => ({
      id: i.id,
      product_id: i.product_id,
      name: i.product_name,
      slug: i.product_id, // use product_id as fallback slug
      image_url: i.image_url || '',
      price: i.price_cents,
      quantity: i.quantity,
      variant_id: i.variant_id,
      variant_name: i.variant_options
        ? Object.values(i.variant_options).join(', ')
        : undefined,
    })),
    subtotal,
    item_count: itemCount,
  };

  return {
    cart: itemCount > 0 ? cart : null,
    isLoading: false,
    addToCart: {
      mutate: ({
        productId,
        quantity,
        variantId,
      }: {
        productId: string;
        quantity: number;
        variantId?: string;
      }) => {
        addItem({
          id: `${productId}-${variantId || 'default'}`,
          product_id: productId,
          product_name: '',
          price_cents: 0,
          quantity,
          variant_id: variantId,
        });
      },
    },
    updateQuantity: {
      mutate: ({ itemId, quantity }: { itemId: string; quantity: number }) => {
        updateQty(itemId, quantity);
      },
    },
    removeFromCart: {
      mutate: (itemId: string) => {
        removeItem(itemId);
      },
    },
    itemCount,
  };
}
