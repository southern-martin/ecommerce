import { describe, it, expect, vi, beforeEach } from 'vitest';
import { authApi } from '../auth.api';

// Mock the api-client module
vi.mock('@/shared/lib/api-client', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn(),
  },
}));

import apiClient from '@/shared/lib/api-client';

const mockApiClient = apiClient as unknown as {
  post: ReturnType<typeof vi.fn>;
  get: ReturnType<typeof vi.fn>;
};

describe('authApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('login', () => {
    it('sends POST to /auth/login with credentials', async () => {
      const mockResponse = {
        data: {
          user_id: 'user-1',
          email: 'test@example.com',
          access_token: 'token-abc',
          refresh_token: 'refresh-abc',
          role: 'buyer',
        },
      };
      mockApiClient.post.mockResolvedValue(mockResponse);

      const result = await authApi.login({
        email: 'test@example.com',
        password: 'password123',
      });

      expect(mockApiClient.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@example.com',
        password: 'password123',
      });
      expect(result.user_id).toBe('user-1');
      expect(result.access_token).toBe('token-abc');
    });
  });

  describe('register', () => {
    it('sends POST to /auth/register', async () => {
      const mockResponse = {
        data: {
          user_id: 'user-2',
          email: 'new@example.com',
          access_token: 'new-token',
          refresh_token: 'new-refresh',
          role: 'buyer',
        },
      };
      mockApiClient.post.mockResolvedValue(mockResponse);

      const result = await authApi.register({
        email: 'new@example.com',
        password: 'password123',
      });

      expect(mockApiClient.post).toHaveBeenCalledWith('/auth/register', {
        email: 'new@example.com',
        password: 'password123',
      });
      expect(result.user_id).toBe('user-2');
    });
  });

  describe('refreshToken', () => {
    it('sends POST to /auth/refresh', async () => {
      const mockResponse = {
        data: {
          access_token: 'new-access',
          refresh_token: 'new-refresh',
        },
      };
      mockApiClient.post.mockResolvedValue(mockResponse);

      const result = await authApi.refreshToken('old-refresh');

      expect(mockApiClient.post).toHaveBeenCalledWith('/auth/refresh', {
        refresh_token: 'old-refresh',
      });
      expect(result.access_token).toBe('new-access');
    });
  });

  describe('forgotPassword', () => {
    it('sends POST to /auth/forgot-password', async () => {
      mockApiClient.post.mockResolvedValue({ data: {} });

      await authApi.forgotPassword({ email: 'test@example.com' });

      expect(mockApiClient.post).toHaveBeenCalledWith('/auth/forgot-password', {
        email: 'test@example.com',
      });
    });
  });

  describe('resetPassword', () => {
    it('sends POST to /auth/reset-password', async () => {
      mockApiClient.post.mockResolvedValue({ data: {} });

      await authApi.resetPassword({
        token: 'reset-token',
        password: 'newpassword123',
        confirm_password: 'newpassword123',
      });

      expect(mockApiClient.post).toHaveBeenCalledWith('/auth/reset-password', {
        token: 'reset-token',
        password: 'newpassword123',
        confirm_password: 'newpassword123',
      });
    });
  });
});
