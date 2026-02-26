import { Badge } from '@/shared/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { formatPrice, formatDate } from '@/shared/lib/utils';
import type { Order } from '@/modules/checkout/services/order.api';

interface SellerOrderTableProps {
  orders: Order[];
}

export function SellerOrderTable({ orders }: SellerOrderTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Order</TableHead>
          <TableHead>Customer</TableHead>
          <TableHead>Date</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Items</TableHead>
          <TableHead className="text-right">Total</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {orders.map((order) => (
          <TableRow key={order.id}>
            <TableCell className="font-medium">#{order.order_number}</TableCell>
            <TableCell>
              {order.shipping_address.first_name} {order.shipping_address.last_name}
            </TableCell>
            <TableCell className="text-muted-foreground">{formatDate(order.created_at)}</TableCell>
            <TableCell>
              <Badge>{order.status}</Badge>
            </TableCell>
            <TableCell>{order.items.length}</TableCell>
            <TableCell className="text-right font-medium">{formatPrice(order.total)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
