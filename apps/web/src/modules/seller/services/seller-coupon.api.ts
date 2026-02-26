import apiClient from '@/shared/lib/api-client';

export interface Coupon {
  id: string;
  code: string;
  type: 'percentage' | 'fixed_amount' | 'free_shipping';
  value: number;
  min_order_amount: number;
  max_uses: number;
  uses_count: number;
  starts_at: string;
  expires_at: string;
  is_active: boolean;
  created_at: string;
}

export const sellerCouponApi = {
  getCoupons: async (params?: { page?: number; page_size?: number }) => {
    const response = await apiClient.get('/seller/coupons', { params });
    return response.data;
  },
  getCoupon: async (id: string): Promise<Coupon> => {
    const response = await apiClient.get(`/seller/coupons/${id}`);
    return response.data.data ?? response.data;
  },
  createCoupon: async (data: Partial<Coupon>): Promise<Coupon> => {
    const response = await apiClient.post('/seller/coupons', data);
    return response.data.data ?? response.data;
  },
  updateCoupon: async (id: string, data: Partial<Coupon>): Promise<Coupon> => {
    const response = await apiClient.patch(`/seller/coupons/${id}`, data);
    return response.data.data ?? response.data;
  },
  deleteCoupon: async (id: string): Promise<void> => {
    await apiClient.delete(`/seller/coupons/${id}`);
  },
};
