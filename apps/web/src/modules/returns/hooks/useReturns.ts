import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { returnApi } from '../services/return.api';
import type { CreateReturnData } from '../services/return.api';

export function useReturns(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['returns', page, pageSize],
    queryFn: () => returnApi.getReturns({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useReturn(id: string) {
  return useQuery({
    queryKey: ['return', id],
    queryFn: () => returnApi.getReturnById(id),
    enabled: !!id,
  });
}

export function useCreateReturn() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateReturnData) => returnApi.createReturn(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['returns'] });
    },
  });
}
