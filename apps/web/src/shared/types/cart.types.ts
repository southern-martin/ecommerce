export interface CartItem {
  id: string;
  product_id: string;
  product_name: string;
  variant_id?: string;
  variant_options?: Record<string, string>;
  quantity: number;
  price_cents: number;
  image_url?: string;
  seller_id?: string;
}
