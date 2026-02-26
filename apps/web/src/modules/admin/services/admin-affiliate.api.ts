import apiClient from '@/shared/lib/api-client';

export interface AffiliateProgram {
  commission_rate: number;
  cookie_duration_days: number;
  min_payout_amount: number;
  payout_schedule: string;
}

export interface AffiliatePayout {
  id: string;
  user_id: string;
  user_email: string;
  amount: number;
  status: string;
  requested_at: string;
}

export const adminAffiliateApi = {
  getProgram: async (): Promise<AffiliateProgram> => {
    const response = await apiClient.get('/admin/affiliates/program');
    return response.data.data ?? response.data;
  },
  updateProgram: async (data: Partial<AffiliateProgram>): Promise<AffiliateProgram> => {
    const response = await apiClient.patch('/admin/affiliates/program', data);
    return response.data.data ?? response.data;
  },
  getPayouts: async (params?: { page?: number }) => {
    const response = await apiClient.get('/admin/affiliates/payouts', { params });
    return response.data;
  },
  updatePayoutStatus: async (id: string, status: string) => {
    const response = await apiClient.patch(`/admin/affiliates/payouts/${id}`, { status });
    return response.data;
  },
};
