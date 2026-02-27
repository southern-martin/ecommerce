import { Link, useSearchParams } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent } from '@/shared/components/ui/card';
import { CheckCircle, Package, ShoppingBag } from 'lucide-react';
import { PageLayout } from '@/shared/components/layout/PageLayout';

export default function CheckoutSuccessPage() {
  const [searchParams] = useSearchParams();
  const orderNumber = searchParams.get('order');

  return (
    <PageLayout
      breadcrumbs={[
        { label: 'Orders', href: '/account/orders' },
        { label: 'Confirmation' },
      ]}
    >
      <div className="mx-auto max-w-lg py-8">
        <Card className="rounded-2xl border bg-card">
          <CardContent className="flex flex-col items-center p-8 text-center">
            <div className="flex h-20 w-20 items-center justify-center rounded-full bg-green-50">
              <CheckCircle className="h-12 w-12 text-green-500" />
            </div>

            <h1 className="mt-6 text-2xl font-bold">Order Confirmed!</h1>

            <p className="mt-2 text-muted-foreground">
              Thank you for your purchase. Your order has been placed successfully.
            </p>

            {orderNumber && (
              <p className="mt-4 text-sm">
                Order number:{' '}
                <span className="rounded-lg bg-muted px-2 py-1 font-mono font-semibold">
                  {orderNumber}
                </span>
              </p>
            )}

            <div className="mt-8 flex w-full flex-col gap-3">
              <Button asChild className="rounded-xl font-semibold" size="lg">
                <Link to="/account/orders">
                  <Package className="mr-2 h-4 w-4" />
                  View My Orders
                </Link>
              </Button>
              <Button asChild variant="outline" className="rounded-xl" size="lg">
                <Link to="/products">
                  <ShoppingBag className="mr-2 h-4 w-4" />
                  Continue Shopping
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </PageLayout>
  );
}
