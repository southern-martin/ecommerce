import { ProductGrid } from '@/modules/shop/components/ProductGrid';
import type { Product } from '@/modules/shop/types/shop.types';

interface SearchResultsProps {
  results: Product[];
  isLoading?: boolean;
  total?: number;
}

export function SearchResults({ results, isLoading, total }: SearchResultsProps) {
  return (
    <div>
      {total !== undefined && !isLoading && (
        <p className="mb-4 text-sm text-muted-foreground">
          {total} result{total !== 1 ? 's' : ''} found
        </p>
      )}
      <ProductGrid products={results} isLoading={isLoading} />
    </div>
  );
}
