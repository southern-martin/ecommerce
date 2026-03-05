import { useMemo } from 'react';
import { useSellerOrders } from './useSellerOrders';
import type { Order } from '@/modules/checkout/services/order.api';

interface RevenueByDay {
  date: string;
  revenue: number;
}

interface StatusCount {
  name: string;
  value: number;
}

interface ProductRevenue {
  name: string;
  revenue: number;
}

export function useSellerAnalytics() {
  const { data: ordersData, isLoading } = useSellerOrders(1, 100);

  const analytics = useMemo(() => {
    if (!ordersData?.data) return null;
    const orders = ordersData.data;

    return {
      revenueByDay: deriveRevenueByDay(orders),
      statusDistribution: deriveStatusDistribution(orders),
      topProducts: deriveTopProducts(orders),
    };
  }, [ordersData]);

  return { analytics, isLoading };
}

function deriveRevenueByDay(orders: Order[]): RevenueByDay[] {
  const now = new Date();
  const days = new Map<string, number>();

  for (let i = 6; i >= 0; i--) {
    const d = new Date(now);
    d.setDate(d.getDate() - i);
    const key = d.toLocaleDateString('en-US', { weekday: 'short' });
    days.set(key, 0);
  }

  orders.forEach((order) => {
    const orderDate = new Date(order.created_at);
    const diffMs = now.getTime() - orderDate.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    if (diffDays < 7 && diffDays >= 0) {
      const key = orderDate.toLocaleDateString('en-US', { weekday: 'short' });
      days.set(key, (days.get(key) || 0) + (order.total ?? 0));
    }
  });

  return Array.from(days.entries()).map(([date, revenue]) => ({ date, revenue }));
}

function deriveStatusDistribution(orders: Order[]): StatusCount[] {
  const counts = new Map<string, number>();
  orders.forEach((order) => {
    const status = order.status;
    counts.set(status, (counts.get(status) || 0) + 1);
  });
  return Array.from(counts.entries()).map(([name, value]) => ({
    name: name.charAt(0).toUpperCase() + name.slice(1),
    value,
  }));
}

function deriveTopProducts(orders: Order[]): ProductRevenue[] {
  const productRevenue = new Map<string, number>();
  orders.forEach((order) => {
    (order.items || []).forEach((item) => {
      const name = item.name;
      productRevenue.set(
        name,
        (productRevenue.get(name) || 0) + item.price * item.quantity
      );
    });
  });
  return Array.from(productRevenue.entries())
    .map(([name, revenue]) => ({ name, revenue }))
    .sort((a, b) => b.revenue - a.revenue)
    .slice(0, 5);
}
