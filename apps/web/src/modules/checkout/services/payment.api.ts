import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface PaymentIntent {
  id: string;
  client_secret: string;
  amount: number;
  currency: string;
  status: string;
}

export const paymentApi = {
  createPaymentIntent: async (orderId: string): Promise<PaymentIntent> => {
    const response = await apiClient.post<ApiResponse<PaymentIntent>>(
      `/payments/intent`,
      { order_id: orderId }
    );
    return response.data.data;
  },
};
