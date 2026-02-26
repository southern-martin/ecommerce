import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface Coupon {
  id: string;
  code: string;
  type: 'percentage' | 'fixed_amount' | 'free_shipping';
  value: number;
  min_order_amount?: number;
  max_discount?: number;
  expires_at: string;
  is_active: boolean;
}

export interface CouponValidation {
  valid: boolean;
  discount_amount: number;
  message?: string;
}

export const couponApi = {
  validateCoupon: async (code: string, orderTotal: number): Promise<CouponValidation> => {
    const response = await apiClient.post<ApiResponse<CouponValidation>>('/coupons/validate', {
      code,
      order_total: orderTotal,
    });
    return response.data.data;
  },

  getActiveCoupons: async (): Promise<Coupon[]> => {
    const response = await apiClient.get<ApiResponse<Coupon[]>>('/coupons/active');
    return response.data.data;
  },
};
