import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { affiliateApi } from '../services/affiliate.api';

export function useAffiliateStats() {
  return useQuery({
    queryKey: ['affiliate', 'stats'],
    queryFn: () => affiliateApi.getStats(),
  });
}

export function useAffiliateClicks() {
  return useQuery({
    queryKey: ['affiliate', 'clicks'],
    queryFn: () => affiliateApi.getClicks(),
  });
}

export function useGenerateLink() {
  return useMutation({
    mutationFn: (productId?: string) => affiliateApi.generateLink(productId),
  });
}

export function useRequestPayout() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => affiliateApi.requestPayout(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['affiliate'] });
    },
  });
}
