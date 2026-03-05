import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { formatPrice, formatDateTime, truncate } from '@/shared/lib/utils';
import type { Payout } from '../services/seller-wallet.api';

interface PayoutTableProps {
  payouts: Payout[];
}

export function PayoutTable({ payouts }: PayoutTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>ID</TableHead>
          <TableHead>Amount</TableHead>
          <TableHead>Method</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Requested</TableHead>
          <TableHead>Completed</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {payouts.map((payout) => (
          <TableRow key={payout.id}>
            <TableCell className="font-mono text-xs text-muted-foreground">
              {truncate(payout.id, 12)}
            </TableCell>
            <TableCell className="font-medium">
              {formatPrice(payout.amount_cents, payout.currency)}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {payout.method.replace(/_/g, ' ')}
            </TableCell>
            <TableCell>
              <StatusBadge status={payout.status} />
            </TableCell>
            <TableCell className="text-muted-foreground">
              {formatDateTime(payout.requested_at)}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {payout.completed_at ? formatDateTime(payout.completed_at) : '—'}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
