import { useQuery } from '@tanstack/react-query';
import { flashSaleApi } from '../services/flash-sale.api';

export function useFlashSales() {
  return useQuery({
    queryKey: ['flash-sales', 'active'],
    queryFn: () => flashSaleApi.getActiveFlashSales(),
    refetchInterval: 60000,
  });
}

export function useFlashSale(id: string) {
  return useQuery({
    queryKey: ['flash-sale', id],
    queryFn: () => flashSaleApi.getFlashSaleById(id),
    enabled: !!id,
  });
}
