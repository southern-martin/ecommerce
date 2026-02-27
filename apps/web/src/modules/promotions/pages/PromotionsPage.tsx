import { useQuery } from '@tanstack/react-query';
import { Zap, Ticket, Tag, Sparkles } from 'lucide-react';
import apiClient from '@/shared/lib/api-client';
import { Badge } from '@/shared/components/ui/badge';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { Skeleton } from '@/shared/components/ui/skeleton';

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

  if (loadingSales || loadingCoupons) {
    return (
      <PageLayout
        title="Deals & Promotions"
        subtitle="Save big with our latest offers"
        icon={Zap}
        breadcrumbs={[{ label: 'Deals' }]}
      >
        <div className="space-y-8">
          <div className="space-y-4">
            <Skeleton className="h-8 w-48" />
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-40 w-full rounded-2xl" />
              ))}
            </div>
          </div>
          <div className="space-y-4">
            <Skeleton className="h-8 w-48" />
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-32 w-full rounded-2xl" />
              ))}
            </div>
          </div>
        </div>
      </PageLayout>
    );
  }

  return (
    <PageLayout
      title="Deals & Promotions"
      subtitle="Save big with our latest offers"
      icon={Zap}
      breadcrumbs={[{ label: 'Deals' }]}
    >
      <div className="space-y-10">
        {/* Flash Sales Section */}
        <section>
          <div className="mb-6 flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-orange-100 text-orange-600">
              <Zap className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-xl font-semibold tracking-tight">Flash Sales</h2>
              <p className="text-sm text-muted-foreground">Limited-time deals you don't want to miss</p>
            </div>
          </div>

          {flashSales.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed py-16">
              <div className="flex h-14 w-14 items-center justify-center rounded-full bg-orange-50 text-orange-400">
                <Sparkles className="h-7 w-7" />
              </div>
              <p className="mt-4 text-sm font-medium text-muted-foreground">
                No active flash sales right now
              </p>
              <p className="mt-1 text-xs text-muted-foreground/70">
                Check back soon for exciting deals!
              </p>
            </div>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {flashSales.map((sale: any) => (
                <div
                  key={sale.id}
                  className="group relative overflow-hidden rounded-2xl bg-gradient-to-br from-orange-50 to-amber-50 p-6 transition-shadow hover:shadow-lg"
                >
                  <div className="absolute right-3 top-3">
                    <Badge
                      variant="secondary"
                      className="bg-orange-100 text-orange-700 hover:bg-orange-100"
                    >
                      {sale.status || 'active'}
                    </Badge>
                  </div>
                  <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-100 text-orange-600">
                    <Zap className="h-5 w-5" />
                  </div>
                  <h3 className="mt-4 text-lg font-semibold text-foreground">{sale.name}</h3>
                  <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                    {sale.description}
                  </p>
                  {sale.discount_percent && (
                    <p className="mt-3 text-2xl font-bold text-orange-600">
                      {sale.discount_percent}% OFF
                    </p>
                  )}
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Coupons Section */}
        <section>
          <div className="mb-6 flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-violet-100 text-violet-600">
              <Ticket className="h-5 w-5" />
            </div>
            <div>
              <h2 className="text-xl font-semibold tracking-tight">Available Coupons</h2>
              <p className="text-sm text-muted-foreground">Apply at checkout to save</p>
            </div>
          </div>

          {coupons.length === 0 ? (
            <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed py-16">
              <div className="flex h-14 w-14 items-center justify-center rounded-full bg-violet-50 text-violet-400">
                <Tag className="h-7 w-7" />
              </div>
              <p className="mt-4 text-sm font-medium text-muted-foreground">
                No coupons available
              </p>
              <p className="mt-1 text-xs text-muted-foreground/70">
                New coupons are added regularly
              </p>
            </div>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {coupons.map((coupon: any) => (
                <div
                  key={coupon.id}
                  className="rounded-2xl border-2 border-dashed border-violet-200 bg-white p-6 transition-colors hover:border-violet-300 hover:bg-violet-50/30"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="inline-flex items-center rounded-lg bg-violet-100 px-3 py-1.5">
                        <span className="font-mono text-sm font-bold tracking-wider text-violet-700">
                          {coupon.code}
                        </span>
                      </div>
                      <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
                        {coupon.description}
                      </p>
                    </div>
                    <Badge className="ml-3 bg-violet-100 text-violet-700 hover:bg-violet-100">
                      {coupon.type}
                    </Badge>
                  </div>
                  {coupon.discount_value && (
                    <p className="mt-3 text-lg font-semibold text-violet-600">
                      {coupon.type === 'percentage'
                        ? `${coupon.discount_value}% off`
                        : `$${coupon.discount_value} off`}
                    </p>
                  )}
                </div>
              ))}
            </div>
          )}
        </section>
      </div>
    </PageLayout>
  );
}
