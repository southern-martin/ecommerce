import { useQuery } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';
import { Card } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';

export default function PromotionsPage() {
  const { data: flashSales = [], isLoading: loadingSales } = useQuery({
    queryKey: ['flash-sales'],
    queryFn: async () => {
      const res = await apiClient.get('/promotions/flash-sales');
      const d = res.data;
      return Array.isArray(d) ? d : (d as any).data || [];
    },
  });

  const { data: coupons = [], isLoading: loadingCoupons } = useQuery({
    queryKey: ['coupons'],
    queryFn: async () => {
      const res = await apiClient.get('/promotions/coupons');
      const d = res.data;
      return Array.isArray(d) ? d : (d as any).data || [];
    },
  });

  if (loadingSales || loadingCoupons) return <div className="p-6">Loading promotions...</div>;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-3xl font-bold mb-6">Promotions & Deals</h1>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-4">Flash Sales</h2>
        {flashSales.length === 0 ? (
          <p className="text-muted-foreground">No active flash sales right now.</p>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {flashSales.map((sale: any) => (
              <Card key={sale.id} className="p-4">
                <h3 className="font-medium">{sale.name}</h3>
                <p className="text-sm text-muted-foreground mt-1">{sale.description}</p>
                <Badge className="mt-2" variant="outline">{sale.status || 'active'}</Badge>
              </Card>
            ))}
          </div>
        )}
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-4">Available Coupons</h2>
        {coupons.length === 0 ? (
          <p className="text-muted-foreground">No coupons available.</p>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {coupons.map((coupon: any) => (
              <Card key={coupon.id} className="p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-mono font-bold text-lg">{coupon.code}</p>
                    <p className="text-sm text-muted-foreground">{coupon.description}</p>
                  </div>
                  <Badge>{coupon.type}</Badge>
                </div>
              </Card>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
