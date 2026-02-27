import apiClient from '@/shared/lib/api-client';

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  parent_id?: string;
  sort_order?: number;
  image_url?: string;
  is_active?: boolean;
  created_at: string;
}

export interface CreateCategoryData {
  name: string;
  description?: string;
  parent_id?: string;
  sort_order?: number;
  image_url?: string;
}

export interface Attribute {
  id: string;
  name: string;
  slug?: string;
  type: 'text' | 'number' | 'select' | 'boolean' | 'multi_select' | 'color' | 'bool';
  required: boolean;
  filterable: boolean;
  options?: string[];
  unit?: string;
  created_at: string;
}

export interface CreateAttributeData {
  name: string;
  type: 'text' | 'number' | 'select' | 'boolean';
  required?: boolean;
  filterable?: boolean;
  options?: string[];
}

export interface CategoryAttribute {
  id: string;
  attribute_id: string;
  attribute: Attribute;
}

export const adminProductApi = {
  // Categories
  getCategories: async (): Promise<Category[]> => {
    const response = await apiClient.get('/categories');
    return response.data.categories || [];
  },

  createCategory: async (data: CreateCategoryData): Promise<Category> => {
    const response = await apiClient.post('/admin/categories', data);
    return response.data;
  },

  // Attributes
  getAttributes: async (): Promise<Attribute[]> => {
    const response = await apiClient.get('/admin/attributes');
    return response.data.attributes || [];
  },

  createAttribute: async (data: CreateAttributeData): Promise<Attribute> => {
    const response = await apiClient.post('/admin/attributes', data);
    return response.data;
  },

  updateAttribute: async (id: string, data: Partial<CreateAttributeData>): Promise<Attribute> => {
    const response = await apiClient.patch(`/admin/attributes/${id}`, data);
    return response.data;
  },

  deleteAttribute: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/attributes/${id}`);
  },

  // Category Attributes
  getCategoryAttributes: async (categoryId: string): Promise<CategoryAttribute[]> => {
    const response = await apiClient.get(`/categories/${categoryId}/attributes`);
    return response.data.attributes || response.data.data || [];
  },

  assignAttribute: async (
    categoryId: string,
    data: { attribute_id: string }
  ): Promise<CategoryAttribute> => {
    const response = await apiClient.post(`/categories/${categoryId}/attributes`, data);
    return response.data;
  },

  removeAttribute: async (categoryId: string, attrId: string): Promise<void> => {
    await apiClient.delete(`/categories/${categoryId}/attributes/${attrId}`);
  },
};
