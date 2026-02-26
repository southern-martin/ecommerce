import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminProductApi } from '../services/admin-product.api';
import type { CreateCategoryData, CreateAttributeData } from '../services/admin-product.api';

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

export function useCategoryAttributes(categoryId: string) {
  return useQuery({
    queryKey: ['category-attributes', categoryId],
    queryFn: () => adminProductApi.getCategoryAttributes(categoryId),
    enabled: !!categoryId,
  });
}

export function useAssignAttribute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ categoryId, attribute_id }: { categoryId: string; attribute_id: string }) =>
      adminProductApi.assignAttribute(categoryId, { attribute_id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['category-attributes'] });
    },
  });
}

export function useRemoveAttribute() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ categoryId, attrId }: { categoryId: string; attrId: string }) =>
      adminProductApi.removeAttribute(categoryId, attrId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['category-attributes'] });
    },
  });
}
