import { useQuery } from '@tanstack/react-query';
import { shippingApi } from '../services/shipping.api';

export function useTracking(trackingNumber: string) {
  return useQuery({
    queryKey: ['tracking', trackingNumber],
    queryFn: () => shippingApi.getTracking(trackingNumber),
    enabled: !!trackingNumber,
    refetchInterval: 5 * 60 * 1000,
  });
}
