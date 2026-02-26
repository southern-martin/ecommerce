import { useQuery } from '@tanstack/react-query';
import { productApi } from '../services/product.api';

export function useProduct(slug: string) {
  return useQuery({
    queryKey: ['product', slug],
    queryFn: () => productApi.getProductBySlug(slug),
    enabled: !!slug,
  });
}
