import apiClient from '@/shared/lib/api-client';
import type { ApiResponse, PaginatedResponse } from '@/shared/types/api.types';

export interface Category {
  id: string;
  name: string;
  slug: string;
  description?: string;
  parent_id?: string;
  created_at: string;
}

export interface CreateCategoryData {
  name: string;
  slug: string;
  description?: string;
  parent_id?: string;
}

export interface Attribute {
  id: string;
  name: string;
  type: 'text' | 'number' | 'select' | 'boolean';
  required: boolean;
  filterable: boolean;
  options?: string[];
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
    const response = await apiClient.get<ApiResponse<Category[]>>('/categories');
    return response.data.data;
  },

  createCategory: async (data: CreateCategoryData): Promise<Category> => {
    const response = await apiClient.post<ApiResponse<Category>>('/admin/categories', data);
    return response.data.data;
  },

  // Attributes
  getAttributes: async (): Promise<Attribute[]> => {
    const response = await apiClient.get<ApiResponse<Attribute[]>>('/admin/attributes');
    return response.data.data;
  },

  createAttribute: async (data: CreateAttributeData): Promise<Attribute> => {
    const response = await apiClient.post<ApiResponse<Attribute>>('/admin/attributes', data);
    return response.data.data;
  },

  updateAttribute: async (id: string, data: Partial<CreateAttributeData>): Promise<Attribute> => {
    const response = await apiClient.patch<ApiResponse<Attribute>>(`/admin/attributes/${id}`, data);
    return response.data.data;
  },

  deleteAttribute: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/attributes/${id}`);
  },

  // Category Attributes
  getCategoryAttributes: async (categoryId: string): Promise<CategoryAttribute[]> => {
    const response = await apiClient.get<ApiResponse<CategoryAttribute[]>>(
      `/categories/${categoryId}/attributes`
    );
    return response.data.data;
  },

  assignAttribute: async (
    categoryId: string,
    data: { attribute_id: string }
  ): Promise<CategoryAttribute> => {
    const response = await apiClient.post<ApiResponse<CategoryAttribute>>(
      `/categories/${categoryId}/attributes`,
      data
    );
    return response.data.data;
  },

  removeAttribute: async (categoryId: string, attrId: string): Promise<void> => {
    await apiClient.delete(`/categories/${categoryId}/attributes/${attrId}`);
  },
};
