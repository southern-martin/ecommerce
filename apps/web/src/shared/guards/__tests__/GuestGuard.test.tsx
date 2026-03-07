import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@/test/test-utils';
import { GuestGuard } from '../GuestGuard';
import { useAuthStore } from '../../stores/auth.store';

describe('GuestGuard', () => {
  beforeEach(() => {
    useAuthStore.setState({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
    });
  });

  it('renders children when not authenticated', () => {
    render(
      <GuestGuard>
        <div>Login Form</div>
      </GuestGuard>
    );

    expect(screen.getByText('Login Form')).toBeInTheDocument();
  });

  it('redirects when authenticated', () => {
    useAuthStore.setState({
      user: {
        id: 'user-1',
        email: 'test@example.com',
        first_name: 'Test',
        last_name: 'User',
        role: 'buyer',
        is_verified: true,
        created_at: '2026-01-01T00:00:00Z',
      },
      isAuthenticated: true,
    });

    render(
      <GuestGuard>
        <div>Login Form</div>
      </GuestGuard>
    );

    expect(screen.queryByText('Login Form')).not.toBeInTheDocument();
  });
});
