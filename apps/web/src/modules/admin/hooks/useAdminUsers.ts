import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminUserApi } from '../services/admin-user.api';
import type { UpdateUserData } from '../services/admin-user.api';

export function useAdminUsers(page = 1, pageSize = 10, role?: string) {
  return useQuery({
    queryKey: ['admin-users', page, pageSize, role],
    queryFn: () => adminUserApi.getUsers({ page, page_size: pageSize, role }),
    placeholderData: keepPreviousData,
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserData }) =>
      adminUserApi.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminUserApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-users'] });
    },
  });
}
