import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { RevenueChart } from '../components/RevenueChart';
import { OrderStatusChart } from '../components/OrderStatusChart';
import { TopProductsChart } from '../components/TopProductsChart';

export default function SellerAnalyticsPage() {
  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Analytics</h1>

      <Card>
        <CardHeader>
          <CardTitle>Revenue (Last 7 Days)</CardTitle>
        </CardHeader>
        <CardContent>
          <RevenueChart />
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Order Status Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <OrderStatusChart />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top Products by Revenue</CardTitle>
          </CardHeader>
          <CardContent>
            <TopProductsChart />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
