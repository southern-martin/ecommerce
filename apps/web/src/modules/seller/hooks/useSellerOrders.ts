import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { sellerOrderApi } from '../services/seller-order.api';

export function useSellerOrders(page = 1, pageSize = 10, status?: string) {
  return useQuery({
    queryKey: ['seller-orders', page, pageSize, status],
    queryFn: () => sellerOrderApi.getOrders({ page, page_size: pageSize, status }),
    placeholderData: keepPreviousData,
  });
}

export function useSellerDashboard() {
  return useQuery({
    queryKey: ['seller-dashboard'],
    queryFn: () => sellerOrderApi.getDashboardStats(),
  });
}
