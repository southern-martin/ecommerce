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
import { cn } from '@/shared/lib/utils';
import type { WalletTransaction } from '../services/seller-wallet.api';

interface WalletTransactionTableProps {
  transactions: WalletTransaction[];
}

export function WalletTransactionTable({ transactions }: WalletTransactionTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Type</TableHead>
          <TableHead>Amount</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Reference</TableHead>
          <TableHead>Date</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {transactions.map((tx) => (
          <TableRow key={tx.id}>
            <TableCell>
              <StatusBadge status={tx.type} />
            </TableCell>
            <TableCell
              className={cn(
                'font-medium',
                tx.amount_cents >= 0 ? 'text-green-600' : 'text-red-600'
              )}
            >
              {tx.amount_cents >= 0 ? '+' : ''}
              {formatPrice(tx.amount_cents)}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {tx.description || '—'}
            </TableCell>
            <TableCell className="font-mono text-xs text-muted-foreground">
              {tx.reference_id ? truncate(tx.reference_id, 12) : '—'}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {formatDateTime(tx.created_at)}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
