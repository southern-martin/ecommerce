import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { formatDate } from '@/shared/lib/utils';
import type { Shipment } from '../services/seller-shipping.api';

interface ShipmentInfoCardProps {
  shipment: Shipment;
}

export function ShipmentInfoCard({ shipment }: ShipmentInfoCardProps) {
  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm text-muted-foreground">Tracking Number</p>
          <p className="font-mono font-medium">{shipment.tracking_number || 'Pending'}</p>
        </div>
        <StatusBadge status={shipment.status} />
      </div>
      <div className="grid gap-4 sm:grid-cols-3">
        <div>
          <p className="text-sm text-muted-foreground">Carrier</p>
          <p className="text-sm font-medium">{shipment.carrier}</p>
        </div>
        <div>
          <p className="text-sm text-muted-foreground">Est. Delivery</p>
          <p className="text-sm font-medium">
            {shipment.estimated_delivery ? formatDate(shipment.estimated_delivery) : '—'}
          </p>
        </div>
        <div>
          <p className="text-sm text-muted-foreground">Created</p>
          <p className="text-sm font-medium">{formatDate(shipment.created_at)}</p>
        </div>
      </div>
    </div>
  );
}
