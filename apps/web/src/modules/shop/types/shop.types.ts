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
  name: string;
  value: string;
  price_modifier: number;
  stock_quantity: number;
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
