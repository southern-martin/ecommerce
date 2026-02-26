import { useAuthStore } from '../stores/auth.store';

export function useAuth() {
  const { user, isAuthenticated, setAuth, setUser, logout } = useAuthStore();
  return {
    user,
    isAuthenticated,
    isAdmin: user?.role === 'admin',
    isSeller: user?.role === 'seller',
    isBuyer: user?.role === 'buyer',
    setAuth,
    setUser,
    logout,
  };
}
