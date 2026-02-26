import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Badge } from '@/shared/components/ui/badge';
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
import { Plus, ChevronLeft, ChevronRight } from 'lucide-react';
import { formatDate, formatPrice } from '@/shared/lib/utils';
import { useReturns } from '../hooks/useReturns';

export default function BuyerReturnsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useReturns(page);
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
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">My Returns</h1>
        <Button asChild>
          <Link to="/account/returns/new">
            <Plus className="mr-2 h-4 w-4" />
            New Return
          </Link>
        </Button>
      </div>

      {data && data.data.length > 0 ? (
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Return ID</TableHead>
                <TableHead>Order</TableHead>
                <TableHead>Reason</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Refund</TableHead>
                <TableHead>Date</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.data.map((ret) => (
                <TableRow key={ret.id}>
                  <TableCell className="font-medium">{ret.id.slice(0, 8)}</TableCell>
                  <TableCell>#{ret.order_number}</TableCell>
                  <TableCell>{ret.reason}</TableCell>
                  <TableCell>
                    <Badge>{ret.status}</Badge>
                  </TableCell>
                  <TableCell>{formatPrice(ret.refund_amount)}</TableCell>
                  <TableCell className="text-muted-foreground">{formatDate(ret.created_at)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
              <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm text-muted-foreground">Page {page} of {totalPages}</span>
              <Button variant="outline" size="sm" disabled={page === totalPages} onClick={() => setPage((p) => p + 1)}>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No return requests yet.</p>
      )}
    </div>
  );
}
