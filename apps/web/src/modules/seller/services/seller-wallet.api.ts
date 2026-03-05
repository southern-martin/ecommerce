import apiClient from '@/shared/lib/api-client';

export interface SellerWallet {
  seller_id: string;
  available_balance: number;
  pending_balance: number;
  currency: string;
  updated_at: string;
}

export interface WalletTransaction {
  id: string;
  seller_id: string;
  type: string;
  amount_cents: number;
  reference_type: string;
  reference_id: string;
  description: string;
  created_at: string;
}

export interface Payout {
  id: string;
  seller_id: string;
  amount_cents: number;
  currency: string;
  method: string;
  status: 'requested' | 'processing' | 'completed' | 'failed';
  requested_at: string;
  completed_at: string | null;
}

export interface RequestPayoutInput {
  amount_cents: number;
  currency?: string;
  method?: string;
}

export const sellerWalletApi = {
  getBalance: async (): Promise<SellerWallet> => {
    const response = await apiClient.get('/payments/wallet');
    return response.data;
  },
  getTransactions: async (params: { page: number; page_size: number }) => {
    const response = await apiClient.get('/payments/wallet/transactions', { params });
    return response.data as {
      transactions: WalletTransaction[];
      total: number;
      page: number;
      page_size: number;
    };
  },
  requestPayout: async (data: RequestPayoutInput): Promise<Payout> => {
    const response = await apiClient.post('/payments/payouts', data);
    return response.data;
  },
  getPayouts: async (params: { page: number; page_size: number }) => {
    const response = await apiClient.get('/payments/payouts', { params });
    return response.data as {
      payouts: Payout[];
      total: number;
      page: number;
      page_size: number;
    };
  },
};
