import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { searchApi } from '../services/search.api';
import type { FilterState } from '@/modules/shop/types/shop.types';

export function useSearch(query: string, filters?: Partial<FilterState>) {
  return useQuery({
    queryKey: ['search', query, filters],
    queryFn: () => searchApi.search(query, filters),
    enabled: query.length >= 2,
    placeholderData: keepPreviousData,
  });
}
