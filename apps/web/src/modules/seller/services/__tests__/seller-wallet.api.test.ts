import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockApiClient = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  patch: vi.fn(),
  delete: vi.fn(),
}));

vi.mock('@/shared/lib/api-client', () => ({
  default: mockApiClient,
}));

import { sellerWalletApi } from '../seller-wallet.api';

describe('sellerWalletApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getBalance', () => {
    it('should fetch wallet balance', async () => {
      const wallet = { seller_id: 's1', available_balance: 10000, pending_balance: 2000, currency: 'USD', updated_at: '2026-01-01' };
      mockApiClient.get.mockResolvedValue({ data: wallet });

      const result = await sellerWalletApi.getBalance();

      expect(mockApiClient.get).toHaveBeenCalledWith('/payments/wallet');
      expect(result).toEqual(wallet);
    });
  });

  describe('getTransactions', () => {
    it('should fetch paginated transactions', async () => {
      const response = { transactions: [{ id: 't1', amount_cents: 500, type: 'credit' }], total: 1, page: 1, page_size: 10 };
      mockApiClient.get.mockResolvedValue({ data: response });

      const result = await sellerWalletApi.getTransactions({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/payments/wallet/transactions', { params: { page: 1, page_size: 10 } });
      expect(result.transactions).toHaveLength(1);
      expect(result.total).toBe(1);
    });
  });

  describe('requestPayout', () => {
    it('should post a payout request and return the payout', async () => {
      const payout = { id: 'pay1', seller_id: 's1', amount_cents: 5000, currency: 'USD', method: 'bank', status: 'requested' };
      mockApiClient.post.mockResolvedValue({ data: payout });

      const result = await sellerWalletApi.requestPayout({ amount_cents: 5000 });

      expect(mockApiClient.post).toHaveBeenCalledWith('/payments/payouts', { amount_cents: 5000 });
      expect(result.status).toBe('requested');
      expect(result.amount_cents).toBe(5000);
    });
  });

  describe('getPayouts', () => {
    it('should fetch paginated payouts', async () => {
      const response = { payouts: [{ id: 'pay1', status: 'completed', amount_cents: 5000 }], total: 1, page: 1, page_size: 10 };
      mockApiClient.get.mockResolvedValue({ data: response });

      const result = await sellerWalletApi.getPayouts({ page: 1, page_size: 10 });

      expect(mockApiClient.get).toHaveBeenCalledWith('/payments/payouts', { params: { page: 1, page_size: 10 } });
      expect(result.payouts).toHaveLength(1);
      expect(result.total).toBe(1);
    });
  });
});
