import { Skeleton } from '@/shared/components/ui/skeleton';
import { Zap } from 'lucide-react';
import { FlashSaleBanner } from '../components/FlashSaleBanner';
import { ProductGrid } from '@/modules/shop/components/ProductGrid';
import { useFlashSales } from '../hooks/useFlashSales';

export default function FlashSalesPage() {
  const { data: flashSales, isLoading } = useFlashSales();

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-24 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-3">
        <Zap className="h-8 w-8 text-orange-500" />
        <h1 className="text-3xl font-bold">Flash Sales</h1>
      </div>

      {flashSales && flashSales.length > 0 ? (
        flashSales.map((sale) => (
          <div key={sale.id} className="space-y-4">
            <FlashSaleBanner sale={sale} />
            <ProductGrid products={sale.products} />
          </div>
        ))
      ) : (
        <div className="flex flex-col items-center py-16">
          <Zap className="h-12 w-12 text-muted-foreground/50" />
          <p className="mt-4 text-lg text-muted-foreground">No active flash sales right now</p>
          <p className="text-sm text-muted-foreground">Check back later for amazing deals!</p>
        </div>
      )}
    </div>
  );
}
