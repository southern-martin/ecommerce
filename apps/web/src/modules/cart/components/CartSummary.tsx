import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Separator } from '@/shared/components/ui/separator';
import { formatPrice } from '@/shared/lib/utils';

interface CartSummaryProps {
  subtotal: number;
  itemCount: number;
  estimatedShipping?: number;
  showCheckoutButton?: boolean;
}

export function CartSummary({
  subtotal,
  itemCount,
  estimatedShipping = 0,
  showCheckoutButton = true,
}: CartSummaryProps) {
  const total = subtotal + estimatedShipping;

  return (
    <div className="rounded-lg border bg-card p-6">
      <h3 className="text-lg font-semibold">Order Summary</h3>

      <div className="mt-4 space-y-3">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">
            Subtotal ({itemCount} item{itemCount !== 1 ? 's' : ''})
          </span>
          <span>{formatPrice(subtotal)}</span>
        </div>

        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">Estimated Shipping</span>
          <span>{estimatedShipping > 0 ? formatPrice(estimatedShipping) : 'Free'}</span>
        </div>

        <Separator />

        <div className="flex justify-between font-semibold">
          <span>Total</span>
          <span>{formatPrice(total)}</span>
        </div>
      </div>

      {showCheckoutButton && (
        <Button asChild className="mt-6 w-full" size="lg">
          <Link to="/checkout">Proceed to Checkout</Link>
        </Button>
      )}
    </div>
  );
}
