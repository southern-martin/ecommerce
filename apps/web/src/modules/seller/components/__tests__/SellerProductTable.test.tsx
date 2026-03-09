import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@/test/test-utils';
import { SellerProductTable } from '../SellerProductTable';
import type { SellerProduct } from '../../services/seller-product.api';

vi.mock('@/shared/lib/utils', () => ({
  formatPrice: (cents: number) => `$${(cents / 100).toFixed(2)}`,
  cn: (...args: string[]) => args.filter(Boolean).join(' '),
}));

const makeProduct = (overrides: Partial<SellerProduct> = {}): SellerProduct => ({
  id: 'p1',
  name: 'Test Product',
  slug: 'test-product',
  description: 'A test product',
  seller_id: 's1',
  category_id: 'c1',
  attribute_group_id: 'g1',
  base_price_cents: 2500,
  currency: 'USD',
  status: 'active',
  product_type: 'simple',
  has_variants: false,
  stock_quantity: 10,
  image_urls: [],
  tags: [],
  options: [],
  variants: [],
  attributes: [],
  created_at: '2026-01-01',
  updated_at: '2026-01-01',
  ...overrides,
});

describe('SellerProductTable', () => {
  it('should render product name and price', () => {
    render(<SellerProductTable products={[makeProduct()]} />);

    expect(screen.getByText('Test Product')).toBeInTheDocument();
    expect(screen.getByText('$25.00')).toBeInTheDocument();
  });

  it('should show stock quantity for simple products', () => {
    render(<SellerProductTable products={[makeProduct({ stock_quantity: 42 })]} />);

    expect(screen.getByText('42 in stock')).toBeInTheDocument();
  });

  it('should show variant count for configurable products', () => {
    render(
      <SellerProductTable
        products={[
          makeProduct({
            product_type: 'configurable',
            variants: [
              { id: 'v1', product_id: 'p1', sku: 'SKU1', name: 'V1', price_cents: 100, compare_at_cents: 0, cost_cents: 0, stock: 5, is_default: true, is_active: true, weight_grams: 100, barcode: '', image_urls: [], option_values: [], created_at: '', updated_at: '' },
              { id: 'v2', product_id: 'p1', sku: 'SKU2', name: 'V2', price_cents: 200, compare_at_cents: 0, cost_cents: 0, stock: 3, is_default: false, is_active: true, weight_grams: 100, barcode: '', image_urls: [], option_values: [], created_at: '', updated_at: '' },
            ],
          }),
        ]}
      />,
    );

    expect(screen.getByText('2 variants')).toBeInTheDocument();
  });

  it('should call onDelete when delete button is clicked', () => {
    const onDelete = vi.fn();
    render(<SellerProductTable products={[makeProduct()]} onDelete={onDelete} />);

    const deleteButtons = screen.getAllByRole('button');
    const trashButton = deleteButtons.find((btn) => btn.classList.contains('text-destructive'));
    expect(trashButton).toBeDefined();
    fireEvent.click(trashButton!);

    expect(onDelete).toHaveBeenCalledWith('p1');
  });

  it('should render status badge', () => {
    render(<SellerProductTable products={[makeProduct({ status: 'draft' })]} />);

    expect(screen.getByText('draft')).toBeInTheDocument();
  });
});
