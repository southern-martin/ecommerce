import { Link, useSearchParams } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent } from '@/shared/components/ui/card';
import { CheckCircle, Package, ShoppingBag } from 'lucide-react';

export default function CheckoutSuccessPage() {
  const [searchParams] = useSearchParams();
  const orderNumber = searchParams.get('order');

  return (
    <div className="mx-auto max-w-lg py-16">
      <Card>
        <CardContent className="flex flex-col items-center p-8 text-center">
          <CheckCircle className="h-16 w-16 text-green-500" />

          <h1 className="mt-6 text-2xl font-bold">Order Confirmed!</h1>

          <p className="mt-2 text-muted-foreground">
            Thank you for your purchase. Your order has been placed successfully.
          </p>

          {orderNumber && (
            <p className="mt-4 text-sm">
              Order number:{' '}
              <span className="font-mono font-semibold">{orderNumber}</span>
            </p>
          )}

          <div className="mt-8 flex w-full flex-col gap-3">
            <Button asChild>
              <Link to={`/account/orders`}>
                <Package className="mr-2 h-4 w-4" />
                View My Orders
              </Link>
            </Button>
            <Button asChild variant="outline">
              <Link to="/products">
                <ShoppingBag className="mr-2 h-4 w-4" />
                Continue Shopping
              </Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
