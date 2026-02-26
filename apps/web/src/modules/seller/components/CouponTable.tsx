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
import { ConfirmDialog } from '@/shared/components/data/ConfirmDialog';
import { formatDate } from '@/shared/lib/utils';
import { Pencil, Trash2 } from 'lucide-react';
import type { Coupon } from '../services/seller-coupon.api';

interface CouponTableProps {
  coupons: Coupon[];
  onEdit: (coupon: Coupon) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

function formatCouponValue(coupon: Coupon): string {
  if (coupon.type === 'percentage') return `${coupon.value}%`;
  if (coupon.type === 'fixed_amount') return `$${(coupon.value / 100).toFixed(2)}`;
  return 'Free Shipping';
}

function getCouponStatus(coupon: Coupon): string {
  if (!coupon.is_active) return 'inactive';
  if (new Date(coupon.expires_at) < new Date()) return 'expired';
  return 'active';
}

export function CouponTable({ coupons, onEdit, onDelete, isDeleting }: CouponTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Code</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Value</TableHead>
          <TableHead>Usage</TableHead>
          <TableHead>Starts</TableHead>
          <TableHead>Expires</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {coupons.map((coupon) => (
          <TableRow key={coupon.id}>
            <TableCell className="font-medium font-mono">{coupon.code}</TableCell>
            <TableCell className="capitalize">{coupon.type.replace(/_/g, ' ')}</TableCell>
            <TableCell>{formatCouponValue(coupon)}</TableCell>
            <TableCell>
              {coupon.uses_count} / {coupon.max_uses}
            </TableCell>
            <TableCell className="text-muted-foreground">{formatDate(coupon.starts_at)}</TableCell>
            <TableCell className="text-muted-foreground">{formatDate(coupon.expires_at)}</TableCell>
            <TableCell>
              <StatusBadge status={getCouponStatus(coupon)} />
            </TableCell>
            <TableCell className="text-right">
              <div className="flex justify-end gap-2">
                <Button variant="ghost" size="sm" onClick={() => onEdit(coupon)}>
                  <Pencil className="h-4 w-4" />
                </Button>
                <ConfirmDialog
                  title="Delete Coupon"
                  description={`Are you sure you want to delete coupon "${coupon.code}"? This action cannot be undone.`}
                  confirmLabel="Delete"
                  onConfirm={() => onDelete(coupon.id)}
                  isPending={isDeleting}
                  trigger={
                    <Button variant="ghost" size="sm">
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  }
                />
              </div>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
