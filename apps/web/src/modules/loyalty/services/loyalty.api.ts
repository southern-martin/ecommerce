import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface Membership {
  id: string;
  tier: 'bronze' | 'silver' | 'gold' | 'platinum';
  points_balance: number;
  total_points_earned: number;
  points_to_next_tier: number;
  next_tier: string | null;
  tier_progress_percentage: number;
  member_since: string;
}

export interface PointsTransaction {
  id: string;
  type: 'earned' | 'redeemed' | 'expired';
  points: number;
  description: string;
  created_at: string;
}

export const loyaltyApi = {
  getMembership: async (): Promise<Membership> => {
    const response = await apiClient.get<ApiResponse<Membership>>('/loyalty/membership');
    return response.data.data;
  },

  getPointsHistory: async (): Promise<PointsTransaction[]> => {
    const response = await apiClient.get<ApiResponse<PointsTransaction[]>>('/loyalty/points/history');
    return response.data.data;
  },

  redeemPoints: async (points: number): Promise<{ discount_code: string; value: number }> => {
    const response = await apiClient.post<ApiResponse<{ discount_code: string; value: number }>>(
      '/loyalty/points/redeem',
      { points }
    );
    return response.data.data;
  },
};
