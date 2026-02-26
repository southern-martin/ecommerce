import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';

export interface Banner {
  id: string;
  title: string;
  subtitle?: string;
  image_url: string;
  link_url?: string;
  is_active: boolean;
  sort_order: number;
  created_at: string;
}

export interface CreateBannerData {
  title: string;
  subtitle?: string;
  image_url: string;
  link_url?: string;
  is_active?: boolean;
  sort_order?: number;
}

export interface Page {
  id: string;
  title: string;
  slug: string;
  content: string;
  meta_title?: string;
  meta_description?: string;
  published: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreatePageData {
  title: string;
  slug: string;
  content: string;
  meta_title?: string;
  meta_description?: string;
}

export const adminCmsApi = {
  // Banners
  getBanners: async (): Promise<Banner[]> => {
    const response = await apiClient.get<ApiResponse<Banner[]>>('/admin/banners');
    return response.data.data;
  },

  createBanner: async (data: CreateBannerData): Promise<Banner> => {
    const response = await apiClient.post<ApiResponse<Banner>>('/admin/banners', data);
    return response.data.data;
  },

  updateBanner: async (id: string, data: Partial<CreateBannerData>): Promise<Banner> => {
    const response = await apiClient.patch<ApiResponse<Banner>>(`/admin/banners/${id}`, data);
    return response.data.data;
  },

  deleteBanner: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/banners/${id}`);
  },

  // Pages
  getPages: async (): Promise<Page[]> => {
    const response = await apiClient.get<ApiResponse<Page[]>>('/admin/pages');
    return response.data.data;
  },

  createPage: async (data: CreatePageData): Promise<Page> => {
    const response = await apiClient.post<ApiResponse<Page>>('/admin/pages', data);
    return response.data.data;
  },

  updatePage: async (id: string, data: Partial<CreatePageData>): Promise<Page> => {
    const response = await apiClient.patch<ApiResponse<Page>>(`/admin/pages/${id}`, data);
    return response.data.data;
  },

  deletePage: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/pages/${id}`);
  },

  publishPage: async (id: string): Promise<Page> => {
    const response = await apiClient.patch<ApiResponse<Page>>(`/admin/pages/${id}/publish`);
    return response.data.data;
  },
};
