import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { RoleGuard } from '../RoleGuard';
import { useAuthStore } from '../../stores/auth.store';
import type { User } from '../../types/user.types';

const createUser = (role: 'buyer' | 'seller' | 'admin'): User => ({
  id: `user-${role}`,
  email: `${role}@example.com`,
  first_name: 'Test',
  last_name: 'User',
  role,
  is_verified: true,
  created_at: '2026-01-01T00:00:00Z',
});

describe('RoleGuard', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
  });

  it('renders children when user has allowed role', () => {
    useAuthStore.setState({ user: createUser('admin'), isAuthenticated: true });

    render(
      <RoleGuard allowedRoles={['admin']}>
        <div>Admin Panel</div>
      </RoleGuard>
    );

    expect(screen.getByText('Admin Panel')).toBeInTheDocument();
  });

  it('redirects when user role is not allowed', () => {
    useAuthStore.setState({ user: createUser('buyer'), isAuthenticated: true });

    render(
      <RoleGuard allowedRoles={['admin', 'seller']}>
        <div>Admin Panel</div>
      </RoleGuard>
    );

    expect(screen.queryByText('Admin Panel')).not.toBeInTheDocument();
  });

  it('redirects to login when not authenticated', () => {
    render(
      <RoleGuard allowedRoles={['admin']}>
        <div>Admin Panel</div>
      </RoleGuard>
    );

    expect(screen.queryByText('Admin Panel')).not.toBeInTheDocument();
  });

  it('accepts roles prop as alias for allowedRoles', () => {
    useAuthStore.setState({ user: createUser('seller'), isAuthenticated: true });

    render(
      <RoleGuard roles={['seller']}>
        <div>Seller Dashboard</div>
      </RoleGuard>
    );

    expect(screen.getByText('Seller Dashboard')).toBeInTheDocument();
  });

  it('allows multiple roles', () => {
    useAuthStore.setState({ user: createUser('seller'), isAuthenticated: true });

    render(
      <RoleGuard allowedRoles={['seller', 'admin']}>
        <div>Privileged Content</div>
      </RoleGuard>
    );

    expect(screen.getByText('Privileged Content')).toBeInTheDocument();
  });
});
