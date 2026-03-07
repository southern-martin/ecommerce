import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useAuth } from '../useAuth';
import { useAuthStore } from '../../stores/auth.store';
import type { User } from '../../types/user.types';

const buyerUser: User = {
  id: 'user-1',
  email: 'buyer@example.com',
  first_name: 'Jane',
  last_name: 'Doe',
  role: 'buyer',
  is_verified: true,
  created_at: '2026-01-01T00:00:00Z',
};

const sellerUser: User = {
  id: 'user-2',
  email: 'seller@example.com',
  first_name: 'John',
  last_name: 'Smith',
  role: 'seller',
  is_verified: true,
  created_at: '2026-01-01T00:00:00Z',
};

const adminUser: User = {
  id: 'user-3',
  email: 'admin@example.com',
  first_name: 'Admin',
  last_name: 'User',
  role: 'admin',
  is_verified: true,
  created_at: '2026-01-01T00:00:00Z',
};

describe('useAuth', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
    localStorage.clear();
  });

  it('returns isAuthenticated false when no user', () => {
    const { result } = renderHook(() => useAuth());
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
  });

  it('returns isAuthenticated true when user is set', () => {
    useAuthStore.getState().setAuth(buyerUser, {
      access_token: 'token',
      refresh_token: 'refresh',
      expires_in: 3600,
    });
    const { result } = renderHook(() => useAuth());
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.user).toEqual(buyerUser);
  });

  it('identifies buyer role', () => {
    useAuthStore.getState().setAuth(buyerUser, {
      access_token: 'token',
      refresh_token: 'refresh',
      expires_in: 3600,
    });
    const { result } = renderHook(() => useAuth());
    expect(result.current.isBuyer).toBe(true);
    expect(result.current.isSeller).toBe(false);
    expect(result.current.isAdmin).toBe(false);
  });

  it('identifies seller role', () => {
    useAuthStore.getState().setAuth(sellerUser, {
      access_token: 'token',
      refresh_token: 'refresh',
      expires_in: 3600,
    });
    const { result } = renderHook(() => useAuth());
    expect(result.current.isSeller).toBe(true);
    expect(result.current.isBuyer).toBe(false);
    expect(result.current.isAdmin).toBe(false);
  });

  it('identifies admin role', () => {
    useAuthStore.getState().setAuth(adminUser, {
      access_token: 'token',
      refresh_token: 'refresh',
      expires_in: 3600,
    });
    const { result } = renderHook(() => useAuth());
    expect(result.current.isAdmin).toBe(true);
    expect(result.current.isBuyer).toBe(false);
    expect(result.current.isSeller).toBe(false);
  });

  it('logout clears authentication', () => {
    useAuthStore.getState().setAuth(buyerUser, {
      access_token: 'token',
      refresh_token: 'refresh',
      expires_in: 3600,
    });
    const { result } = renderHook(() => useAuth());
    act(() => {
      result.current.logout();
    });
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.user).toBeNull();
  });
});
