export interface Order {
  id: string;
  buyer_id: string;
  seller_id: string;
  status: string;
  items: OrderItem[];
  subtotal_cents: number;
  shipping_cents: number;
  tax_cents: number;
  discount_cents: number;
  total_cents: number;
  shipping_address: Address;
  created_at: string;
  updated_at: string;
}

export interface OrderItem {
  id: string;
  product_id: string;
  product_name: string;
  variant_id?: string;
  quantity: number;
  price_cents: number;
  image_url?: string;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  zip_code: string;
  country: string;
}
