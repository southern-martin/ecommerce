import apiClient from '@/shared/lib/api-client';
import type { ApiResponse, PaginatedResponse } from '@/shared/types/api.types';

export interface Dispute {
  id: string;
  order_id: string;
  order_number: string;
  buyer_id: string;
  buyer_name: string;
  seller_id: string;
  seller_name: string;
  reason: string;
  status: 'open' | 'under_review' | 'resolved' | 'closed';
  resolution_type?: 'refund' | 'replacement' | 'rejected' | 'partial_refund';
  notes?: string;
  created_at: string;
}

export interface ResolveDisputeData {
  resolution_type: 'refund' | 'replacement' | 'rejected' | 'partial_refund';
  notes: string;
}

export const adminDisputeApi = {
  getDisputes: async (params?: {
    page?: number;
    page_size?: number;
  }): Promise<PaginatedResponse<Dispute>> => {
    const response = await apiClient.get<PaginatedResponse<Dispute>>('/admin/disputes', {
      params,
    });
    return response.data;
  },

  resolveDispute: async (id: string, data: ResolveDisputeData): Promise<Dispute> => {
    const response = await apiClient.patch<ApiResponse<Dispute>>(
      `/admin/disputes/${id}/resolve`,
      data
    );
    return response.data.data;
  },
};
