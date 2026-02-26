import { useQuery } from '@tanstack/react-query';
import { adminUserApi } from '../services/admin-user.api';

export function useAdminDashboard() {
  return useQuery({
    queryKey: ['admin-dashboard'],
    queryFn: () => adminUserApi.getDashboardStats(),
  });
}
