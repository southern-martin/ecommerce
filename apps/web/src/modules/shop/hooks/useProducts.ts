import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { productApi } from '../services/product.api';
import type { FilterState } from '../types/shop.types';

export function useProducts(filters: Partial<FilterState>) {
  return useQuery({
    queryKey: ['products', filters],
    queryFn: () => productApi.getProducts(filters),
    placeholderData: keepPreviousData,
  });
}
