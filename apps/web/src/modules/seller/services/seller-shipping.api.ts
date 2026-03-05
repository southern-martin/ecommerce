import apiClient from '@/shared/lib/api-client';

export interface Shipment {
  id: string;
  order_id: string;
  tracking_number: string;
  carrier: string;
  status: string;
  estimated_delivery: string;
  created_at: string;
}

export interface SellerCarrier {
  id: string;
  carrier_code: string;
  carrier_name: string;
  account_number: string;
  is_active: boolean;
}

export interface ShippingAddress {
  street: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
}

export interface ShipmentItemInput {
  product_id: string;
  variant_id?: string;
  product_name: string;
  quantity: number;
}

export interface CreateShipmentInput {
  order_id: string;
  carrier_code: string;
  service_code?: string;
  origin: ShippingAddress;
  destination: ShippingAddress;
  weight_grams: number;
  rate_cents?: number;
  currency?: string;
  items: ShipmentItemInput[];
}

export const sellerShippingApi = {
  getShipments: async (params?: { page?: number; page_size?: number }) => {
    const response = await apiClient.get('/seller/shipments', { params });
    return response.data;
  },
  getCarriers: async () => {
    const response = await apiClient.get('/seller/carriers');
    return response.data.data ?? response.data;
  },
  setupCarrier: async (data: { carrier_code: string; account_number: string }) => {
    const response = await apiClient.post('/seller/carriers', data);
    return response.data.data ?? response.data;
  },
  createShipment: async (data: CreateShipmentInput): Promise<Shipment> => {
    const response = await apiClient.post('/shipments', data);
    return response.data.shipment ?? response.data.data ?? response.data;
  },
  getShipmentsByOrderId: async (orderId: string): Promise<Shipment | null> => {
    const response = await apiClient.get('/seller/shipments', {
      params: { order_id: orderId, page_size: 1 },
    });
    const shipments = response.data?.data ?? [];
    return shipments.length > 0 ? shipments[0] : null;
  },
};
