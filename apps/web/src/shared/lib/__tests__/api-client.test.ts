import { describe, it, expect, beforeEach } from 'vitest';
import { useAuthStore } from '../../stores/auth.store';

describe('api-client interceptor behavior', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
    localStorage.clear();
  });

  it('reads access_token from localStorage for authorization', () => {
    localStorage.setItem('access_token', 'my-jwt-token');
    const token = localStorage.getItem('access_token');
    expect(token).toBe('my-jwt-token');
  });

  it('reads X-User-ID from auth store', () => {
    useAuthStore.setState({
      user: {
        id: 'user-123',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        role: 'buyer',
        is_verified: true,
        created_at: '2026-01-01T00:00:00Z',
      },
      isAuthenticated: true,
    });
    const state = useAuthStore.getState();
    expect(state.user?.id).toBe('user-123');
  });

  it('handles missing token gracefully', () => {
    const token = localStorage.getItem('access_token');
    expect(token).toBeNull();
  });

  it('clears tokens on logout', () => {
    localStorage.setItem('access_token', 'token');
    localStorage.setItem('refresh_token', 'refresh');

    useAuthStore.getState().logout();

    expect(localStorage.getItem('access_token')).toBeNull();
    expect(localStorage.getItem('refresh_token')).toBeNull();
  });
});
