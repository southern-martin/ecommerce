import { useMutation } from '@tanstack/react-query';
import { useNavigate, useLocation } from 'react-router-dom';
import { authApi } from '../services/auth.api';
import type { LoginRequest } from '../types/auth.types';
import { useAuthStore } from '@/shared/stores/auth.store';

export function useLogin() {
  const navigate = useNavigate();
  const location = useLocation();
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (response) => {
      const user = {
        id: response.user_id,
        email: response.email,
        first_name: '',
        last_name: '',
        role: response.role,
        is_verified: true,
        created_at: new Date().toISOString(),
      };
      setAuth(user, {
        access_token: response.access_token,
        refresh_token: response.refresh_token,
        expires_in: 3600,
      });
      const from = (location.state as any)?.from?.pathname || '/';
      navigate(from);
    },
  });
}
