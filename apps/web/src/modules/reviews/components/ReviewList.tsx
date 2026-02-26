import { Separator } from '@/shared/components/ui/separator';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ReviewCard } from './ReviewCard';
import { reviewApi } from '../services/review.api';
import type { Review } from '../services/review.api';

interface ReviewListProps {
  reviews: Review[];
  isLoading?: boolean;
}

export function ReviewList({ reviews, isLoading }: ReviewListProps) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        {Array.from({ length: 3 }).map((_, i) => (
          <div key={i} className="space-y-2">
            <Skeleton className="h-4 w-32" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-3/4" />
          </div>
        ))}
      </div>
    );
  }

  if (reviews.length === 0) {
    return (
      <p className="py-8 text-center text-muted-foreground">
        No reviews yet. Be the first to review this product.
      </p>
    );
  }

  return (
    <div className="divide-y">
      {reviews.map((review) => (
        <ReviewCard
          key={review.id}
          review={review}
          onMarkHelpful={(id) => reviewApi.markHelpful(id)}
        />
      ))}
    </div>
  );
}
