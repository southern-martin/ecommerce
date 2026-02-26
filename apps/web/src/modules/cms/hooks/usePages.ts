import { useQuery } from '@tanstack/react-query';
import { cmsApi } from '../services/cms.api';

export function usePages() {
  return useQuery({
    queryKey: ['cms', 'pages'],
    queryFn: () => cmsApi.getPages(),
  });
}

export function usePage(slug: string) {
  return useQuery({
    queryKey: ['cms', 'page', slug],
    queryFn: () => cmsApi.getPageBySlug(slug),
    enabled: !!slug,
  });
}
