import { useQuery } from '@tanstack/react-query';
import { cmsApi } from '../services/cms.api';

export function useBanners() {
  return useQuery({
    queryKey: ['cms', 'banners'],
    queryFn: () => cmsApi.getBanners(),
    staleTime: 5 * 60 * 1000,
  });
}
