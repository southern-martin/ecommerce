import apiClient from '@/shared/lib/api-client';
import type { ApiResponse } from '@/shared/types/api.types';
import type { User } from '@/shared/types/user.types';

export interface UpdateProfileData {
  first_name?: string;
  last_name?: string;
  phone?: string;
  avatar_url?: string;
}

export const profileApi = {
  getProfile: async (): Promise<User> => {
    const response = await apiClient.get<ApiResponse<User>>('/users/me');
    return response.data.data;
  },

  updateProfile: async (data: UpdateProfileData): Promise<User> => {
    const response = await apiClient.patch<ApiResponse<User>>('/users/me', data);
    return response.data.data;
  },
};
