import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { CheckoutSummary } from '../CheckoutSummary';
import type { CartItem } from '@/modules/cart/services/cart.api';

const mockItems: CartItem[] = [
  {
    id: 'ci-1',
    product_id: 'p1',
    name: 'Premium Widget',
    slug: 'premium-widget',
    image_url: 'https://example.com/widget.jpg',
    price: 2500,
    quantity: 2,
  },
  {
    id: 'ci-2',
    product_id: 'p2',
    name: 'Basic Gadget',
    slug: 'basic-gadget',
    image_url: 'https://example.com/gadget.jpg',
    price: 1500,
    quantity: 1,
  },
];

describe('CheckoutSummary', () => {
  it('renders the Order Summary heading', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    expect(screen.getByText('Order Summary')).toBeInTheDocument();
  });

  it('displays each item with name and quantity', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    expect(screen.getByText('Premium Widget')).toBeInTheDocument();
    expect(screen.getByText('Basic Gadget')).toBeInTheDocument();
    expect(screen.getByText('Qty: 2')).toBeInTheDocument();
    expect(screen.getByText('Qty: 1')).toBeInTheDocument();
  });

  it('displays item line totals (price * quantity)', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    // Premium Widget: 2500 * 2 = 5000 cents = $50.00
    expect(screen.getByText('$50.00')).toBeInTheDocument();
    // Basic Gadget: 1500 * 1 = 1500 cents = $15.00
    expect(screen.getByText('$15.00')).toBeInTheDocument();
  });

  it('shows subtotal, shipping, and tax labels', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    expect(screen.getByText('Subtotal')).toBeInTheDocument();
    expect(screen.getByText('Shipping')).toBeInTheDocument();
    expect(screen.getByText('Tax')).toBeInTheDocument();
  });

  it('shows "Free" when shipping is 0', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={0} tax={585} discount={0} total={7085} />
    );
    expect(screen.getByText('Free')).toBeInTheDocument();
  });

  it('shows discount when greater than 0', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={1000} total={6585} />
    );
    expect(screen.getByText('Discount')).toBeInTheDocument();
    expect(screen.getByText('-$10.00')).toBeInTheDocument();
  });

  it('does not show discount row when discount is 0', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    expect(screen.queryByText('Discount')).not.toBeInTheDocument();
  });

  it('displays the total', () => {
    render(
      <CheckoutSummary items={mockItems} subtotal={6500} shipping={500} tax={585} discount={0} total={7585} />
    );
    expect(screen.getByText('Total')).toBeInTheDocument();
    expect(screen.getByText('$75.85')).toBeInTheDocument();
  });
});
