import { useParams, Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import apiClient from '@/shared/lib/api-client';
import type { Order } from '@/modules/checkout/services/order.api';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Separator } from '@/shared/components/ui/separator';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { formatPrice, formatDateTime } from '@/shared/lib/utils';
import { ArrowLeft } from 'lucide-react';
import { OrderItemsTable } from '../components/OrderItemsTable';
import { OrderStatusUpdater } from '../components/OrderStatusUpdater';
import { CreateShipmentForm } from '../components/CreateShipmentForm';
import { ShipmentInfoCard } from '../components/ShipmentInfoCard';
import { useCreateShipment, useShipmentByOrderId, useSellerCarriers } from '../hooks/useSellerShipping';

async function getSellerOrder(id: string): Promise<Order> {
  const response = await apiClient.get(`/seller/orders/${id}`);
  return response.data.data ?? response.data;
}

async function updateOrderStatus(id: string, status: string): Promise<Order> {
  const response = await apiClient.patch(`/seller/orders/${id}/status`, { status });
  return response.data.data ?? response.data;
}

export default function SellerOrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();

  const { data: order, isLoading } = useQuery({
    queryKey: ['seller-orders', id],
    queryFn: () => getSellerOrder(id!),
    enabled: !!id,
  });

  const updateStatus = useMutation({
    mutationFn: (status: string) => updateOrderStatus(id!, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['seller-orders', id] });
      queryClient.invalidateQueries({ queryKey: ['seller-orders'] });
    },
  });

  const { data: existingShipment, isLoading: shipmentLoading } = useShipmentByOrderId(id);
  const { data: carriers } = useSellerCarriers();
  const createShipment = useCreateShipment();

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!order) {
    return (
      <div className="py-8 text-center text-muted-foreground">Order not found.</div>
    );
  }

  const canCreateShipment =
    !existingShipment &&
    (order.status === 'confirmed' || order.status === 'processing');

  // The backend may return shipping_address in a different shape than the frontend type
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const addr = order.shipping_address as any;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/seller/orders">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Orders
          </Link>
        </Button>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Order #{order.order_number}</h1>
        <span className="text-sm text-muted-foreground">{formatDateTime(order.created_at)}</span>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Order Status</CardTitle>
        </CardHeader>
        <CardContent>
          <OrderStatusUpdater
            currentStatus={order.status}
            onUpdate={(status) => updateStatus.mutate(status)}
            isPending={updateStatus.isPending}
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Shipment</CardTitle>
        </CardHeader>
        <CardContent>
          {shipmentLoading ? (
            <Skeleton className="h-24" />
          ) : existingShipment ? (
            <ShipmentInfoCard shipment={existingShipment} />
          ) : canCreateShipment ? (
            <CreateShipmentForm
              orderId={order.id}
              destinationAddress={{
                street: addr.address_line1 || addr.line1 || '',
                city: addr.city || '',
                state: addr.state || '',
                postal_code: addr.postal_code || '',
                country: addr.country || addr.country_code || '',
              }}
              items={(order.items || []).map((item) => ({
                product_id: item.product_id,
                product_name: item.name,
                quantity: item.quantity,
              }))}
              carriers={carriers ?? []}
              onSubmit={(data) => {
                createShipment.mutate(data, {
                  onSuccess: () => {
                    updateStatus.mutate('shipped');
                    queryClient.invalidateQueries({ queryKey: ['seller-shipment-by-order', id] });
                  },
                });
              }}
              isPending={createShipment.isPending}
            />
          ) : (
            <p className="text-sm text-muted-foreground">
              {order.status === 'pending'
                ? 'Confirm the order before creating a shipment.'
                : 'Shipment has already been created or the order is not eligible.'}
            </p>
          )}
        </CardContent>
      </Card>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Customer</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="font-medium">
              {addr.first_name || addr.full_name}{' '}
              {addr.last_name || ''}
            </p>
            {addr.phone && (
              <p className="text-sm text-muted-foreground">{addr.phone}</p>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Shipping Address</CardTitle>
          </CardHeader>
          <CardContent className="text-sm">
            <p>{addr.address_line1 || addr.line1}</p>
            {(addr.address_line2 || addr.line2) && (
              <p>{addr.address_line2 || addr.line2}</p>
            )}
            <p>
              {addr.city}, {addr.state}{' '}
              {addr.postal_code}
            </p>
            <p>{addr.country || addr.country_code}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Order Items</CardTitle>
        </CardHeader>
        <CardContent>
          <OrderItemsTable items={order.items} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Order Totals</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Subtotal</span>
              <span>{formatPrice(order.subtotal)}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Shipping</span>
              <span>{formatPrice(order.shipping_cost)}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Tax</span>
              <span>{formatPrice(order.tax)}</span>
            </div>
            {order.discount > 0 && (
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Discount</span>
                <span className="text-green-600">-{formatPrice(order.discount)}</span>
              </div>
            )}
            <Separator />
            <div className="flex justify-between font-medium">
              <span>Total</span>
              <span>{formatPrice(order.total)}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
