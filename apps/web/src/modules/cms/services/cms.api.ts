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
}

export interface StaticPage {
  id: string;
  title: string;
  slug: string;
  content: string;
  meta_title?: string;
  meta_description?: string;
  published: boolean;
  updated_at: string;
}

export const cmsApi = {
  getBanners: async (): Promise<Banner[]> => {
    const response = await apiClient.get<ApiResponse<Banner[]>>('/cms/banners');
    return response.data.data;
  },

  getPages: async (): Promise<StaticPage[]> => {
    const response = await apiClient.get<ApiResponse<StaticPage[]>>('/cms/pages');
    return response.data.data;
  },

  getPageBySlug: async (slug: string): Promise<StaticPage> => {
    const response = await apiClient.get<ApiResponse<StaticPage>>(`/cms/pages/${slug}`);
    return response.data.data;
  },
};
