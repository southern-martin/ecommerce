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

export interface AttributeGroup {
  id: string;
  name: string;
  slug: string;
  description?: string;
  sort_order: number;
  attributes?: Attribute[];
  created_at: string;
  updated_at: string;
}

export interface CreateAttributeGroupData {
  name: string;
  description?: string;
  sort_order?: number;
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

  updateCategory: async (id: string, data: Partial<CreateCategoryData & { is_active?: boolean }>): Promise<Category> => {
    const response = await apiClient.patch(`/admin/categories/${id}`, data);
    return response.data;
  },

  deleteCategory: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/categories/${id}`);
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

  // Attribute Groups
  getAttributeGroups: async (): Promise<AttributeGroup[]> => {
    const response = await apiClient.get('/admin/attribute-groups');
    return response.data.attribute_groups || [];
  },

  createAttributeGroup: async (data: CreateAttributeGroupData): Promise<AttributeGroup> => {
    const response = await apiClient.post('/admin/attribute-groups', data);
    return response.data;
  },

  updateAttributeGroup: async (id: string, data: Partial<CreateAttributeGroupData>): Promise<AttributeGroup> => {
    const response = await apiClient.patch(`/admin/attribute-groups/${id}`, data);
    return response.data;
  },

  deleteAttributeGroup: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/attribute-groups/${id}`);
  },

  getGroupAttributes: async (groupId: string): Promise<Attribute[]> => {
    const response = await apiClient.get(`/attribute-groups/${groupId}/attributes`);
    return response.data.attributes || [];
  },

  addAttributeToGroup: async (
    groupId: string,
    data: { attribute_id: string; sort_order?: number }
  ): Promise<void> => {
    await apiClient.post(`/admin/attribute-groups/${groupId}/attributes`, data);
  },

  removeAttributeFromGroup: async (groupId: string, attrId: string): Promise<void> => {
    await apiClient.delete(`/admin/attribute-groups/${groupId}/attributes/${attrId}`);
  },
};
