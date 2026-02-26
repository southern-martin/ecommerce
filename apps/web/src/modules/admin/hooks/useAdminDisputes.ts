import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminDisputeApi } from '../services/admin-dispute.api';
import type { ResolveDisputeData } from '../services/admin-dispute.api';

export function useAdminDisputes(page = 1) {
  return useQuery({
    queryKey: ['admin-disputes', page],
    queryFn: () => adminDisputeApi.getDisputes({ page, page_size: 10 }),
    placeholderData: keepPreviousData,
  });
}

export function useResolveDispute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ResolveDisputeData }) =>
      adminDisputeApi.resolveDispute(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-disputes'] });
    },
  });
}
