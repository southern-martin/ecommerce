import apiClient from '@/shared/lib/api-client';

export interface ShippingAddress {
  first_name: string;
  last_name: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  phone: string;
}

export interface CreateOrderData {
  shipping_address: ShippingAddress;
  payment_method_id: string;
  coupon_code?: string;
  notes?: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  name: string;
  image_url: string;
  price: number;
  quantity: number;
  variant_name?: string;
}

export interface Order {
  id: string;
  order_number: string;
  status: string;
  items: OrderItem[];
  subtotal: number;
  shipping_cost: number;
  tax: number;
  discount: number;
  total: number;
  shipping_address: ShippingAddress;
  created_at: string;
}

/**
 * Maps the frontend checkout data + cart items to the backend CreateOrder request shape.
 * Backend expects: buyer_id, currency, shipping_address{full_name,line1,...}, items[{product_id,product_name,quantity,unit_price_cents,seller_id,...}]
 */
export const orderApi = {
  createOrder: async (
    data: CreateOrderData,
    cartItems: Array<{
      product_id: string;
      product_name: string;
      quantity: number;
      price_cents: number;
      image_url?: string;
      variant_id?: string;
      seller_id?: string;
    }>,
    userId: string
  ): Promise<Order> => {
    const payload = {
      buyer_id: userId,
      currency: 'USD',
      shipping_address: {
        full_name: `${data.shipping_address.first_name} ${data.shipping_address.last_name}`,
        line1: data.shipping_address.address_line1,
        line2: data.shipping_address.address_line2 || '',
        city: data.shipping_address.city,
        state: data.shipping_address.state,
        postal_code: data.shipping_address.postal_code,
        country_code: data.shipping_address.country,
        phone: data.shipping_address.phone,
      },
      items: cartItems.map((item) => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const orderItem: any = {
          product_id: item.product_id,
          product_name: item.product_name,
          quantity: item.quantity,
          unit_price_cents: item.price_cents,
          seller_id: item.seller_id || userId,
          image_url: item.image_url || '',
        };
        // Only include variant_id if it's a valid non-empty string (backend expects UUID or omit)
        if (item.variant_id) {
          orderItem.variant_id = item.variant_id;
        }
        return orderItem;
      }),
    };

    const response = await apiClient.post('/orders', payload);
    const raw = response.data.data || response.data;
    return {
      id: raw.id,
      order_number: raw.order_number || raw.id,
      status: raw.status,
      items: raw.items || [],
      subtotal: raw.subtotal_cents || 0,
      shipping_cost: raw.shipping_cents || 0,
      tax: raw.tax_cents || 0,
      discount: raw.discount_cents || 0,
      total: raw.total_cents || 0,
      shipping_address: data.shipping_address,
      created_at: raw.created_at || new Date().toISOString(),
    };
  },
};
