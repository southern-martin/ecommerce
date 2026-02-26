import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface AffiliateStats {
  total_clicks: number;
  total_conversions: number;
  total_earnings: number;
  conversion_rate: number;
  pending_payout: number;
  referral_code: string;
  referral_link: string;
}

export interface ReferralClick {
  id: string;
  ip_hash: string;
  converted: boolean;
  earned: number;
  created_at: string;
}

export const affiliateApi = {
  getStats: async (): Promise<AffiliateStats> => {
    const response = await apiClient.get<ApiResponse<AffiliateStats>>('/affiliate/stats');
    return response.data.data;
  },

  getClicks: async (): Promise<ReferralClick[]> => {
    const response = await apiClient.get<ApiResponse<ReferralClick[]>>('/affiliate/clicks');
    return response.data.data;
  },

  generateLink: async (productId?: string): Promise<{ link: string; code: string }> => {
    const response = await apiClient.post<ApiResponse<{ link: string; code: string }>>(
      '/affiliate/generate-link',
      { product_id: productId }
    );
    return response.data.data;
  },

  requestPayout: async (): Promise<void> => {
    await apiClient.post('/affiliate/payout');
  },
};
