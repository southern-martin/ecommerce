import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { formatDate } from '@/shared/lib/utils';
import type { Shipment } from '../services/seller-shipping.api';

interface ShipmentTableProps {
  shipments: Shipment[];
}

export function ShipmentTable({ shipments }: ShipmentTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Tracking Number</TableHead>
          <TableHead>Carrier</TableHead>
          <TableHead>Order ID</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Est. Delivery</TableHead>
          <TableHead>Created</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {shipments.map((shipment) => (
          <TableRow key={shipment.id}>
            <TableCell className="font-medium font-mono">{shipment.tracking_number}</TableCell>
            <TableCell>{shipment.carrier}</TableCell>
            <TableCell className="text-muted-foreground">{shipment.order_id}</TableCell>
            <TableCell>
              <StatusBadge status={shipment.status} />
            </TableCell>
            <TableCell className="text-muted-foreground">
              {formatDate(shipment.estimated_delivery)}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {formatDate(shipment.created_at)}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
