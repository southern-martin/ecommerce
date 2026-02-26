import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { loyaltyApi } from '../services/loyalty.api';

export function usePointsHistory() {
  return useQuery({
    queryKey: ['loyalty', 'points-history'],
    queryFn: () => loyaltyApi.getPointsHistory(),
  });
}

export function useRedeemPoints() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (points: number) => loyaltyApi.redeemPoints(points),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['loyalty'] });
    },
  });
}
