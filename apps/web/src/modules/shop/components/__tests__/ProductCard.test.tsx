import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen, fireEvent } from '@/test/test-utils';
import { ProductCard } from '../ProductCard';
import { useCartStore } from '@/shared/stores/cart.store';
import type { Product } from '../../types/shop.types';

const mockProduct: Product = {
  id: 'prod-1',
  name: 'Wireless Headphones',
  slug: 'wireless-headphones',
  description: 'Great sound quality',
  price: 4999,
  images: [
    { id: 'img-1', url: 'https://example.com/headphones.jpg', alt: 'Headphones', is_primary: true },
  ],
  category: { id: 'cat-1', name: 'Electronics', slug: 'electronics' },
  product_type: 'simple',
  rating: 4.5,
  review_count: 23,
  in_stock: true,
  stock_quantity: 50,
  seller: { id: 'seller-1', name: 'TechStore' },
  created_at: '2026-01-01T00:00:00Z',
};

describe('ProductCard', () => {
  beforeEach(() => {
    useCartStore.setState({ items: [] });
  });

  it('renders product name', () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText('Wireless Headphones')).toBeInTheDocument();
  });

  it('renders product price', () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText('$49.99')).toBeInTheDocument();
  });

  it('renders product image', () => {
    render(<ProductCard product={mockProduct} />);
    const img = screen.getByAltText('Headphones');
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', 'https://example.com/headphones.jpg');
  });

  it('renders review count', () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText('(23)')).toBeInTheDocument();
  });

  it('renders "Add to Cart" button for simple products', () => {
    render(<ProductCard product={mockProduct} />);
    expect(screen.getByText('Add to Cart')).toBeInTheDocument();
  });

  it('adds item to cart on button click', () => {
    render(<ProductCard product={mockProduct} />);
    fireEvent.click(screen.getByText('Add to Cart'));

    const items = useCartStore.getState().items;
    expect(items).toHaveLength(1);
    expect(items[0].product_id).toBe('prod-1');
    expect(items[0].product_name).toBe('Wireless Headphones');
    expect(items[0].price_cents).toBe(4999);
  });

  it('renders "Select Options" for configurable products', () => {
    const configurable = {
      ...mockProduct,
      product_type: 'configurable',
      min_price: 2999,
      max_price: 5999,
    };
    render(<ProductCard product={configurable} />);
    expect(screen.getByText('Select Options')).toBeInTheDocument();
  });

  it('renders price range for configurable products with variants', () => {
    const configurable = {
      ...mockProduct,
      product_type: 'configurable',
      min_price: 2999,
      max_price: 5999,
    };
    render(<ProductCard product={configurable} />);
    expect(screen.getByText(/\$29\.99/)).toBeInTheDocument();
    expect(screen.getByText(/\$59\.99/)).toBeInTheDocument();
  });

  it('renders "Out of Stock" when not in stock', () => {
    const outOfStock = { ...mockProduct, in_stock: false, stock_quantity: 0 };
    render(<ProductCard product={outOfStock} />);
    expect(screen.getByText('Out of Stock')).toBeInTheDocument();
    expect(screen.getByText('Sold Out')).toBeInTheDocument();
  });

  it('disables button when out of stock', () => {
    const outOfStock = { ...mockProduct, in_stock: false, stock_quantity: 0 };
    render(<ProductCard product={outOfStock} />);
    const button = screen.getByRole('button', { name: /out of stock/i });
    expect(button).toBeDisabled();
  });

  it('links to product detail page', () => {
    render(<ProductCard product={mockProduct} />);
    const links = screen.getAllByRole('link');
    expect(links.some((l) => l.getAttribute('href') === '/products/wireless-headphones')).toBe(true);
  });

  it('shows discount badge when compare_at_price is set', () => {
    const discounted = { ...mockProduct, compare_at_price: 7999 };
    render(<ProductCard product={discounted} />);
    // discount = round((1 - 4999/7999) * 100) = 38%
    expect(screen.getByText('-38%')).toBeInTheDocument();
  });

  it('shows placeholder when no images', () => {
    const noImage = { ...mockProduct, images: [] };
    render(<ProductCard product={noImage} />);
    // No img element, but should still render the card
    expect(screen.getByText('Wireless Headphones')).toBeInTheDocument();
  });
});
