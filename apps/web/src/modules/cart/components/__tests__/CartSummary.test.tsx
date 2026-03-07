import { describe, it, expect } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { CartSummary } from '../CartSummary';

describe('CartSummary', () => {
  it('renders order summary heading', () => {
    render(<CartSummary subtotal={5000} itemCount={3} />);
    expect(screen.getByText('Order Summary')).toBeInTheDocument();
  });

  it('displays subtotal with item count', () => {
    render(<CartSummary subtotal={5000} itemCount={3} />);
    expect(screen.getByText('Subtotal (3 items)')).toBeInTheDocument();
    // subtotal and total may both show $50.00 when shipping is free
    expect(screen.getAllByText('$50.00').length).toBeGreaterThanOrEqual(1);
  });

  it('uses singular "item" for count of 1', () => {
    render(<CartSummary subtotal={1999} itemCount={1} />);
    expect(screen.getByText('Subtotal (1 item)')).toBeInTheDocument();
  });

  it('shows Free when estimated shipping is 0', () => {
    render(<CartSummary subtotal={5000} itemCount={2} estimatedShipping={0} />);
    expect(screen.getByText('Free')).toBeInTheDocument();
  });

  it('shows shipping cost when provided', () => {
    render(<CartSummary subtotal={5000} itemCount={2} estimatedShipping={999} />);
    expect(screen.getByText('$9.99')).toBeInTheDocument();
  });

  it('calculates total including shipping', () => {
    render(<CartSummary subtotal={5000} itemCount={2} estimatedShipping={999} />);
    // Total = 5000 + 999 = 5999 = $59.99
    expect(screen.getByText('$59.99')).toBeInTheDocument();
  });

  it('shows different subtotal and total when shipping is non-zero', () => {
    const { container } = render(
      <CartSummary subtotal={3000} itemCount={1} estimatedShipping={500} />
    );
    // Subtotal $30.00, Shipping $5.00, Total $35.00
    expect(screen.getByText('$30.00')).toBeInTheDocument();
    expect(screen.getByText('$5.00')).toBeInTheDocument();
    expect(screen.getByText('$35.00')).toBeInTheDocument();
  });

  it('shows checkout button by default', () => {
    render(<CartSummary subtotal={5000} itemCount={2} />);
    expect(screen.getByText('Proceed to Checkout')).toBeInTheDocument();
  });

  it('hides checkout button when showCheckoutButton is false', () => {
    render(
      <CartSummary subtotal={5000} itemCount={2} showCheckoutButton={false} />
    );
    expect(screen.queryByText('Proceed to Checkout')).not.toBeInTheDocument();
  });

  it('shows secure checkout text', () => {
    render(<CartSummary subtotal={5000} itemCount={2} />);
    expect(screen.getByText('Secure checkout')).toBeInTheDocument();
  });

  it('links checkout button to /checkout', () => {
    render(<CartSummary subtotal={5000} itemCount={2} />);
    const link = screen.getByText('Proceed to Checkout').closest('a');
    expect(link).toHaveAttribute('href', '/checkout');
  });
});
