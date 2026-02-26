import { useMutation } from '@tanstack/react-query';
import { shippingApi } from '../services/shipping.api';

export function useShippingRates() {
  return useMutation({
    mutationFn: (params: { postal_code: string; country: string; weight: number }) =>
      shippingApi.getRates(params),
  });
}
