import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { sellerReturnApi } from '../services/seller-return.api';

export function useSellerReturns(page = 1, pageSize = 10, status?: string) {
  return useQuery({
    queryKey: ['seller-returns', page, pageSize, status],
    queryFn: () => sellerReturnApi.getReturns({ page, page_size: pageSize, status }),
    placeholderData: keepPreviousData,
  });
}

export function useApproveReturn() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => sellerReturnApi.approveReturn(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-returns'] });
    },
  });
}

export function useRejectReturn() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => sellerReturnApi.rejectReturn(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-returns'] });
    },
  });
}
