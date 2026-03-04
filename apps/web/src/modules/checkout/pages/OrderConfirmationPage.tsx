import { useParams, Link, useLocation } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent } from '@/shared/components/ui/card';
import { Separator } from '@/shared/components/ui/separator';
import { CheckCircle, Package, ShoppingBag } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { formatPrice } from '@/shared/lib/utils';
import type { Order } from '../services/order.api';

export default function OrderConfirmationPage() {
  const { orderId } = useParams<{ orderId: string }>();
  const location = useLocation();
  const order = (location.state as { order?: Order } | null)?.order;

  return (
    <PageLayout
      breadcrumbs={[
        { label: 'Orders', href: '/account/orders' },
        { label: 'Confirmation' },
      ]}
    >
      <div className="mx-auto max-w-lg py-8 space-y-6">
        <Card className="rounded-2xl border bg-card">
          <CardContent className="flex flex-col items-center p-8 text-center">
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-green-50">
              <CheckCircle className="h-12 w-12 text-green-500" />
            </div>

            <h1 className="mt-6 text-2xl font-bold">Order Confirmed!</h1>

            <p className="mt-2 text-muted-foreground">
              Your order has been placed successfully.
            </p>

            <p className="mt-4 text-sm">
              Order ID:{' '}
              <span className="rounded-lg bg-muted px-2 py-1 font-mono font-semibold">
                {orderId}
              </span>
            </p>
          </CardContent>
        </Card>

        {order && order.items && order.items.length > 0 && (
          <Card className="rounded-2xl border bg-card">
            <CardContent className="p-6">
              <h3 className="text-sm font-semibold">Items Ordered</h3>
              <div className="mt-3 space-y-3">
                {order.items.map((item, idx) => (
                  <div key={item.id || idx} className="flex items-center gap-3 text-sm">
                    {item.image_url && (
                      <img
                        src={item.image_url}
                        alt={item.name}
                        className="h-10 w-10 rounded-lg object-cover bg-muted"
                      />
                    )}
                    <div className="flex-1">
                      <p className="font-medium">{item.name}</p>
                      {item.variant_name && (
                        <p className="text-xs text-muted-foreground">{item.variant_name}</p>
                      )}
                      <p className="text-xs text-muted-foreground">Qty: {item.quantity}</p>
                    </div>
                    <span className="font-medium">{formatPrice(item.price * item.quantity)}</span>
                  </div>
                ))}
              </div>

              <Separator className="my-4" />

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Subtotal</span>
                  <span>{formatPrice(order.subtotal)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Shipping</span>
                  <span>{order.shipping_cost > 0 ? formatPrice(order.shipping_cost) : 'Free'}</span>
                </div>
                {order.tax > 0 && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Tax</span>
                    <span>{formatPrice(order.tax)}</span>
                  </div>
                )}
                {order.discount > 0 && (
                  <div className="flex justify-between text-green-600">
                    <span>Discount</span>
                    <span>-{formatPrice(order.discount)}</span>
                  </div>
                )}
                <Separator />
                <div className="flex justify-between text-base font-bold">
                  <span>Total</span>
                  <span>{formatPrice(order.total)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {order?.shipping_address && (
          <Card className="rounded-2xl border bg-card">
            <CardContent className="p-6">
              <h3 className="text-sm font-semibold">Shipping To</h3>
              <p className="mt-2 text-sm text-muted-foreground">
                {order.shipping_address.first_name} {order.shipping_address.last_name}
                <br />
                {order.shipping_address.address_line1}
                {order.shipping_address.address_line2 && (
                  <>, {order.shipping_address.address_line2}</>
                )}
                <br />
                {order.shipping_address.city}, {order.shipping_address.state}{' '}
                {order.shipping_address.postal_code}
                <br />
                {order.shipping_address.country}
              </p>
            </CardContent>
          </Card>
        )}

        <div className="flex flex-col gap-3">
          <Button asChild className="rounded-xl font-semibold" size="lg">
            <Link to="/account/orders">
              <Package className="mr-2 h-4 w-4" />
              View Orders
            </Link>
          </Button>
          <Button asChild variant="outline" className="rounded-xl" size="lg">
            <Link to="/products">
              <ShoppingBag className="mr-2 h-4 w-4" />
              Continue Shopping
            </Link>
          </Button>
        </div>
      </div>
    </PageLayout>
  );
}
