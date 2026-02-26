import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { accountOrderApi } from '../services/order.api';

export function useOrders(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['orders', page, pageSize],
    queryFn: () => accountOrderApi.getOrders({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useOrder(id: string) {
  return useQuery({
    queryKey: ['order', id],
    queryFn: () => accountOrderApi.getOrderById(id),
    enabled: !!id,
  });
}
