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
};
