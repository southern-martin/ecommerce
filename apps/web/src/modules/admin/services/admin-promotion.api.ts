import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface AdminCoupon {
  id: string;
  code: string;
  type: 'percentage' | 'fixed_amount' | 'free_shipping';
  value: number;
  min_order_amount?: number;
  max_uses?: number;
  used_count: number;
  starts_at: string;
  expires_at: string;
  is_active: boolean;
  created_at: string;
}

export interface FlashSale {
  id: string;
  name: string;
  discount_percentage: number;
  starts_at: string;
  ends_at: string;
  is_active: boolean;
  created_at: string;
}

export interface Bundle {
  id: string;
  name: string;
  description?: string;
  discount_percentage: number;
  is_active: boolean;
  created_at: string;
}

export interface CreateCouponData {
  code: string;
  type: 'percentage' | 'fixed_amount' | 'free_shipping';
  value: number;
  min_order_amount?: number;
  max_uses?: number;
  starts_at: string;
  expires_at: string;
}

export interface CreateFlashSaleData {
  name: string;
  discount_percentage: number;
  starts_at: string;
  ends_at: string;
}

export interface CreateBundleData {
  name: string;
  description?: string;
  discount_percentage: number;
}

export const adminPromotionApi = {
  // Coupons
  getCoupons: async (): Promise<AdminCoupon[]> => {
    const response = await apiClient.get<ApiResponse<AdminCoupon[]>>('/admin/promotions/coupons');
    return response.data.data;
  },

  createCoupon: async (data: CreateCouponData): Promise<AdminCoupon> => {
    const response = await apiClient.post<ApiResponse<AdminCoupon>>(
      '/admin/promotions/coupons',
      data
    );
    return response.data.data;
  },

  updateCoupon: async (id: string, data: Partial<CreateCouponData>): Promise<AdminCoupon> => {
    const response = await apiClient.patch<ApiResponse<AdminCoupon>>(
      `/admin/promotions/coupons/${id}`,
      data
    );
    return response.data.data;
  },

  deleteCoupon: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/promotions/coupons/${id}`);
  },

  // Flash Sales
  getFlashSales: async (): Promise<FlashSale[]> => {
    const response = await apiClient.get<ApiResponse<FlashSale[]>>(
      '/admin/promotions/flash-sales'
    );
    return response.data.data;
  },

  createFlashSale: async (data: CreateFlashSaleData): Promise<FlashSale> => {
    const response = await apiClient.post<ApiResponse<FlashSale>>(
      '/admin/promotions/flash-sales',
      data
    );
    return response.data.data;
  },

  updateFlashSale: async (id: string, data: Partial<CreateFlashSaleData>): Promise<FlashSale> => {
    const response = await apiClient.patch<ApiResponse<FlashSale>>(
      `/admin/promotions/flash-sales/${id}`,
      data
    );
    return response.data.data;
  },

  // Bundles
  getBundles: async (): Promise<Bundle[]> => {
    const response = await apiClient.get<ApiResponse<Bundle[]>>('/admin/promotions/bundles');
    return response.data.data;
  },

  createBundle: async (data: CreateBundleData): Promise<Bundle> => {
    const response = await apiClient.post<ApiResponse<Bundle>>(
      '/admin/promotions/bundles',
      data
    );
    return response.data.data;
  },

  updateBundle: async (id: string, data: Partial<CreateBundleData>): Promise<Bundle> => {
    const response = await apiClient.patch<ApiResponse<Bundle>>(
      `/admin/promotions/bundles/${id}`,
      data
    );
    return response.data.data;
  },

  deleteBundle: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/promotions/bundles/${id}`);
  },
};
