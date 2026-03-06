import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useAuthStore } from '../auth.store';
import type { User } from '../../types/user.types';

const mockUser: User = {
  id: 'user-123',
  email: 'test@example.com',
  first_name: 'John',
  last_name: 'Doe',
  role: 'buyer',
  is_verified: true,
  created_at: '2026-01-01T00:00:00Z',
};

const mockTokens = {
  access_token: 'access-token-xyz',
  refresh_token: 'refresh-token-xyz',
  expires_in: 900,
};

describe('useAuthStore', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
    localStorage.clear();
  });

  describe('setAuth', () => {
    it('sets user, tokens, and isAuthenticated', () => {
      useAuthStore.getState().setAuth(mockUser, mockTokens);

      const state = useAuthStore.getState();
      expect(state.user).toEqual(mockUser);
      expect(state.accessToken).toBe('access-token-xyz');
      expect(state.refreshToken).toBe('refresh-token-xyz');
      expect(state.isAuthenticated).toBe(true);
    });

    it('stores tokens in localStorage', () => {
      useAuthStore.getState().setAuth(mockUser, mockTokens);

      expect(localStorage.getItem('access_token')).toBe('access-token-xyz');
      expect(localStorage.getItem('refresh_token')).toBe('refresh-token-xyz');
    });
  });

  describe('setUser', () => {
    it('updates user without changing tokens', () => {
      useAuthStore.getState().setAuth(mockUser, mockTokens);

      const updatedUser = { ...mockUser, first_name: 'Jane' };
      useAuthStore.getState().setUser(updatedUser);

      const state = useAuthStore.getState();
      expect(state.user?.first_name).toBe('Jane');
      expect(state.accessToken).toBe('access-token-xyz');
    });
  });

  describe('logout', () => {
    it('clears all auth state', () => {
      useAuthStore.getState().setAuth(mockUser, mockTokens);
      useAuthStore.getState().logout();

      const state = useAuthStore.getState();
      expect(state.user).toBeNull();
      expect(state.accessToken).toBeNull();
      expect(state.refreshToken).toBeNull();
      expect(state.isAuthenticated).toBe(false);
    });

    it('removes tokens from localStorage', () => {
      useAuthStore.getState().setAuth(mockUser, mockTokens);
      useAuthStore.getState().logout();

      expect(localStorage.getItem('access_token')).toBeNull();
      expect(localStorage.getItem('refresh_token')).toBeNull();
    });
  });
});
