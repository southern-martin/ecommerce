import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminPromotionApi } from '../services/admin-promotion.api';
import type {
  CreateCouponData,
  CreateFlashSaleData,
  CreateBundleData,
} from '../services/admin-promotion.api';

// Coupons
export function useAdminCoupons() {
  return useQuery({
    queryKey: ['admin-coupons'],
    queryFn: () => adminPromotionApi.getCoupons(),
  });
}

export function useCreateCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateCouponData) => adminPromotionApi.createCoupon(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-coupons'] });
    },
  });
}

export function useUpdateCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateCouponData> }) =>
      adminPromotionApi.updateCoupon(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-coupons'] });
    },
  });
}

export function useDeleteCoupon() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminPromotionApi.deleteCoupon(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-coupons'] });
    },
  });
}

// Flash Sales
export function useAdminFlashSales() {
  return useQuery({
    queryKey: ['admin-flash-sales'],
    queryFn: () => adminPromotionApi.getFlashSales(),
  });
}

export function useCreateFlashSale() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateFlashSaleData) => adminPromotionApi.createFlashSale(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-flash-sales'] });
    },
  });
}

export function useUpdateFlashSale() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateFlashSaleData> }) =>
      adminPromotionApi.updateFlashSale(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-flash-sales'] });
    },
  });
}

// Bundles
export function useAdminBundles() {
  return useQuery({
    queryKey: ['admin-bundles'],
    queryFn: () => adminPromotionApi.getBundles(),
  });
}

export function useCreateBundle() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateBundleData) => adminPromotionApi.createBundle(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-bundles'] });
    },
  });
}

export function useUpdateBundle() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateBundleData> }) =>
      adminPromotionApi.updateBundle(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-bundles'] });
    },
  });
}

export function useDeleteBundle() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminPromotionApi.deleteBundle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-bundles'] });
    },
  });
}
