import { ProductCard } from '@/modules/shop/components/ProductCard';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { useRecommendations } from '../hooks/useRecommendations';

interface RecommendationCarouselProps {
  productId?: string;
  category?: string;
  title?: string;
}

export function RecommendationCarousel({
  productId,
  category,
  title = 'Recommended for You',
}: RecommendationCarouselProps) {
  const { data: products, isLoading } = useRecommendations({
    product_id: productId,
    category,
    limit: 8,
  });

  if (isLoading) {
    return (
      <div>
        <h3 className="mb-4 text-lg font-semibold">{title}</h3>
        <div className="flex gap-4 overflow-hidden">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="w-64 flex-shrink-0 space-y-2">
              <Skeleton className="aspect-square w-full rounded-lg" />
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-4 w-1/2" />
            </div>
          ))}
        </div>
      </div>
    );
  }

  if (!products || products.length === 0) return null;

  return (
    <div>
      <h3 className="mb-4 text-lg font-semibold">{title}</h3>
      <div className="flex gap-4 overflow-x-auto pb-4">
        {products.map((product) => (
          <div key={product.id} className="w-64 flex-shrink-0">
            <ProductCard product={product} />
          </div>
        ))}
      </div>
    </div>
  );
}
