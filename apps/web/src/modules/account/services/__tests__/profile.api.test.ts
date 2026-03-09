import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('@/shared/lib/api-client', () => ({
  default: {
    get: vi.fn(),
    patch: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';
import { profileApi } from '../profile.api';

const mockApiClient = apiClient as unknown as {
  get: ReturnType<typeof vi.fn>;
  patch: ReturnType<typeof vi.fn>;
};

describe('profileApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getProfile', () => {
    it('sends GET to /users/me and returns user data', async () => {
      const mockUser = { id: 'u1', email: 'test@example.com', first_name: 'John' };
      mockApiClient.get.mockResolvedValue({ data: { data: mockUser } });

      const result = await profileApi.getProfile();

      expect(mockApiClient.get).toHaveBeenCalledWith('/users/me');
      expect(result).toEqual(mockUser);
    });
  });

  describe('updateProfile', () => {
    it('sends PATCH to /users/me with update data', async () => {
      const updateData = { first_name: 'Jane', phone: '555-1234' };
      const mockUser = { id: 'u1', email: 'test@example.com', first_name: 'Jane', phone: '555-1234' };
      mockApiClient.patch.mockResolvedValue({ data: { data: mockUser } });

      const result = await profileApi.updateProfile(updateData);

      expect(mockApiClient.patch).toHaveBeenCalledWith('/users/me', updateData);
      expect(result).toEqual(mockUser);
    });

    it('sends partial update with only changed fields', async () => {
      const updateData = { avatar_url: 'https://example.com/avatar.jpg' };
      mockApiClient.patch.mockResolvedValue({ data: { data: { id: 'u1', avatar_url: updateData.avatar_url } } });

      await profileApi.updateProfile(updateData);

      expect(mockApiClient.patch).toHaveBeenCalledWith('/users/me', updateData);
    });
  });
});
