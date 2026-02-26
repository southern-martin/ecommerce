import { useQuery } from '@tanstack/react-query';
import { loyaltyApi } from '../services/loyalty.api';

export function useMembership() {
  return useQuery({
    queryKey: ['loyalty', 'membership'],
    queryFn: () => loyaltyApi.getMembership(),
  });
}
