import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { adminCmsApi } from '../services/admin-cms.api';
import type { CreateBannerData, CreatePageData } from '../services/admin-cms.api';

// Banners
export function useAdminBanners() {
  return useQuery({
    queryKey: ['admin-banners'],
    queryFn: () => adminCmsApi.getBanners(),
  });
}

export function useCreateBanner() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateBannerData) => adminCmsApi.createBanner(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-banners'] });
    },
  });
}

export function useUpdateBanner() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreateBannerData> }) =>
      adminCmsApi.updateBanner(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-banners'] });
    },
  });
}

export function useDeleteBanner() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminCmsApi.deleteBanner(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-banners'] });
    },
  });
}

// Pages
export function useAdminPages() {
  return useQuery({
    queryKey: ['admin-pages'],
    queryFn: () => adminCmsApi.getPages(),
  });
}

export function useCreatePage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePageData) => adminCmsApi.createPage(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-pages'] });
    },
  });
}

export function useUpdatePage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<CreatePageData> }) =>
      adminCmsApi.updatePage(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-pages'] });
    },
  });
}

export function useDeletePage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminCmsApi.deletePage(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-pages'] });
    },
  });
}

export function usePublishPage() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminCmsApi.publishPage(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-pages'] });
    },
  });
}
