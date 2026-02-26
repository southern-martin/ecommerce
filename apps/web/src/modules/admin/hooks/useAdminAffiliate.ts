import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminAffiliateApi } from '../services/admin-affiliate.api';
import type { AffiliateProgram } from '../services/admin-affiliate.api';

export function useAffiliateProgram() {
  return useQuery({
    queryKey: ['admin-affiliate-program'],
    queryFn: adminAffiliateApi.getProgram,
  });
}

export function useUpdateAffiliateProgram() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<AffiliateProgram>) => adminAffiliateApi.updateProgram(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-affiliate-program'] });
    },
  });
}

export function useAffiliatePayouts(page = 1) {
  return useQuery({
    queryKey: ['admin-affiliate-payouts', page],
    queryFn: () => adminAffiliateApi.getPayouts({ page }),
    placeholderData: keepPreviousData,
  });
}

export function useUpdatePayoutStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      adminAffiliateApi.updatePayoutStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-affiliate-payouts'] });
    },
  });
}
