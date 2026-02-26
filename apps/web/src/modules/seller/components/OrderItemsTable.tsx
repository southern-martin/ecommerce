import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { formatPrice } from '@/shared/lib/utils';

interface OrderItem {
  product_name?: string;
  name?: string;
  quantity: number;
  unit_price_cents?: number;
  price?: number;
  image_url?: string;
}

interface OrderItemsTableProps {
  items: OrderItem[];
}

export function OrderItemsTable({ items }: OrderItemsTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Product</TableHead>
          <TableHead className="text-center">Quantity</TableHead>
          <TableHead className="text-right">Unit Price</TableHead>
          <TableHead className="text-right">Subtotal</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item, index) => {
          const name = item.product_name || item.name || 'Unknown Product';
          const unitPriceCents = item.unit_price_cents ?? item.price ?? 0;
          const subtotalCents = unitPriceCents * item.quantity;

          return (
            <TableRow key={index}>
              <TableCell>
                <div className="flex items-center gap-3">
                  {item.image_url && (
                    <img
                      src={item.image_url}
                      alt={name}
                      className="h-10 w-10 rounded object-cover"
                    />
                  )}
                  <span className="font-medium">{name}</span>
                </div>
              </TableCell>
              <TableCell className="text-center">{item.quantity}</TableCell>
              <TableCell className="text-right">{formatPrice(unitPriceCents)}</TableCell>
              <TableCell className="text-right font-medium">{formatPrice(subtotalCents)}</TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
