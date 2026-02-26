import { useQuery } from '@tanstack/react-query';
import { aiApi } from '../services/ai.api';

export function useRecommendations(params?: {
  product_id?: string;
  category?: string;
  limit?: number;
}) {
  return useQuery({
    queryKey: ['recommendations', params],
    queryFn: () => aiApi.getRecommendations(params),
    staleTime: 5 * 60 * 1000,
  });
}
