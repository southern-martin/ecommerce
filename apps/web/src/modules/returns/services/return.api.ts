import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse, ApiResponse } from '@/shared/types/api.types';

export interface ReturnRequest {
  id: string;
  order_id: string;
  order_number: string;
  reason: string;
  description: string;
  status: string;
  items: ReturnItem[];
  refund_amount: number;
  created_at: string;
  updated_at: string;
}

export interface ReturnItem {
  id: string;
  product_id: string;
  product_name: string;
  quantity: number;
  image_url: string;
}

export interface CreateReturnData {
  order_id: string;
  reason: string;
  description: string;
  items: { product_id: string; quantity: number }[];
}

export const returnApi = {
  getReturns: async (params: { page: number; page_size: number }): Promise<PaginatedResponse<ReturnRequest>> => {
    const response = await apiClient.get<PaginatedResponse<ReturnRequest>>('/returns', { params });
    return response.data;
  },

  createReturn: async (data: CreateReturnData): Promise<ReturnRequest> => {
    const response = await apiClient.post<ApiResponse<ReturnRequest>>('/returns', data);
    return response.data.data;
  },

  getReturnById: async (id: string): Promise<ReturnRequest> => {
    const response = await apiClient.get<ApiResponse<ReturnRequest>>(`/returns/${id}`);
    return response.data.data;
  },
};
