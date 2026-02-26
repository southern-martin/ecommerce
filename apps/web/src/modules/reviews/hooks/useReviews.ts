import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query';
import { reviewApi } from '../services/review.api';
import type { CreateReviewData } from '../services/review.api';

export function useReviews(productId: string, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['reviews', productId, page, pageSize],
    queryFn: () => reviewApi.getProductReviews(productId, { page, page_size: pageSize }),
    enabled: !!productId,
    placeholderData: keepPreviousData,
  });
}

export function useCreateReview() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateReviewData) => reviewApi.createReview(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['reviews', variables.product_id] });
    },
  });
}
