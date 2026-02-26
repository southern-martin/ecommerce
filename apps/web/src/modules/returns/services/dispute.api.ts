import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface Dispute {
  id: string;
  return_id: string;
  reason: string;
  description: string;
  status: 'open' | 'under_review' | 'resolved' | 'closed';
  resolution?: string;
  created_at: string;
}

export const disputeApi = {
  createDispute: async (returnId: string, data: { reason: string; description: string }): Promise<Dispute> => {
    const response = await apiClient.post<ApiResponse<Dispute>>(`/returns/${returnId}/dispute`, data);
    return response.data.data;
  },

  getDispute: async (returnId: string): Promise<Dispute> => {
    const response = await apiClient.get<ApiResponse<Dispute>>(`/returns/${returnId}/dispute`);
    return response.data.data;
  },
};
