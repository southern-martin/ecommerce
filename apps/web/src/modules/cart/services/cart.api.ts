import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface CartItem {
  id: string;
  product_id: string;
  name: string;
  slug: string;
  image_url: string;
  price: number;
  quantity: number;
  variant_id?: string;
  variant_name?: string;
}

export interface Cart {
  id: string;
  items: CartItem[];
  subtotal: number;
  item_count: number;
}

export const cartApi = {
  getCart: async (): Promise<Cart> => {
    const response = await apiClient.get<ApiResponse<Cart>>('/cart');
    return response.data.data;
  },

  addToCart: async (productId: string, quantity: number, variantId?: string): Promise<Cart> => {
    const response = await apiClient.post<ApiResponse<Cart>>('/cart/items', {
      product_id: productId,
      quantity,
      variant_id: variantId,
    });
    return response.data.data;
  },

  updateQuantity: async (itemId: string, quantity: number): Promise<Cart> => {
    const response = await apiClient.patch<ApiResponse<Cart>>('/cart/items', {
      item_id: itemId,
      quantity,
    });
    return response.data.data;
  },

  removeFromCart: async (itemId: string): Promise<Cart> => {
    const response = await apiClient.delete<ApiResponse<Cart>>('/cart/items', {
      data: { item_id: itemId },
    });
    return response.data.data;
  },
};
