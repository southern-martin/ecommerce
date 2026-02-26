import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { ProductGrid } from '../components/ProductGrid';
import { FilterPanel } from '../components/FilterPanel';
import { SortDropdown } from '../components/SortDropdown';
import { useProducts } from '../hooks/useProducts';
import { useCategories } from '../hooks/useCategories';
import type { FilterState, SortOption } from '../types/shop.types';

export default function ProductListPage() {
  const [searchParams, setSearchParams] = useSearchParams();

  const [filters, setFilters] = useState<Partial<FilterState>>({
    category: searchParams.get('category') ?? undefined,
    sort: (searchParams.get('sort') as SortOption) ?? 'newest',
    page: Number(searchParams.get('page')) || 1,
    page_size: 20,
  });

  const { data, isLoading } = useProducts(filters);
  const { data: categories } = useCategories();

  const handleFilterChange = (newFilters: Partial<FilterState>) => {
    const updated = { ...filters, ...newFilters, page: 1 };
    setFilters(updated);
    const params = new URLSearchParams();
    Object.entries(updated).forEach(([key, value]) => {
      if (value !== undefined) params.set(key, String(value));
    });
    setSearchParams(params);
  };

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }));
  };

  const handleReset = () => {
    setFilters({ sort: 'newest', page: 1, page_size: 20 });
    setSearchParams({});
  };

  const totalPages = data ? Math.ceil(data.total / (filters.page_size ?? 20)) : 0;

  return (
    <div className="flex gap-8">
      <FilterPanel
        filters={filters}
        categories={categories ?? []}
        onFilterChange={handleFilterChange}
        onReset={handleReset}
      />

      <div className="flex-1">
        <div className="mb-6 flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            {data ? `${data.total} products found` : 'Loading...'}
          </p>
          <SortDropdown
            value={filters.sort ?? 'newest'}
            onChange={(sort) => handleFilterChange({ sort })}
          />
        </div>

        <ProductGrid products={data?.data ?? []} isLoading={isLoading} />

        {totalPages > 1 && (
          <div className="mt-8 flex items-center justify-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={filters.page === 1}
              onClick={() => handlePageChange((filters.page ?? 1) - 1)}
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {filters.page} of {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={filters.page === totalPages}
              onClick={() => handlePageChange((filters.page ?? 1) + 1)}
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
