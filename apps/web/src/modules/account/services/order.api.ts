import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';
import type { Order } from '@/modules/checkout/services/order.api';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function mapOrder(raw: any): Order {
  return {
    id: raw.id,
    order_number: raw.order_number || raw.id,
    status: raw.status,
    items: (raw.items || []).map((item: any) => ({ // eslint-disable-line @typescript-eslint/no-explicit-any
      id: item.id,
      product_id: item.product_id,
      name: item.product_name || item.name,
      image_url: item.image_url || '',
      price: item.unit_price_cents || item.price || 0,
      quantity: item.quantity,
      variant_name: item.variant_name,
    })),
    subtotal: raw.subtotal_cents ?? raw.subtotal ?? 0,
    shipping_cost: raw.shipping_cents ?? raw.shipping_cost ?? 0,
    tax: raw.tax_cents ?? raw.tax ?? 0,
    discount: raw.discount_cents ?? raw.discount ?? 0,
    total: raw.total_cents ?? raw.total ?? 0,
    shipping_address: raw.shipping_address || {},
    created_at: raw.created_at || '',
  };
}

export const accountOrderApi = {
  getOrders: async (params: { page: number; page_size: number }): Promise<PaginatedResponse<Order>> => {
    const response = await apiClient.get('/orders', { params });
    const raw = response.data;
    // Backend may return { data: [...], total, page, page_size } or { orders: [...] }
    const orders = raw.data || raw.orders || [];
    return {
      data: orders.map(mapOrder),
      total: raw.total || orders.length,
      page: raw.page || params.page,
      page_size: raw.page_size || params.page_size,
    };
  },

  getOrderById: async (id: string): Promise<Order> => {
    const response = await apiClient.get(`/orders/${id}`);
    const raw = response.data.data || response.data;
    return mapOrder(raw);
  },
};
