import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { RETURN_STATUSES } from '@/shared/lib/constants';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { ReturnTable } from '../components/ReturnTable';
import {
  useSellerReturns,
  useApproveReturn,
  useRejectReturn,
} from '../hooks/useSellerReturns';

export default function SellerReturnsPage() {
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState<string | undefined>(undefined);

  const { data, isLoading } = useSellerReturns(page, 10, statusFilter);
  const approveReturn = useApproveReturn();
  const rejectReturn = useRejectReturn();

  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;
  const returns = data?.data ?? [];

  const handleApprove = (id: string) => {
    approveReturn.mutate(id);
  };

  const handleReject = (id: string) => {
    rejectReturn.mutate(id);
  };

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Returns</h1>
        <Select
          value={statusFilter ?? 'all'}
          onValueChange={(val) => {
            setStatusFilter(val === 'all' ? undefined : val);
            setPage(1);
          }}
        >
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Filter by status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            {RETURN_STATUSES.map((status) => (
              <SelectItem key={status} value={status}>
                {status.replace(/_/g, ' ')}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {isLoading ? (
        <Skeleton className="h-64 w-full" />
      ) : returns.length > 0 ? (
        <>
          <ReturnTable
            returns={returns}
            onApprove={handleApprove}
            onReject={handleReject}
            isPending={approveReturn.isPending || rejectReturn.isPending}
          />
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
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
        </>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No returns found.</p>
      )}
    </div>
  );
}
