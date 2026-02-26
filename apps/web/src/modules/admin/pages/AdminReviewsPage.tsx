import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { ChevronLeft, ChevronRight, CheckCircle, Star } from 'lucide-react';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { useAdminReviews, useApproveReview } from '../hooks/useAdminReviews';
import { truncate } from '@/shared/lib/utils';

function StarRating({ rating }: { rating: number }) {
  return (
    <div className="flex items-center gap-0.5">
      {Array.from({ length: 5 }).map((_, i) => (
        <Star
          key={i}
          className={`h-4 w-4 ${
            i < rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'
          }`}
        />
      ))}
    </div>
  );
}

export default function AdminReviewsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useAdminReviews(page);
  const approveReview = useApproveReview();
  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Reviews</h1>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Product</TableHead>
            <TableHead>User</TableHead>
            <TableHead>Rating</TableHead>
            <TableHead>Comment</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="w-[80px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {(!data || data.data.length === 0) ? (
            <TableRow>
              <TableCell colSpan={6} className="text-center text-muted-foreground">
                No reviews found.
              </TableCell>
            </TableRow>
          ) : (
            data.data.map((review) => (
              <TableRow key={review.id}>
                <TableCell className="font-medium">{review.product_name}</TableCell>
                <TableCell>{review.user_name}</TableCell>
                <TableCell>
                  <StarRating rating={review.rating} />
                </TableCell>
                <TableCell className="max-w-[250px]">
                  {truncate(review.comment, 80)}
                </TableCell>
                <TableCell>
                  <StatusBadge status={review.status} />
                </TableCell>
                <TableCell>
                  {review.status === 'pending' && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => approveReview.mutate(review.id)}
                      disabled={approveReview.isPending}
                    >
                      <CheckCircle className="h-4 w-4 text-green-600" />
                    </Button>
                  )}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            <ChevronLeft className="h-4 w-4" />
          </Button>
          <span className="text-sm text-muted-foreground">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page === totalPages}
            onClick={() => setPage((p) => p + 1)}
          >
            <ChevronRight className="h-4 w-4" />
          </Button>
        </div>
      )}
    </div>
  );
}
