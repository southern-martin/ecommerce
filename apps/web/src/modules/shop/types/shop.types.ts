export interface Product {
  id: string;
  name: string;
  slug: string;
  description: string;
  price: number;
  compare_at_price?: number;
  images: ProductImage[];
  category: Category;
  rating: number;
  review_count: number;
  in_stock: boolean;
  stock_quantity: number;
  variants?: ProductVariant[];
  options?: ProductOption[];
  attributes?: ProductAttribute[];
  seller: { id: string; name: string };
  created_at: string;
}

export interface ProductImage {
  id: string;
  url: string;
  alt: string;
  is_primary: boolean;
}

export interface ProductVariant {
  id: string;
  sku: string;
  name: string;
  price_cents: number;
  compare_at_cents?: number;
  cost_cents?: number;
  stock: number;
  is_default: boolean;
  is_active: boolean;
  weight_grams?: number;
  barcode?: string;
  image_urls?: string[];
  option_values: VariantOptionValue[];
}

export interface VariantOptionValue {
  variant_id: string;
  option_id: string;
  option_value_id: string;
  option_name: string;
  value: string;
}

export interface ProductOption {
  id: string;
  product_id: string;
  name: string;
  sort_order: number;
  values: ProductOptionValue[];
}

export interface ProductOptionValue {
  id: string;
  option_id: string;
  value: string;
  color_hex?: string;
  sort_order: number;
}

export interface ProductAttribute {
  id: string;
  product_id: string;
  attribute_id: string;
  attribute_name: string;
  value: string;
  values?: string[];
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  parent_id?: string;
  image_url?: string;
  children?: Category[];
}

export interface FilterState {
  category?: string;
  min_price?: number;
  max_price?: number;
  rating?: number;
  in_stock?: boolean;
  page: number;
  page_size: number;
  sort: SortOption;
  search?: string;
}

export type SortOption =
  | 'newest'
  | 'price_asc'
  | 'price_desc'
  | 'rating'
  | 'popular';
