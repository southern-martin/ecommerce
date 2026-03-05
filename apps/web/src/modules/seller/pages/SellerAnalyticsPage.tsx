import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { RevenueChart } from '../components/RevenueChart';
import { OrderStatusChart } from '../components/OrderStatusChart';
import { TopProductsChart } from '../components/TopProductsChart';
import { useSellerAnalytics } from '../hooks/useSellerAnalytics';

export default function SellerAnalyticsPage() {
  const { analytics, isLoading } = useSellerAnalytics();

  if (isLoading) {
    return (
      <div className="space-y-8">
        <h1 className="text-2xl font-bold">Analytics</h1>
        <Skeleton className="h-[350px]" />
        <div className="grid gap-6 md:grid-cols-2">
          <Skeleton className="h-[350px]" />
          <Skeleton className="h-[350px]" />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Analytics</h1>

      <Card>
        <CardHeader>
          <CardTitle>Revenue (Last 7 Days)</CardTitle>
        </CardHeader>
        <CardContent>
          <RevenueChart data={analytics?.revenueByDay ?? []} />
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Order Status Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <OrderStatusChart data={analytics?.statusDistribution ?? []} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top Products by Revenue</CardTitle>
          </CardHeader>
          <CardContent>
            <TopProductsChart data={analytics?.topProducts ?? []} />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
