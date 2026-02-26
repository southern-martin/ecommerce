import { Link } from 'react-router-dom';
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

interface OrderListProps {
  orders: Order[];
}

const statusVariant: Record<string, 'default' | 'secondary' | 'destructive' | 'outline'> = {
  pending: 'outline',
  confirmed: 'secondary',
  processing: 'secondary',
  shipped: 'default',
  delivered: 'default',
  cancelled: 'destructive',
  refunded: 'destructive',
};

export function OrderList({ orders }: OrderListProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Order</TableHead>
          <TableHead>Date</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Items</TableHead>
          <TableHead className="text-right">Total</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {orders.map((order) => (
          <TableRow key={order.id}>
            <TableCell>
              <Link to={`/account/orders/${order.id}`} className="font-medium hover:text-primary">
                #{order.order_number}
              </Link>
            </TableCell>
            <TableCell className="text-muted-foreground">{formatDate(order.created_at)}</TableCell>
            <TableCell>
              <Badge variant={statusVariant[order.status] ?? 'outline'}>
                {order.status}
              </Badge>
            </TableCell>
            <TableCell>{order.items.length}</TableCell>
            <TableCell className="text-right font-medium">{formatPrice(order.total)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
