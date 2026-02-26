import apiClient from '@/shared/lib/api-client';

export interface Carrier {
  id?: string;
  code: string;
  name: string;
  tracking_url_template: string;
  is_active: boolean;
}

export const adminShippingApi = {
  getCarriers: async (): Promise<Carrier[]> => {
    const response = await apiClient.get('/admin/carriers');
    return response.data.data ?? response.data;
  },
  createCarrier: async (data: Partial<Carrier>): Promise<Carrier> => {
    const response = await apiClient.post('/admin/carriers', data);
    return response.data.data ?? response.data;
  },
  updateCarrier: async (code: string, data: Partial<Carrier>): Promise<Carrier> => {
    const response = await apiClient.patch(`/admin/carriers/${code}`, data);
    return response.data.data ?? response.data;
  },
};
