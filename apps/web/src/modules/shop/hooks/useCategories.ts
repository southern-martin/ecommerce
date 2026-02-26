import { useQuery } from '@tanstack/react-query';
import { categoryApi } from '../services/category.api';

export function useCategories() {
  return useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.getCategories(),
    staleTime: 5 * 60 * 1000,
  });
}
