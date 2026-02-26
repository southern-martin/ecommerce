import { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { SearchBar } from '../components/SearchBar';
import { SearchResults } from '../components/SearchResults';
import { SortDropdown } from '@/modules/shop/components/SortDropdown';
import { useSearch } from '../hooks/useSearch';
import type { SortOption } from '@/modules/shop/types/shop.types';

export default function SearchResultsPage() {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q') ?? '';
  const [sort, setSort] = useState<SortOption>('newest');
  const [page, setPage] = useState(1);

  const { data, isLoading } = useSearch(query, { sort, page, page_size: 20 });
  const totalPages = data ? Math.ceil(data.total / 20) : 0;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <SearchBar />
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">
          {query ? `Results for "${query}"` : 'Search'}
        </h1>
        <SortDropdown value={sort} onChange={setSort} />
      </div>

      <SearchResults
        results={data?.data ?? []}
        isLoading={isLoading}
        total={data?.total}
      />

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      )}
    </div>
  );
}
