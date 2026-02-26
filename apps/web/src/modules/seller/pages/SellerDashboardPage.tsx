import { DollarSign, Package, ShoppingCart, AlertCircle } from 'lucide-react';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { StatsCard } from '../components/StatsCard';
import { SellerOrderTable } from '../components/SellerOrderTable';
import { useSellerDashboard } from '../hooks/useSellerOrders';
import { useSellerOrders } from '../hooks/useSellerOrders';
import { formatPrice } from '@/shared/lib/utils';

export default function SellerDashboardPage() {
  const { data: stats, isLoading: statsLoading } = useSellerDashboard();
  const { data: recentOrders, isLoading: ordersLoading } = useSellerOrders(1, 5);

  if (statsLoading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-28" />
          ))}
        </div>
        <Skeleton className="h-64" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      {stats && (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            icon={DollarSign}
            label="Total Revenue"
            value={formatPrice(stats.total_revenue)}
            trend={stats.revenue_trend}
          />
          <StatsCard
            icon={ShoppingCart}
            label="Total Orders"
            value={stats.total_orders.toLocaleString()}
            trend={stats.orders_trend}
          />
          <StatsCard
            icon={Package}
            label="Total Products"
            value={stats.total_products.toLocaleString()}
          />
          <StatsCard
            icon={AlertCircle}
            label="Pending Orders"
            value={stats.pending_orders.toLocaleString()}
          />
        </div>
      )}

      <div>
        <h2 className="mb-4 text-lg font-semibold">Recent Orders</h2>
        {ordersLoading ? (
          <Skeleton className="h-48" />
        ) : recentOrders && recentOrders.data.length > 0 ? (
          <SellerOrderTable orders={recentOrders.data} />
        ) : (
          <p className="py-8 text-center text-muted-foreground">No orders yet.</p>
        )}
      </div>
    </div>
  );
}
