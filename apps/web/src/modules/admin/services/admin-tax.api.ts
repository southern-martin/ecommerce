import apiClient from '@/shared/lib/api-client';

export interface TaxRule {
  id: string;
  name: string;
  country: string;
  state: string;
  tax_rate: number;
  product_category: string;
  is_active: boolean;
}

export const adminTaxApi = {
  getRules: async (): Promise<TaxRule[]> => {
    const response = await apiClient.get('/admin/tax/rules');
    return response.data.data ?? response.data;
  },
  createRule: async (data: Partial<TaxRule>): Promise<TaxRule> => {
    const response = await apiClient.post('/admin/tax/rules', data);
    return response.data.data ?? response.data;
  },
  updateRule: async (id: string, data: Partial<TaxRule>): Promise<TaxRule> => {
    const response = await apiClient.patch(`/admin/tax/rules/${id}`, data);
    return response.data.data ?? response.data;
  },
  deleteRule: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/tax/rules/${id}`);
  },
};
