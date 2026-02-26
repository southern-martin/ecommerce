import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminTaxApi } from '../services/admin-tax.api';
import type { TaxRule } from '../services/admin-tax.api';

export function useAdminTaxRules() {
  return useQuery({
    queryKey: ['admin-tax-rules'],
    queryFn: adminTaxApi.getRules,
  });
}

export function useCreateTaxRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<TaxRule>) => adminTaxApi.createRule(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-tax-rules'] });
    },
  });
}

export function useUpdateTaxRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<TaxRule> }) =>
      adminTaxApi.updateRule(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-tax-rules'] });
    },
  });
}

export function useDeleteTaxRule() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminTaxApi.deleteRule(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-tax-rules'] });
    },
  });
}
