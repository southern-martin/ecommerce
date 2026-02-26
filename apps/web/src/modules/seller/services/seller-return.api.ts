import apiClient from '@/shared/lib/api-client';

export interface SellerReturn {
  id: string;
  order_id: string;
  order_number: string;
  reason: string;
  description: string;
  status: string;
  refund_amount: number;
  created_at: string;
}

export const sellerReturnApi = {
  getReturns: async (params?: { page?: number; page_size?: number; status?: string }) => {
    const response = await apiClient.get('/seller/returns', { params });
    return response.data;
  },
  approveReturn: async (id: string) => {
    const response = await apiClient.patch(`/seller/returns/${id}/approve`);
    return response.data;
  },
  rejectReturn: async (id: string) => {
    const response = await apiClient.patch(`/seller/returns/${id}/reject`);
    return response.data;
  },
  updateReturnStatus: async (id: string, status: string) => {
    const response = await apiClient.patch(`/seller/returns/${id}/status`, { status });
    return response.data;
  },
};
