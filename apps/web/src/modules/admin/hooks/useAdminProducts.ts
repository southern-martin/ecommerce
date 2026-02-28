import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminProductApi } from '../services/admin-product.api';
import type { CreateCategoryData, CreateAttributeData, CreateAttributeGroupData } from '../services/admin-product.api';

export function useCategories() {
  return useQuery({
    queryKey: ['categories'],
    queryFn: () => adminProductApi.getCategories(),
  });
}

export function useCreateCategory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateCategoryData) => adminProductApi.createCategory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });
}

export function useUpdateCategory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateCategoryData & { is_active?: boolean }> }) =>
      adminProductApi.updateCategory(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });
}

export function useDeleteCategory() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminProductApi.deleteCategory(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
    },
  });
}

export function useAdminAttributes() {
  return useQuery({
    queryKey: ['admin-attributes'],
    queryFn: () => adminProductApi.getAttributes(),
  });
}

export function useCreateAttribute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateAttributeData) => adminProductApi.createAttribute(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-attributes'] });
    },
  });
}

export function useUpdateAttribute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateAttributeData> }) =>
      adminProductApi.updateAttribute(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-attributes'] });
    },
  });
}

export function useDeleteAttribute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminProductApi.deleteAttribute(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-attributes'] });
    },
  });
}

// Attribute Groups
export function useAttributeGroups() {
  return useQuery({
    queryKey: ['attribute-groups'],
    queryFn: () => adminProductApi.getAttributeGroups(),
  });
}

export function useCreateAttributeGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateAttributeGroupData) => adminProductApi.createAttributeGroup(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attribute-groups'] });
    },
  });
}

export function useUpdateAttributeGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateAttributeGroupData> }) =>
      adminProductApi.updateAttributeGroup(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attribute-groups'] });
    },
  });
}

export function useDeleteAttributeGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminProductApi.deleteAttributeGroup(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attribute-groups'] });
    },
  });
}

export function useGroupAttributes(groupId: string) {
  return useQuery({
    queryKey: ['group-attributes', groupId],
    queryFn: () => adminProductApi.getGroupAttributes(groupId),
    enabled: !!groupId,
  });
}

export function useAssignAttributeToGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      groupId,
      attribute_id,
      sort_order,
    }: {
      groupId: string;
      attribute_id: string;
      sort_order?: number;
    }) => adminProductApi.addAttributeToGroup(groupId, { attribute_id, sort_order }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['group-attributes'] });
      queryClient.invalidateQueries({ queryKey: ['attribute-groups'] });
    },
  });
}

export function useRemoveAttributeFromGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ groupId, attrId }: { groupId: string; attrId: string }) =>
      adminProductApi.removeAttributeFromGroup(groupId, attrId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['group-attributes'] });
      queryClient.invalidateQueries({ queryKey: ['attribute-groups'] });
    },
  });
}
