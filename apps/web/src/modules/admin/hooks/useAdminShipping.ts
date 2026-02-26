import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminShippingApi } from '../services/admin-shipping.api';
import type { Carrier } from '../services/admin-shipping.api';

export function useAdminCarriers() {
  return useQuery({
    queryKey: ['admin-carriers'],
    queryFn: adminShippingApi.getCarriers,
  });
}

export function useCreateCarrier() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<Carrier>) => adminShippingApi.createCarrier(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-carriers'] });
    },
  });
}

export function useUpdateCarrier() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ code, data }: { code: string; data: Partial<Carrier> }) =>
      adminShippingApi.updateCarrier(code, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-carriers'] });
    },
  });
}
