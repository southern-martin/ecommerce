import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { DashboardStats } from '../DashboardStats';

vi.mock('lucide-react', () => ({
  DollarSign: (props: Record<string, unknown>) => <svg data-testid="icon-dollar" {...props} />,
  Users: (props: Record<string, unknown>) => <svg data-testid="icon-users" {...props} />,
  ShoppingCart: (props: Record<string, unknown>) => <svg data-testid="icon-cart" {...props} />,
  Package: (props: Record<string, unknown>) => <svg data-testid="icon-package" {...props} />,
}));

describe('DashboardStats', () => {
  const mockStats = {
    total_users: 1500,
    total_orders: 320,
    total_revenue: 9999900, // in cents = $99,999.00
    total_products: 85,
    new_users_today: 12,
    orders_today: 8,
  };

  it('should render all four stat cards', () => {
    render(<DashboardStats stats={mockStats} />);

    expect(screen.getByText('Total Revenue')).toBeDefined();
    expect(screen.getByText('Total Users')).toBeDefined();
    expect(screen.getByText('Total Orders')).toBeDefined();
    expect(screen.getByText('Total Products')).toBeDefined();
  });

  it('should display formatted revenue value', () => {
    render(<DashboardStats stats={mockStats} />);

    expect(screen.getByText('$99,999.00')).toBeDefined();
  });

  it('should display formatted user count', () => {
    render(<DashboardStats stats={mockStats} />);

    expect(screen.getByText('1,500')).toBeDefined();
  });

  it('should display formatted order and product counts', () => {
    render(<DashboardStats stats={mockStats} />);

    expect(screen.getByText('320')).toBeDefined();
    expect(screen.getByText('85')).toBeDefined();
  });
});
