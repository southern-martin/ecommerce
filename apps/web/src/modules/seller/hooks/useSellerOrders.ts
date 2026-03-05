import { useMemo } from 'react';
import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { sellerOrderApi } from '../services/seller-order.api';
import { sellerProductApi } from '../services/seller-product.api';

export function useSellerOrders(page = 1, pageSize = 10, status?: string) {
  return useQuery({
    queryKey: ['seller-orders', page, pageSize, status],
    queryFn: () => sellerOrderApi.getOrders({ page, page_size: pageSize, status }),
    placeholderData: keepPreviousData,
  });
}

export function useSellerDashboardStats() {
  const { data: ordersData, isLoading: ordersLoading } = useQuery({
    queryKey: ['seller-orders', 1, 100],
    queryFn: () => sellerOrderApi.getOrders({ page: 1, page_size: 100 }),
  });

  const { data: productsData, isLoading: productsLoading } = useQuery({
    queryKey: ['seller-products', 1, 1],
    queryFn: () => sellerProductApi.getProducts({ page: 1, page_size: 1 }),
  });

  const stats = useMemo(() => {
    if (!ordersData) return null;
    const orders = ordersData.data ?? [];

    const total_revenue = orders.reduce((sum, o) => sum + (o.total ?? 0), 0);
    const total_orders = ordersData.total ?? orders.length;
    const total_products = productsData?.total ?? 0;
    const pending_orders = orders.filter((o) => o.status === 'pending').length;

    return { total_revenue, total_orders, total_products, pending_orders };
  }, [ordersData, productsData]);

  return { data: stats, isLoading: ordersLoading || productsLoading };
}
