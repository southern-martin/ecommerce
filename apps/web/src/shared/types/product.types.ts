export interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price_cents: number;
  compare_at_price_cents?: number;
  category_id: string;
  category_name?: string;
  seller_id: string;
  seller_name?: string;
  images: string[];
  attributes: Record<string, string>;
  variants?: ProductVariant[];
  rating_avg: number;
  rating_count: number;
  stock: number;
  status: 'active' | 'draft' | 'archived';
  created_at: string;
}

export interface ProductVariant {
  id: string;
  sku: string;
  price_cents: number;
  stock: number;
  options: Record<string, string>;
  images: string[];
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  parent_id?: string;
  children?: Category[];
  image_url?: string;
}
