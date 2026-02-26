import { useQuery } from '@tanstack/react-query';
import { searchApi } from '../services/search.api';

export function useSearchSuggestions(query: string) {
  return useQuery({
    queryKey: ['search-suggestions', query],
    queryFn: () => searchApi.suggest(query),
    enabled: query.length >= 2,
    staleTime: 30 * 1000,
  });
}
