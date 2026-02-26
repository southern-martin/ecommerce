import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { ChevronLeft, ChevronRight, Gavel, Loader2 } from 'lucide-react';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { useAdminDisputes, useResolveDispute } from '../hooks/useAdminDisputes';
import type { Dispute, ResolveDisputeData } from '../services/admin-dispute.api';

export default function AdminDisputesPage() {
  const [page, setPage] = useState(1);
  const [resolvingDispute, setResolvingDispute] = useState<Dispute | null>(null);
  const [resolutionType, setResolutionType] = useState<ResolveDisputeData['resolution_type']>('refund');
  const [notes, setNotes] = useState('');

  const { data, isLoading } = useAdminDisputes(page);
  const resolveDispute = useResolveDispute();
  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;

  const handleResolve = () => {
    if (!resolvingDispute) return;
    resolveDispute.mutate(
      { id: resolvingDispute.id, data: { resolution_type: resolutionType, notes } },
      {
        onSuccess: () => {
          setResolvingDispute(null);
          setNotes('');
          setResolutionType('refund');
        },
      }
    );
  };

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
      <h1 className="text-2xl font-bold">Disputes</h1>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Order</TableHead>
            <TableHead>Buyer</TableHead>
            <TableHead>Seller</TableHead>
            <TableHead>Reason</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="w-[80px]">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {(!data || data.data.length === 0) ? (
            <TableRow>
              <TableCell colSpan={6} className="text-center text-muted-foreground">
                No disputes found.
              </TableCell>
            </TableRow>
          ) : (
            data.data.map((dispute) => (
              <TableRow key={dispute.id}>
                <TableCell className="font-medium">{dispute.order_number}</TableCell>
                <TableCell>{dispute.buyer_name}</TableCell>
                <TableCell>{dispute.seller_name}</TableCell>
                <TableCell className="max-w-[200px] truncate">{dispute.reason}</TableCell>
                <TableCell>
                  <StatusBadge status={dispute.status} />
                </TableCell>
                <TableCell>
                  {(dispute.status === 'open' || dispute.status === 'under_review') && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setResolvingDispute(dispute)}
                    >
                      <Gavel className="h-4 w-4" />
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

      {/* Resolve Dispute Dialog */}
      <Dialog
        open={!!resolvingDispute}
        onOpenChange={(open) => {
          if (!open) {
            setResolvingDispute(null);
            setNotes('');
            setResolutionType('refund');
          }
        }}
      >
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Resolve Dispute</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            {resolvingDispute && (
              <p className="text-sm text-muted-foreground">
                Dispute for order <strong>{resolvingDispute.order_number}</strong> between{' '}
                <strong>{resolvingDispute.buyer_name}</strong> and{' '}
                <strong>{resolvingDispute.seller_name}</strong>.
              </p>
            )}

            <div className="space-y-2">
              <Label>Resolution Type</Label>
              <Select
                value={resolutionType}
                onValueChange={(value) => setResolutionType(value as any)}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="refund">Full Refund</SelectItem>
                  <SelectItem value="partial_refund">Partial Refund</SelectItem>
                  <SelectItem value="replacement">Replacement</SelectItem>
                  <SelectItem value="rejected">Reject Dispute</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="resolve-notes">Notes</Label>
              <textarea
                id="resolve-notes"
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                rows={4}
                className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                placeholder="Describe the resolution..."
              />
            </div>

            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => setResolvingDispute(null)}>
                Cancel
              </Button>
              <Button onClick={handleResolve} disabled={resolveDispute.isPending || !notes}>
                {resolveDispute.isPending && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Resolve
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
