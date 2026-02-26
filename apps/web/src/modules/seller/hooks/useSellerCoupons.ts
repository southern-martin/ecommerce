import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { sellerCouponApi } from '../services/seller-coupon.api';
import type { Coupon } from '../services/seller-coupon.api';

export function useSellerCoupons(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['seller-coupons', page, pageSize],
    queryFn: () => sellerCouponApi.getCoupons({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useCreateCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: Partial<Coupon>) => sellerCouponApi.createCoupon(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-coupons'] });
    },
  });
}

export function useUpdateCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Coupon> }) =>
      sellerCouponApi.updateCoupon(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-coupons'] });
    },
  });
}

export function useDeleteCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => sellerCouponApi.deleteCoupon(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-coupons'] });
    },
  });
}
