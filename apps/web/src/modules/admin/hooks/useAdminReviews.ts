import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { adminReviewApi } from '../services/admin-review.api';

export function useAdminReviews(page = 1) {
  return useQuery({
    queryKey: ['admin-reviews', page],
    queryFn: () => adminReviewApi.getReviews({ page, page_size: 10 }),
    placeholderData: keepPreviousData,
  });
}

export function useApproveReview() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => adminReviewApi.approveReview(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-reviews'] });
    },
  });
}
