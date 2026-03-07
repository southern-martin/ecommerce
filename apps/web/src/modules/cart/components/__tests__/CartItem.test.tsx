import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@/test/test-utils';
import { CartItem } from '../CartItem';
import type { CartItem as CartItemType } from '../../services/cart.api';

const mockItem: CartItemType = {
  id: 'item-1',
  product_id: 'prod-1',
  name: 'Wireless Headphones',
  slug: 'wireless-headphones',
  image_url: 'https://example.com/headphones.jpg',
  price: 4999,
  quantity: 2,
};

describe('CartItem', () => {
  it('renders the product name', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    expect(screen.getByText('Wireless Headphones')).toBeInTheDocument();
  });

  it('renders the product image', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    const img = screen.getByAltText('Wireless Headphones');
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', 'https://example.com/headphones.jpg');
  });

  it('renders the quantity', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('renders the total price (price * quantity)', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    // 4999 * 2 = 9998 cents = $99.98
    expect(screen.getByText('$99.98')).toBeInTheDocument();
  });

  it('renders variant name when provided', () => {
    const itemWithVariant = { ...mockItem, variant_name: 'Black / Large' };
    render(
      <CartItem
        item={itemWithVariant}
        onUpdateQuantity={vi.fn()}
        onRemove={vi.fn()}
      />
    );
    expect(screen.getByText('Black / Large')).toBeInTheDocument();
  });

  it('does not render variant text when not provided', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    // No variant_name, so no extra paragraph
    expect(screen.queryByText('Black / Large')).not.toBeInTheDocument();
  });

  it('calls onUpdateQuantity with decremented value on minus click', () => {
    const onUpdate = vi.fn();
    render(
      <CartItem item={mockItem} onUpdateQuantity={onUpdate} onRemove={vi.fn()} />
    );

    // Find minus button (first icon button)
    const buttons = screen.getAllByRole('button');
    const minusBtn = buttons[0]; // First button is minus
    fireEvent.click(minusBtn);

    expect(onUpdate).toHaveBeenCalledWith('item-1', 1); // max(1, 2-1) = 1
  });

  it('calls onUpdateQuantity with incremented value on plus click', () => {
    const onUpdate = vi.fn();
    render(
      <CartItem item={mockItem} onUpdateQuantity={onUpdate} onRemove={vi.fn()} />
    );

    const buttons = screen.getAllByRole('button');
    const plusBtn = buttons[1]; // Second button is plus
    fireEvent.click(plusBtn);

    expect(onUpdate).toHaveBeenCalledWith('item-1', 3); // 2+1 = 3
  });

  it('does not go below quantity 1 on minus click', () => {
    const singleItem = { ...mockItem, quantity: 1 };
    const onUpdate = vi.fn();
    render(
      <CartItem
        item={singleItem}
        onUpdateQuantity={onUpdate}
        onRemove={vi.fn()}
      />
    );

    const buttons = screen.getAllByRole('button');
    fireEvent.click(buttons[0]); // minus

    expect(onUpdate).toHaveBeenCalledWith('item-1', 1); // max(1, 1-1) = 1
  });

  it('calls onRemove when delete button is clicked', () => {
    const onRemove = vi.fn();
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={onRemove} />
    );

    const buttons = screen.getAllByRole('button');
    const deleteBtn = buttons[2]; // Third button is delete
    fireEvent.click(deleteBtn);

    expect(onRemove).toHaveBeenCalledWith('item-1');
  });

  it('links to product detail page', () => {
    render(
      <CartItem item={mockItem} onUpdateQuantity={vi.fn()} onRemove={vi.fn()} />
    );
    const links = screen.getAllByRole('link');
    expect(links.some((l) => l.getAttribute('href') === '/products/wireless-headphones')).toBe(true);
  });
});
