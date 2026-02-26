import { Avatar, AvatarFallback, AvatarImage } from '@/shared/components/ui/avatar';
import { Button } from '@/shared/components/ui/button';
import { Star, ThumbsUp } from 'lucide-react';
import { formatDate } from '@/shared/lib/utils';
import type { Review } from '../services/review.api';

interface ReviewCardProps {
  review: Review;
  onMarkHelpful?: (reviewId: string) => void;
}

export function ReviewCard({ review, onMarkHelpful }: ReviewCardProps) {
  return (
    <div className="space-y-3 py-4">
      <div className="flex items-center gap-3">
        <Avatar className="h-8 w-8">
          <AvatarImage src={review.user_avatar} />
          <AvatarFallback>{review.user_name[0]}</AvatarFallback>
        </Avatar>
        <div>
          <p className="text-sm font-medium">{review.user_name}</p>
          <p className="text-xs text-muted-foreground">{formatDate(review.created_at)}</p>
        </div>
      </div>

      <div className="flex items-center gap-1">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            className={`h-4 w-4 ${
              i < review.rating ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground/30'
            }`}
          />
        ))}
      </div>

      <h4 className="font-medium">{review.title}</h4>
      <p className="text-sm text-muted-foreground">{review.comment}</p>

      {review.images && review.images.length > 0 && (
        <div className="flex gap-2">
          {review.images.map((img, i) => (
            <img key={i} src={img} alt="" className="h-16 w-16 rounded-md object-cover" />
          ))}
        </div>
      )}

      <Button
        variant="ghost"
        size="sm"
        className="text-muted-foreground"
        onClick={() => onMarkHelpful?.(review.id)}
      >
        <ThumbsUp className="mr-1 h-3.5 w-3.5" />
        Helpful ({review.helpful_count})
      </Button>
    </div>
  );
}
