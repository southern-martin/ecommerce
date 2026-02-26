import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface ShippingRate {
  id: string;
  carrier: string;
  service: string;
  estimated_days: number;
  price: number;
}

export interface TrackingEvent {
  timestamp: string;
  status: string;
  location: string;
  description: string;
}

export interface TrackingInfo {
  tracking_number: string;
  carrier: string;
  status: string;
  estimated_delivery: string;
  events: TrackingEvent[];
}

export const shippingApi = {
  getRates: async (params: {
    postal_code: string;
    country: string;
    weight: number;
  }): Promise<ShippingRate[]> => {
    const response = await apiClient.post<ApiResponse<ShippingRate[]>>('/shipping/rates', params);
    return response.data.data;
  },

  getTracking: async (trackingNumber: string): Promise<TrackingInfo> => {
    const response = await apiClient.get<ApiResponse<TrackingInfo>>(
      `/shipping/tracking/${trackingNumber}`
    );
    return response.data.data;
  },
};
