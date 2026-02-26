import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';
import { Card } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';

export default function ShippingTrackingPage() {
  const { orderId } = useParams<{ orderId: string }>();

  const { data: shipment, isLoading } = useQuery({
    queryKey: ['shipping', orderId],
    queryFn: async () => {
      const res = await apiClient.get(`/shipping/order/${orderId}`);
      return res.data;
    },
    enabled: !!orderId,
  });

  if (isLoading) return <div className="p-6">Loading tracking info...</div>;

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Shipping Tracking</h1>
      <Card className="p-6">
        {shipment ? (
          <div className="space-y-4">
            <div className="flex justify-between">
              <span className="text-muted-foreground">Order ID</span>
              <span className="font-mono text-sm">{orderId}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Status</span>
              <Badge variant="outline">{(shipment as any).status || 'pending'}</Badge>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Carrier</span>
              <span>{(shipment as any).carrier || 'N/A'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-muted-foreground">Tracking Number</span>
              <span className="font-mono text-sm">{(shipment as any).tracking_number || 'Pending'}</span>
            </div>
          </div>
        ) : (
          <p className="text-muted-foreground text-center">No shipping information available for this order.</p>
        )}
      </Card>
    </div>
  );
}
