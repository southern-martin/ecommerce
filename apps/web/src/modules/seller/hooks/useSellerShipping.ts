import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { sellerShippingApi, type CreateShipmentInput } from '../services/seller-shipping.api';

export function useSellerShipments(page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['seller-shipments', page, pageSize],
    queryFn: () => sellerShippingApi.getShipments({ page, page_size: pageSize }),
    placeholderData: keepPreviousData,
  });
}

export function useSellerCarriers() {
  return useQuery({
    queryKey: ['seller-carriers'],
    queryFn: () => sellerShippingApi.getCarriers(),
  });
}

export function useSetupCarrier() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { carrier_code: string; account_number: string }) =>
      sellerShippingApi.setupCarrier(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-carriers'] });
    },
  });
}

export function useCreateShipment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateShipmentInput) => sellerShippingApi.createShipment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-shipments'] });
      queryClient.invalidateQueries({ queryKey: ['seller-orders'] });
    },
  });
}

export function useShipmentByOrderId(orderId: string | undefined) {
  return useQuery({
    queryKey: ['seller-shipment-by-order', orderId],
    queryFn: () => sellerShippingApi.getShipmentsByOrderId(orderId!),
    enabled: !!orderId,
  });
}
