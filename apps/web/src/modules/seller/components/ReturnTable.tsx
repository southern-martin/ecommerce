import { Button } from '@/shared/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { formatPrice, formatDate } from '@/shared/lib/utils';
import { Check, X } from 'lucide-react';
import type { SellerReturn } from '../services/seller-return.api';

interface ReturnTableProps {
  returns: SellerReturn[];
  onApprove: (id: string) => void;
  onReject: (id: string) => void;
  isPending: boolean;
}

export function ReturnTable({ returns, onApprove, onReject, isPending }: ReturnTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Order</TableHead>
          <TableHead>Reason</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Refund Amount</TableHead>
          <TableHead>Date</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {returns.map((ret) => (
          <TableRow key={ret.id}>
            <TableCell className="font-medium">#{ret.order_number}</TableCell>
            <TableCell className="max-w-[200px] truncate">{ret.reason}</TableCell>
            <TableCell>
              <StatusBadge status={ret.status} />
            </TableCell>
            <TableCell>{formatPrice(ret.refund_amount)}</TableCell>
            <TableCell className="text-muted-foreground">{formatDate(ret.created_at)}</TableCell>
            <TableCell className="text-right">
              {ret.status === 'requested' && (
                <div className="flex justify-end gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onApprove(ret.id)}
                    disabled={isPending}
                  >
                    <Check className="mr-1 h-4 w-4" />
                    Approve
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => onReject(ret.id)}
                    disabled={isPending}
                  >
                    <X className="mr-1 h-4 w-4" />
                    Reject
                  </Button>
                </div>
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
