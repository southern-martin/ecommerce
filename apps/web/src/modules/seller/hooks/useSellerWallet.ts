import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { sellerWalletApi, type RequestPayoutInput } from '../services/seller-wallet.api';

export function useSellerWalletBalance() {
  return useQuery({
    queryKey: ['seller-wallet-balance'],
    queryFn: () => sellerWalletApi.getBalance(),
  });
}

export function useSellerWalletTransactions(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['seller-wallet-transactions', page, pageSize],
    queryFn: () => sellerWalletApi.getTransactions({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useSellerPayouts(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['seller-payouts', page, pageSize],
    queryFn: () => sellerWalletApi.getPayouts({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useRequestPayout() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: RequestPayoutInput) => sellerWalletApi.requestPayout(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-wallet-balance'] });
      queryClient.invalidateQueries({ queryKey: ['seller-wallet-transactions'] });
      queryClient.invalidateQueries({ queryKey: ['seller-payouts'] });
    },
  });
}
