import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

export function RoleGuard({ allowedRoles, roles, children }: { allowedRoles?: string[]; roles?: string[]; children?: React.ReactNode }) {
  const { user, isAuthenticated } = useAuth();
  const roleList = allowedRoles || roles || [];
  if (!isAuthenticated) return <Navigate to="/login" replace />;
  if (!user || !roleList.includes(user.role)) return <Navigate to="/" replace />;
  return children ? <>{children}</> : <Outlet />;
}
