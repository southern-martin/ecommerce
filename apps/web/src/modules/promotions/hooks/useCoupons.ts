import { useQuery, useMutation } from '@tanstack/react-query';
import { couponApi } from '../services/coupon.api';

export function useCoupons() {
  return useQuery({
    queryKey: ['coupons', 'active'],
    queryFn: () => couponApi.getActiveCoupons(),
  });
}

export function useValidateCoupon() {
  return useMutation({
    mutationFn: ({ code, orderTotal }: { code: string; orderTotal: number }) =>
      couponApi.validateCoupon(code, orderTotal),
  });
}
